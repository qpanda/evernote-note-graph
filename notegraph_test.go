package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNote(t *testing.T) {
	noteGraph := NewNoteGraph()
	assert.False(t, noteGraph.AddNote(Note{GUID: "1"}, []NoteLink{}))
	assert.True(t, noteGraph.AddNote(Note{GUID: "2"}, []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}}))

	assert.ElementsMatch(t, *noteGraph.GetNotes(), []Note{{GUID: "1"}, {GUID: "2"}})
	assert.ElementsMatch(t, *noteGraph.GetNoteLinks(), []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})
}

func TestGetLinkedNotes(t *testing.T) {
	noteGraph := NewNoteGraph()
	noteGraph.AddNote(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	noteGraph.AddNote(Note{GUID: "2"}, []NoteLink{})
	noteGraph.AddNote(Note{GUID: "3"}, []NoteLink{})

	assert.ElementsMatch(t, *noteGraph.GetNotes(), []Note{{GUID: "1"}, {GUID: "2"}, {GUID: "3"}})
	assert.ElementsMatch(t, *noteGraph.GetLinkedNotes(), []Note{{GUID: "1"}, {GUID: "2"}})
}

func TestValidate(t *testing.T) {
	noteGraphA := NewNoteGraph()
	noteGraphA.AddNote(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "1"}})
	assert.True(t, noteGraphA.Validate())

	noteGraphB := NewNoteGraph()
	noteGraphB.AddNote(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	assert.False(t, noteGraphB.Validate())

	noteGraphC := NewNoteGraph()
	noteGraphC.AddNote(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	noteGraphC.AddNote(Note{GUID: "2"}, []NoteLink{})
	assert.True(t, noteGraphC.Validate())

	noteGraphD := NewNoteGraph()
	noteGraphD.AddNote(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	noteGraphD.AddNote(Note{GUID: "2"}, []NoteLink{})
	noteGraphD.AddNote(Note{GUID: "3"}, []NoteLink{})
	assert.True(t, noteGraphD.Validate())
}
