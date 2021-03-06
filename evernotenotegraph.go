package main

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/dreampuf/evernote-sdk-golang/edam"
	"github.com/sirupsen/logrus"
)

// DefaultPageSize specifies the default number of notes to fetch metadata from the Evernote API from in one go
const DefaultPageSize = 100

// EvernoteNoteGraph generates a NoteGraph of all Evernote notes and stores the graph as GraphML document
type EvernoteNoteGraph struct {
	EvernoteClient IEvernoteClient
	NoteLinkParser *NoteLinkParser
	NoteURLType    URLType
	GraphMLUtil    *GraphMLUtil
	PageSize       int32
}

// NewEvernoteNoteGraph creates a new instance of EvernoteNoteGraph
func NewEvernoteNoteGraph(evernoteClient IEvernoteClient, noteLinkParser *NoteLinkParser, noteURLType URLType) *EvernoteNoteGraph {
	return &EvernoteNoteGraph{
		EvernoteClient: evernoteClient,
		NoteLinkParser: noteLinkParser,
		NoteURLType:    noteURLType,
		GraphMLUtil:    &GraphMLUtil{},
		PageSize:       DefaultPageSize}
}

// SetPageSize sets the pagesize
func (eng *EvernoteNoteGraph) SetPageSize(pageSize int32) {
	eng.PageSize = pageSize
}

// GetPageSize gets the pagesize
func (eng *EvernoteNoteGraph) GetPageSize() int32 {
	return eng.PageSize
}

// CreateNoteGraph creates a NoteGraph based on all Evernote notes in the Evernote account
func (eng *EvernoteNoteGraph) CreateNoteGraph() (*NoteGraph, error) {
	offset := int32(0)
	noteGraph := NewNoteGraph()
	for {
		logrus.Infof("Processing metadata of Evernote notes from offset [%d] with page size [%d]", offset, eng.PageSize)
		evernoteNoteMetadataList, err := eng.EvernoteClient.FindAllNotesMetadata(offset, eng.PageSize)
		if err != nil {
			return nil, fmt.Errorf("Failed to process metadata of Evernote notes from offset [%d] with page size [%d]: %w", offset, eng.PageSize, err)
		}

		for _, evernoteNoteMetadata := range evernoteNoteMetadataList.GetNotes() {
			note, noteLinks, err := eng.ProcessEvernoteNote(evernoteNoteMetadata)
			if err != nil {
				return nil, fmt.Errorf("Failed to process Evernote note with GUID [%s] and title [%s]: %w", evernoteNoteMetadata.GetGUID(), evernoteNoteMetadata.GetTitle(), err)
			}

			noteGraph.Add(*note, noteLinks)
		}

		remainingNotes := evernoteNoteMetadataList.TotalNotes - (evernoteNoteMetadataList.StartIndex + int32(len(evernoteNoteMetadataList.Notes)))
		if remainingNotes == 0 {
			break
		}

		offset += eng.PageSize
	}

	return noteGraph, nil
}

// ProcessEvernoteNote extracts Note and NoteLinks for the NoteGraph from an Evernote note
func (eng *EvernoteNoteGraph) ProcessEvernoteNote(evernoteNoteMetadata *edam.NoteMetadata) (*Note, []NoteLink, error) {
	logrus.Infof("Processing Evernote note with GUID [%s] and title [%s]", evernoteNoteMetadata.GetGUID(), evernoteNoteMetadata.GetTitle())

	evernoteNote, err := eng.FetchNote(evernoteNoteMetadata)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to fetch Evernote note with GUID [%s] and title [%s]: %w", evernoteNoteMetadata.GetGUID(), evernoteNoteMetadata.GetTitle(), err)
	}

	note, err := eng.CreateNote(evernoteNote)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create Note for Evernote note with GUID [%s] and title [%s]: %w", evernoteNote.GetGUID(), evernoteNote.GetTitle(), err)
	}

	noteLinks, err := eng.ExtractNoteLinks(evernoteNote)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to extract NoteLinks from Evernote note with GUID [%s] and title [%s]: %w", evernoteNote.GetGUID(), evernoteNote.GetTitle(), err)
	}

	selectedNoteLinks := eng.SelectNoteLinks(note, noteLinks)
	return note, selectedNoteLinks, nil
}

