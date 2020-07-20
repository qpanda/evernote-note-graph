package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	// two Notes, one broken NoteLink
	noteGraphA := NewNoteGraph()
	assert.False(t, noteGraphA.Add(Note{GUID: "1"}, []NoteLink{}))
	assert.True(t, noteGraphA.Add(Note{GUID: "2"}, []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}}))
	assert.ElementsMatch(t, *noteGraphA.GetNotes(), []Note{{GUID: "1"}, {GUID: "2"}})
	assert.ElementsMatch(t, *noteGraphA.GetNoteLinks(), []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})

	// two Notes, one valid NoteLink
	noteGraphB := NewNoteGraph()
	assert.False(t, noteGraphB.Add(Note{GUID: "1"}, []NoteLink{}))
	assert.True(t, noteGraphB.Add(Note{GUID: "2"}, []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "1"}}))
	assert.ElementsMatch(t, *noteGraphB.GetNotes(), []Note{{GUID: "1"}, {GUID: "2"}})
	assert.ElementsMatch(t, *noteGraphB.GetNoteLinks(), []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "1"}})
}

func TestGetLinkedNotes(t *testing.T) {
	// single Note, no NoteLinks
	noteGraphA := NewNoteGraph()
	noteGraphA.Add(Note{GUID: "1"}, []NoteLink{})
	assert.ElementsMatch(t, *noteGraphA.GetNotes(), []Note{{GUID: "1"}})
	assert.ElementsMatch(t, *noteGraphA.GetLinkedNotes(), []Note{})

	// single Note, with NoteLink to itself
	noteGraphB := NewNoteGraph()
	noteGraphB.Add(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "1"}})
	assert.ElementsMatch(t, *noteGraphB.GetNotes(), []Note{{GUID: "1"}})
	assert.ElementsMatch(t, *noteGraphB.GetLinkedNotes(), []Note{{GUID: "1"}})

	// two Notes, broken NoteLink
	noteGraphC := NewNoteGraph()
	noteGraphC.Add(Note{GUID: "1"}, []NoteLink{})
	noteGraphC.Add(Note{GUID: "2"}, []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})
	assert.ElementsMatch(t, *noteGraphC.GetNotes(), []Note{{GUID: "1"}, {GUID: "2"}})
	assert.ElementsMatch(t, *noteGraphC.GetLinkedNotes(), []Note{})

	// two Notes, valid NoteLink
	noteGraphD := NewNoteGraph()
	noteGraphD.Add(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	noteGraphD.Add(Note{GUID: "2"}, []NoteLink{})
	assert.ElementsMatch(t, *noteGraphD.GetNotes(), []Note{{GUID: "1"}, {GUID: "2"}})
	assert.ElementsMatch(t, *noteGraphD.GetLinkedNotes(), []Note{{GUID: "1"}, {GUID: "2"}})

	// two Notes, valid bidirectional NoteLinks
	noteGraphE := NewNoteGraph()
	noteGraphE.Add(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	noteGraphE.Add(Note{GUID: "2"}, []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "1"}})
	assert.ElementsMatch(t, *noteGraphE.GetNotes(), []Note{{GUID: "1"}, {GUID: "2"}})
	assert.ElementsMatch(t, *noteGraphE.GetLinkedNotes(), []Note{{GUID: "1"}, {GUID: "2"}})

	// three Notes, valid NoteLinks
	noteGraphF := NewNoteGraph()
	noteGraphF.Add(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	noteGraphF.Add(Note{GUID: "2"}, []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})
	noteGraphF.Add(Note{GUID: "3"}, []NoteLink{})
	assert.ElementsMatch(t, *noteGraphF.GetNotes(), []Note{{GUID: "1"}, {GUID: "2"}, {GUID: "3"}})
	assert.ElementsMatch(t, *noteGraphF.GetLinkedNotes(), []Note{{GUID: "1"}, {GUID: "2"}, {GUID: "3"}})

	// four Notes, valid NoteLinks, disconnected graph
	noteGraphG := NewNoteGraph()
	noteGraphG.Add(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	noteGraphG.Add(Note{GUID: "2"}, []NoteLink{})
	noteGraphG.Add(Note{GUID: "3"}, []NoteLink{{SourceNoteGUID: "3", TargetNoteGUID: "4"}})
	noteGraphG.Add(Note{GUID: "4"}, []NoteLink{})
	assert.ElementsMatch(t, *noteGraphG.GetNotes(), []Note{{GUID: "1"}, {GUID: "2"}, {GUID: "3"}, {GUID: "4"}})
	assert.ElementsMatch(t, *noteGraphG.GetLinkedNotes(), []Note{{GUID: "1"}, {GUID: "2"}, {GUID: "3"}, {GUID: "4"}})
}

