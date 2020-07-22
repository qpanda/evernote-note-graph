package main

import (
	"context"
	"os"
	"testing"

	"github.com/dreampuf/evernote-sdk-golang/edam"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testNoteTitle = "Test Note"
var testNoteContent = `<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd"><en-note><div>Test</div></en-note>`

func TestEvernoteClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	if os.Getenv("EDAM_AUTHTOKEN") == "" {
		t.Skip("skipping test because EDAM_AUTHTOKEN env variable is not set")
	}

	authToken := os.Getenv("EDAM_AUTHTOKEN")
	evernoteClient := NewEvernoteClient(authToken, true)

	// ensure we are running tests against the sandbox Evernote API endpoint
	require.Equal(t, SandboxEvernoteCom, evernoteClient.GetHost())

	// create test note
	note, createNoteErr := evernoteClient.CreateNote(testNoteTitle, testNoteContent)
	if createNoteErr != nil {
		panic(createNoteErr)
	}

	// fetch metadata of first 10 notes sorted descending by creation date which should include the newly created test note
	noteMetadataList, findNoteMetadataErr := evernoteClient.FindAllNotesMetadata(0, 10)
	if findNoteMetadataErr != nil {
		panic(findNoteMetadataErr)
	}

	// retrieve content of first 10 notes
	notes := map[edam.GUID]*edam.Note{}
	noteMetadataArray := noteMetadataList.GetNotes()
	for _, noteMetadata := range noteMetadataArray {
		note, getNoteContentErr := evernoteClient.GetNoteWithContent(noteMetadata.GetGUID())
		if getNoteContentErr != nil {
			panic(getNoteContentErr)
		}

		notes[noteMetadata.GetGUID()] = note
	}

	// verify newly created note has been retrieved
	assert.NotNil(t, notes[note.GetGUID()])
	assert.Equal(t, note.GetTitle(), notes[note.GetGUID()].GetTitle())
	assert.Equal(t, note.GetContent(), notes[note.GetGUID()].GetContent())

	// expunge test note
	_, expungeNoteErr := evernoteClient.ExpungeNote(note.GetGUID())
	if expungeNoteErr != nil {
		panic(expungeNoteErr)
	}
}

func (ec *EvernoteClient) CreateNote(title, content string) (*edam.Note, error) {
	noteStoreClient, err := ec.GetNoteStoreClient()
	if err != nil {
		return nil, err
	}

	context, cancelFunc := context.WithTimeout(context.Background(), Timeout)
	defer func() {
		logrus.Tracef("CreateNote context completed")
		cancelFunc()
	}()

	note, err := noteStoreClient.CreateNote(context, ec.AuthToken, &edam.Note{Title: &title, Content: &content})
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (ec *EvernoteClient) ExpungeNote(guid edam.GUID) (*int32, error) {
	noteStoreClient, err := ec.GetNoteStoreClient()
	if err != nil {
		return nil, err
	}

	context, cancelFunc := context.WithTimeout(context.Background(), Timeout)
	defer func() {
		logrus.Tracef("CreateNote context completed")
		cancelFunc()
	}()

	sequenceNumber, err := noteStoreClient.ExpungeNote(context, ec.AuthToken, guid)
	if err != nil {
		return nil, err
	}

	return &sequenceNumber, nil
}