// CreateNote extracts Note for the NoteGraph from the Evernote note metadata
func (eng *EvernoteNoteGraph) CreateNote(evernoteNote *edam.Note) (*Note, error) {
	logrus.Debugf("Creating Note representation of Evernote note with GUID [%s] and title [%s]", evernoteNote.GetGUID(), evernoteNote.GetTitle())

	noteGUID := string(evernoteNote.GetGUID())
	noteTitle := evernoteNote.GetTitle()
	noteURL, noteURLType, err := eng.CreateNoteURL(noteGUID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Note URL for Evernote note with GUID [%s] and title [%s]: %w", evernoteNote.GetGUID(), evernoteNote.GetTitle(), err)
	}

	return &Note{GUID: noteGUID, Title: noteTitle, Description: noteTitle, URL: *noteURL, URLType: *noteURLType}, nil
}

// CreateNoteURL creates the URL for the Note with EvernoteNoteGraph.NoteURLType
func (eng *EvernoteNoteGraph) CreateNoteURL(noteGUID string) (*url.URL, *URLType, error) {
	if eng.NoteURLType == WebLink {
		noteURL, err := eng.NoteLinkParser.CreateWebLinkURL(noteGUID)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to create WebLink URL for Note with GUID [%s]: %w", noteGUID, err)
		}

		return noteURL, &eng.NoteURLType, nil
	} else if eng.NoteURLType == AppLink {
		noteURL, err := eng.NoteLinkParser.CreateAppLinkURL(noteGUID)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to create AppLink URL for Note with GUID [%s]: %w", noteGUID, err)
		}

		return noteURL, &eng.NoteURLType, nil
	}

	return nil, nil, errors.New("Failed to create URL for Note with GUID [" + noteGUID + "]: Invalid/Unsupported NoteURLType [" + eng.NoteURLType.String() + "]")
}

// FetchNote fetches the Evernote note
func (eng *EvernoteNoteGraph) FetchNote(evernoteNoteMetadata *edam.NoteMetadata) (*edam.Note, error) {
	logrus.Debugf("Fetching Evernote note and note content with GUID [%s] and title [%s]", evernoteNoteMetadata.GetGUID(), evernoteNoteMetadata.GetTitle())
	evernoteNote, err := eng.EvernoteClient.GetNoteWithContent(evernoteNoteMetadata.GetGUID())
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch Evernote note and note content with GUID [%s] and title [%s]: %w", evernoteNoteMetadata.GetGUID(), evernoteNoteMetadata.GetTitle(), err)
	}

	return evernoteNote, nil
}

// ExtractNoteLinks extracts NoteLinks for the NoteGraph from the Evernote note
func (eng *EvernoteNoteGraph) ExtractNoteLinks(evernoteNote *edam.Note) ([]NoteLink, error) {
	logrus.Debugf("Parsing content of Evernote note with GUID [%s] and title [%s]", evernoteNote.GetGUID(), evernoteNote.GetTitle())
	noteLinks, err := eng.NoteLinkParser.ExtractNoteLinks(string(evernoteNote.GetGUID()), evernoteNote.GetContent())
	if err != nil {
		return nil, fmt.Errorf("Failed to parse content of Evernote note with GUID [%s] and title [%s]: %w", evernoteNote.GetGUID(), evernoteNote.GetTitle(), err)
	}

	logrus.Debugf("Detected [%d] NoteLinks in Evernote note with GUID [%s] and title [%s]", len(noteLinks), evernoteNote.GetGUID(), evernoteNote.GetTitle())
	logrus.Tracef("Evernote note with GUID [%s] and title [%s] has [%d] NoteLinks: [%s]", evernoteNote.GetGUID(), evernoteNote.GetTitle(), len(noteLinks), noteLinks)
	return noteLinks, nil
}

// SelectNoteLinks selects the type of NoteLinks to include in the NoteGraph
func (eng *EvernoteNoteGraph) SelectNoteLinks(note *Note, noteLinks []NoteLink) []NoteLink {
	selectedNoteLinks := []NoteLink{}

	// We only include NoteLinks with URLs of type AppLink and WebLink. Including URLs of type ShortenedLink and PublicLink - which may
	// point to notes of other Evernote accounts - would require us to (a) generate Notes (with partial information) and (b) include
	// the generated Notes in the NoteGraph
	for _, noteLink := range noteLinks {
		if noteLink.URLType == AppLink || noteLink.URLType == WebLink {
			selectedNoteLinks = append(selectedNoteLinks, noteLink)
		}
	}

	logrus.Tracef("Selected [%d] out of [%d] links with types AppLink and WebLink for Note with GUID [%s] and title [%s] to be included in NoteGraph", len(noteLinks), len(selectedNoteLinks), note.GUID, note.Title)

	return selectedNoteLinks
}
