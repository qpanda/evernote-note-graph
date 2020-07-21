package main

import (
	"github.com/dreampuf/evernote-sdk-golang/edam"
	"github.com/sirupsen/logrus"
)

// PageSize specifies the number of notes to fetch metadata from the Evernote API from in one go
const PageSize = 100

// EvernoteNoteGraph generates a NoteGraph of all Evernote notes and stores the graph as GraphML document
type EvernoteNoteGraph struct {
	EvernoteClient IEvernoteClient
	NoteLinkParser *NoteLinkParser
	GraphMLUtil    *GraphMLUtil
}

// NewEvernoteNoteGraph creates a new instance of EvernoteNoteGraph
func NewEvernoteNoteGraph(evernoteClient IEvernoteClient, noteLinkParser *NoteLinkParser) *EvernoteNoteGraph {
	return &EvernoteNoteGraph{
		EvernoteClient: evernoteClient,
		NoteLinkParser: noteLinkParser,
		GraphMLUtil:    &GraphMLUtil{}}
}

// CreateNoteGraph creates a NoteGraph based on all Evernote notes in the Evernote account
func (eng *EvernoteNoteGraph) CreateNoteGraph() (*NoteGraph, error) {
	offset := int32(0)
	noteGraph := NewNoteGraph()
	for {
		logrus.Debugf("Processing Evernote notes from offset [%d] with page size [%d]", offset, PageSize)
		noteMetadataList, err := eng.EvernoteClient.FindAllNotesMetadata(offset, PageSize)
		if err != nil {
			return nil, err
		}

		for _, noteMetadata := range noteMetadataList.GetNotes() {
			note, noteLinks, err := eng.ProcessEvernoteNote(noteMetadata)
			if err != nil {
				return nil, err
			}

			noteGraph.Add(*note, noteLinks)
		}

		remainingNotes := noteMetadataList.TotalNotes - (noteMetadataList.StartIndex + int32(len(noteMetadataList.Notes)))
		if remainingNotes == 0 {
			break
		}
	}

	return noteGraph, nil
}

// ProcessEvernoteNote extracts Note and NoteLinks for the NoteGraph from an Evernote note
func (eng *EvernoteNoteGraph) ProcessEvernoteNote(noteMetadata *edam.NoteMetadata) (*Note, []NoteLink, error) {
	logrus.Infof("Processing Evernote note with GUID [%s] and title [%s]", noteMetadata.GetGUID(), noteMetadata.GetTitle())

	note, err := eng.CreateNote(noteMetadata)
	if err != nil {
		return nil, nil, err
	}

	noteLinks, err := eng.FetchContentAndExtractNoteLinks(noteMetadata)
	if err != nil {
		return nil, nil, err
	}

	selectedNoteLinks := eng.SelectNoteLinks(note, noteLinks)
	return note, selectedNoteLinks, nil
}

// CreateNote extracts Note for the NoteGraph from the Evernote note metadata
func (eng *EvernoteNoteGraph) CreateNote(noteMetadata *edam.NoteMetadata) (*Note, error) {
	logrus.Debugf("Creating Note representation of Evernote note with GUID [%s] and title [%s]", noteMetadata.GetGUID(), noteMetadata.GetTitle())

	noteGUID := string(noteMetadata.GetGUID())
	noteTitle := noteMetadata.GetTitle()
	noteURL, err := eng.NoteLinkParser.CreateWebLinkURL(noteGUID)
	if err != nil {
		return nil, err
	}

	return &Note{GUID: noteGUID, Title: noteTitle, Description: noteTitle, URL: *noteURL, URLType: WebLink}, nil
}

// FetchContentAndExtractNoteLinks extracts NoteLinks for the NoteGraph from the Evernote note metadata and content
func (eng *EvernoteNoteGraph) FetchContentAndExtractNoteLinks(noteMetadata *edam.NoteMetadata) ([]NoteLink, error) {
	logrus.Debugf("Fetching content of Evernote note with GUID [%s] and title [%s]", noteMetadata.GetGUID(), noteMetadata.GetTitle())
	evernoteNote, err := eng.EvernoteClient.GetNoteWithContent(noteMetadata.GetGUID())
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Parsing ENML of Evernote note with GUID [%s] and title [%s]", noteMetadata.GetGUID(), noteMetadata.GetTitle())
	noteLinks, err := eng.NoteLinkParser.ExtractNoteLinks(string(evernoteNote.GetGUID()), evernoteNote.GetContent())
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Detected [%d] NoteLinks in Evernote note with GUID [%s] and title [%s]", len(noteLinks), evernoteNote.GetGUID(), evernoteNote.GetTitle())
	logrus.Tracef("Evernote note with GUID [%s] and title [%s] has [%d] NoteLinks: [%s]", evernoteNote.GetGUID(), evernoteNote.GetTitle(), len(noteLinks), noteLinks)
	return noteLinks, nil
}

// SelectNoteLinks selects the type of NoteLinks to include in the NoteGraph
func (eng *EvernoteNoteGraph) SelectNoteLinks(note *Note, noteLinks []NoteLink) []NoteLink {
	selectedNoteLinks := []NoteLink{}

	// for now - to keep things simple - we only include NoteLinks of type AppLink and WebLink
	// NoteLinks of type ShortenedLink and PublicLink may point to Notes that are not stored in the users Evernote account, including these types
	// of NotesLinks would require generating Notes with a special type to be included in the NoteGraph
	for _, noteLink := range noteLinks {
		if noteLink.URLType == AppLink || noteLink.URLType == WebLink {
			selectedNoteLinks = append(selectedNoteLinks, noteLink)
		}
	}

	logrus.Tracef("Selected [%d] out of [%d] links with types AppLink and WebLink for Note with GUID [%s] and title [%s] to be included in NoteGraph", len(noteLinks), len(selectedNoteLinks), note.GUID, note.Title)

	return selectedNoteLinks
}
