package main

import (
	"testing"

	"github.com/dreampuf/evernote-sdk-golang/edam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEvernoteClient struct {
	mock.Mock
}

func (m *MockEvernoteClient) GetHost() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEvernoteClient) GetUser() (*edam.User, error) {
	args := m.Called()
	return args.Get(0).(*edam.User), args.Error(1)
}

func (m *MockEvernoteClient) GetUserStoreURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEvernoteClient) FindAllNotesMetadata(offset int32, maxNotes int32) (*edam.NotesMetadataList, error) {
	args := m.Called(offset, maxNotes)
	return args.Get(0).(*edam.NotesMetadataList), args.Error(1)
}

func (m *MockEvernoteClient) GetNoteWithContent(guid edam.GUID) (*edam.Note, error) {
	args := m.Called(guid)
	return args.Get(0).(*edam.Note), args.Error(1)
}

func TestSelectNoteLinks(t *testing.T) {
	evernoteNoteGraph := NewEvernoteNoteGraph(nil, nil, WebLink)

	note := &Note{GUID: "1"}
	noteLinks := []NoteLink{{SourceNoteGUID: note.GUID, TargetNoteGUID: "1", URLType: AppLink}, {SourceNoteGUID: note.GUID, TargetNoteGUID: "2", URLType: WebLink}, {SourceNoteGUID: note.GUID, TargetNoteGUID: "3", URLType: PublicLink}, {SourceNoteGUID: note.GUID, TargetNoteGUID: "4", URLType: ShortenedLink}}
	selectedNoteLinks := evernoteNoteGraph.SelectNoteLinks(note, noteLinks)

	assert.ElementsMatch(t, selectedNoteLinks, []NoteLink{{SourceNoteGUID: note.GUID, TargetNoteGUID: "1", URLType: AppLink}, {SourceNoteGUID: note.GUID, TargetNoteGUID: "2", URLType: WebLink}})
}

func TestCreateNote(t *testing.T) {
	noteLinkParser := NewNoteLinkParser(SandboxEvernoteCom, "userId", "shardId")
	evernoteNoteGraph := NewEvernoteNoteGraph(nil, noteLinkParser, WebLink)

	evernoteNoteGUID := "1"
	evernoteNoteTitle := "Test"
	evernoteNoteMetadata := &edam.NoteMetadata{GUID: edam.GUID(evernoteNoteGUID), Title: &evernoteNoteTitle}
	createdNote, err := evernoteNoteGraph.CreateNote(evernoteNoteMetadata)
	if err != nil {
		panic(err)
	}

	url, err := noteLinkParser.CreateWebLinkURL(evernoteNoteGUID)
	if err != nil {
		panic(err)
	}

	expectedNote := &Note{GUID: evernoteNoteGUID, Title: evernoteNoteTitle, Description: evernoteNoteTitle, URL: *url, URLType: WebLink}
	assert.Equal(t, expectedNote, createdNote)
}

func TestCreateNoteURL(t *testing.T) {
	noteLinkParser := NewNoteLinkParser(SandboxEvernoteCom, "userId", "shardId")

	webLinkEvernoteNoteGraph := NewEvernoteNoteGraph(nil, noteLinkParser, WebLink)
	_, webLinkURLType, webLinkErr := webLinkEvernoteNoteGraph.CreateNoteURL("1")
	if webLinkErr != nil {
		panic(webLinkErr)
	}
	assert.Equal(t, WebLink, *webLinkURLType)

	appLinkEvernoteNoteGraph := NewEvernoteNoteGraph(nil, noteLinkParser, AppLink)
	_, appLinkURLType, appLinkErr := appLinkEvernoteNoteGraph.CreateNoteURL("1")
	if appLinkErr != nil {
		panic(webLinkErr)
	}
	assert.Equal(t, AppLink, *appLinkURLType)

	invalidLinkEvernoteNoteGraph := NewEvernoteNoteGraph(nil, noteLinkParser, PublicLink)
	_, _, invalidLinkErr := invalidLinkEvernoteNoteGraph.CreateNoteURL("1")
	assert.NotNil(t, invalidLinkErr)
}

func TestFetchContentAndExtractNoteLinksWithoutNoteLinks(t *testing.T) {
	mockEvernoteClient := new(MockEvernoteClient)
	noteLinkParser := NewNoteLinkParser(SandboxEvernoteCom, "userId", "shardId")
	evernoteNoteGraph := NewEvernoteNoteGraph(mockEvernoteClient, noteLinkParser, WebLink)

	evernoteNoteGUID := edam.GUID("1")
	evernoteNoteTitle := "Test"
	evernoteNoteContent := `<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd"><en-note><div>Test</div></en-note>`
	evernoteNoteMetadata := &edam.NoteMetadata{GUID: evernoteNoteGUID, Title: &evernoteNoteTitle}
	mockEvernoteClient.On("GetNoteWithContent", evernoteNoteGUID).Return(&edam.Note{GUID: &evernoteNoteGUID, Title: &evernoteNoteTitle, Content: &evernoteNoteContent}, nil)

	noteLinks, err := evernoteNoteGraph.FetchContentAndExtractNoteLinks(evernoteNoteMetadata)
	if err != nil {
		panic(err)
	}

	assert.Empty(t, noteLinks)
}