func TestGetValidNoteLinks(t *testing.T) {
	// single Note, no NoteLinks
	noteGraphA := NewNoteGraph()
	noteGraphA.Add(Note{GUID: "1"}, []NoteLink{})
	assert.ElementsMatch(t, *noteGraphA.GetNoteLinks(), []NoteLink{})
	assert.ElementsMatch(t, *noteGraphA.GetValidNoteLinks(), []NoteLink{})

	// single Note, with NoteLink to itself
	noteGraphB := NewNoteGraph()
	noteGraphB.Add(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "1"}})
	assert.ElementsMatch(t, *noteGraphB.GetNoteLinks(), []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "1"}})
	assert.ElementsMatch(t, *noteGraphB.GetValidNoteLinks(), []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "1"}})

	// two Notes, broken NoteLink
	noteGraphC := NewNoteGraph()
	noteGraphC.Add(Note{GUID: "1"}, []NoteLink{})
	noteGraphC.Add(Note{GUID: "2"}, []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})
	assert.ElementsMatch(t, *noteGraphC.GetNoteLinks(), []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})
	assert.ElementsMatch(t, *noteGraphC.GetValidNoteLinks(), []NoteLink{})

	// two Notes, valid NoteLink
	noteGraphD := NewNoteGraph()
	noteGraphD.Add(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	noteGraphD.Add(Note{GUID: "2"}, []NoteLink{})
	assert.ElementsMatch(t, *noteGraphD.GetNoteLinks(), []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	assert.ElementsMatch(t, *noteGraphD.GetValidNoteLinks(), []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
}

func TestGetBrokenNoteLinks(t *testing.T) {
	// single Note, no NoteLinks
	noteGraphA := NewNoteGraph()
	noteGraphA.Add(Note{GUID: "1"}, []NoteLink{})
	assert.ElementsMatch(t, *noteGraphA.GetNoteLinks(), []NoteLink{})
	assert.ElementsMatch(t, *noteGraphA.GetBrokenNoteLinks(), []NoteLink{})

	// single Note, with NoteLink to itself
	noteGraphB := NewNoteGraph()
	noteGraphB.Add(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "1"}})
	assert.ElementsMatch(t, *noteGraphB.GetNoteLinks(), []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "1"}})
	assert.ElementsMatch(t, *noteGraphB.GetBrokenNoteLinks(), []NoteLink{})

	// two Notes, broken NoteLink
	noteGraphC := NewNoteGraph()
	noteGraphC.Add(Note{GUID: "1"}, []NoteLink{})
	noteGraphC.Add(Note{GUID: "2"}, []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})
	assert.ElementsMatch(t, *noteGraphC.GetNoteLinks(), []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})
	assert.ElementsMatch(t, *noteGraphC.GetBrokenNoteLinks(), []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})

	// two Notes, valid NoteLink
	noteGraphD := NewNoteGraph()
	noteGraphD.Add(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	noteGraphD.Add(Note{GUID: "2"}, []NoteLink{})
	assert.ElementsMatch(t, *noteGraphD.GetNoteLinks(), []NoteLink{{SourceNoteGUID: "1", TargetNoteGUID: "2"}})
	assert.ElementsMatch(t, *noteGraphD.GetBrokenNoteLinks(), []NoteLink{})

	// one Notes, broken NoteLink
	noteGraphE := NewNoteGraph()
	noteGraphE.Add(Note{GUID: "1"}, []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})
	assert.ElementsMatch(t, *noteGraphE.GetNoteLinks(), []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})
	assert.ElementsMatch(t, *noteGraphE.GetBrokenNoteLinks(), []NoteLink{{SourceNoteGUID: "2", TargetNoteGUID: "3"}})
}
