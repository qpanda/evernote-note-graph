package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNote(t *testing.T) {
	noteGraph := NewNoteGraph()
	assert.False(t, noteGraph.AddNote(Note{GUID: "1"}, []NoteLink{}))
	assert.True(t, noteGraph.AddNote(Note{GUID: "2"}, []NoteLink{{NoteGUID: "3"}}))

	assert.ElementsMatch(t, *noteGraph.GetNotes(), []Note{{GUID: "1"}, {GUID: "2"}})
	assert.ElementsMatch(t, *noteGraph.GetNoteLinks(), []NoteLink{{NoteGUID: "3"}})
}

func TestGetLinkedNotes(t *testing.T) {
	noteGraph := NewNoteGraph()
	noteGraph.AddNote(Note{GUID: "1"}, []NoteLink{{NoteGUID: "2"}})
	noteGraph.AddNote(Note{GUID: "2"}, []NoteLink{})
	noteGraph.AddNote(Note{GUID: "3"}, []NoteLink{})

	assert.ElementsMatch(t, *noteGraph.GetNotes(), []Note{{GUID: "1"}, {GUID: "2"}, {GUID: "3"}})
	assert.ElementsMatch(t, *noteGraph.GetLinkedNotes(), []Note{{GUID: "1"}, {GUID: "2"}})
}

func TestValidate(t *testing.T) {
	noteGraphA := NewNoteGraph()
	noteGraphA.AddNote(Note{GUID: "1"}, []NoteLink{{NoteGUID: "1"}})
	assert.True(t, noteGraphA.Validate())

	noteGraphB := NewNoteGraph()
	noteGraphB.AddNote(Note{GUID: "1"}, []NoteLink{{NoteGUID: "2"}})
	assert.False(t, noteGraphB.Validate())

	noteGraphC := NewNoteGraph()
	noteGraphC.AddNote(Note{GUID: "1"}, []NoteLink{{NoteGUID: "2"}})
	noteGraphC.AddNote(Note{GUID: "2"}, []NoteLink{})
	assert.True(t, noteGraphC.Validate())
}