func TestFetchContentAndExtractNoteLinksWithLinks(t *testing.T) {
	mockEvernoteClient := new(MockEvernoteClient)
	noteLinkParser := NewNoteLinkParser(EvernoteCom, "76136038", "s12")
	evernoteNoteGraph := NewEvernoteNoteGraph(mockEvernoteClient, noteLinkParser, WebLink)

	evernoteNoteGUID := edam.GUID("1")
	evernoteNoteTitle := "Test"
	evernoteNoteContent := `<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd"><en-note><div><a href="https://example.org/">NonNoteLink</a></div><div><a href="https://www.evernote.com/shard/s12/nl/76136038/d72dfad0-7d58-41b5-b2c9-4ca434abd543/">WebLink</a></div><div><a href="evernote:///view/76136038/s12/4d971333-8b65-45d6-857b-243c850cabf5/4d971333-8b65-45d6-857b-243c850cabf5/">AppLink</a></div><div><a href="https://www.evernote.com/shard/s12/sh/4d971333-8b65-45d6-857b-243c850cabf5/25771cdb535e9183/">PublicLink</a></div><div><a href="https://www.evernote.com/l/AAxNlxMzi2VF1oV7JDyFDKv1JXcc21NekYM">ShortenedLink</a></div></en-note>`
	evernoteNoteMetadata := &edam.NoteMetadata{GUID: evernoteNoteGUID, Title: &evernoteNoteTitle}
	mockEvernoteClient.On("GetNoteWithContent", evernoteNoteGUID).Return(&edam.Note{GUID: &evernoteNoteGUID, Title: &evernoteNoteTitle, Content: &evernoteNoteContent}, nil)

	noteLinks, err := evernoteNoteGraph.FetchContentAndExtractNoteLinks(evernoteNoteMetadata)
	if err != nil {
		panic(err)
	}

	assert.Len(t, noteLinks, 4)
}

func TestCreateNoteGraphWithNoNotes(t *testing.T) {
	mockEvernoteClient := new(MockEvernoteClient)
	noteLinkParser := NewNoteLinkParser(SandboxEvernoteCom, "userId", "shardId")
	evernoteNoteGraph := NewEvernoteNoteGraph(mockEvernoteClient, noteLinkParser, WebLink)

	mockEvernoteClient.On("FindAllNotesMetadata", int32(0), mock.Anything).Return(&edam.NotesMetadataList{}, nil)

	noteGraph, err := evernoteNoteGraph.CreateNoteGraph()
	if err != nil {
		panic(err)
	}

	assert.Empty(t, noteGraph.Notes)
	assert.Empty(t, noteGraph.NoteLinks)
}

func TestCreateNoteGraphWithNotes(t *testing.T) {
	mockEvernoteClient := new(MockEvernoteClient)
	noteLinkParser := NewNoteLinkParser(EvernoteCom, "76136038", "s12")
	evernoteNoteGraph := NewEvernoteNoteGraph(mockEvernoteClient, noteLinkParser, WebLink)

	offset := int32(0)
	evernoteNoteGUID := edam.GUID("1")
	evernoteNoteTitle := "Test"
	evernoteNoteContent := `<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd"><en-note><div><a href="https://example.org/">NonNoteLink</a></div><div><a href="https://www.evernote.com/shard/s12/nl/76136038/d72dfad0-7d58-41b5-b2c9-4ca434abd543/">WebLink</a></div><div><a href="evernote:///view/76136038/s12/4d971333-8b65-45d6-857b-243c850cabf5/4d971333-8b65-45d6-857b-243c850cabf5/">AppLink</a></div><div><a href="https://www.evernote.com/shard/s12/sh/4d971333-8b65-45d6-857b-243c850cabf5/25771cdb535e9183/">PublicLink</a></div><div><a href="https://www.evernote.com/l/AAxNlxMzi2VF1oV7JDyFDKv1JXcc21NekYM">ShortenedLink</a></div></en-note>`
	evernoteNoteMetadata := []*edam.NoteMetadata{{GUID: evernoteNoteGUID, Title: &evernoteNoteTitle}}
	evernoteNoteMetadataList := &edam.NotesMetadataList{StartIndex: offset, TotalNotes: int32(len(evernoteNoteMetadata)), Notes: evernoteNoteMetadata}

	mockEvernoteClient.On("FindAllNotesMetadata", offset, mock.Anything).Return(evernoteNoteMetadataList, nil)
	mockEvernoteClient.On("GetNoteWithContent", evernoteNoteGUID).Return(&edam.Note{GUID: &evernoteNoteGUID, Title: &evernoteNoteTitle, Content: &evernoteNoteContent}, nil)

	noteGraph, err := evernoteNoteGraph.CreateNoteGraph()
	if err != nil {
		panic(err)
	}

	assert.Len(t, noteGraph.Notes, 1)
	assert.Len(t, noteGraph.NoteLinks, 2)
}
