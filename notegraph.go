package main

import (
	"fmt"
	"net/url"
)

// Note represents an Evernote note
type Note struct {
	GUID        string
	Title       string
	Description string
	URL         url.URL
}

func (n Note) String() string {
	return fmt.Sprintf("{GUID: %s, Title: %s, Description: %s, URL %s}", n.GUID, n.Title, n.Description, n.URL.String())
}

// LinkType identifies the type of a NoteLink
type LinkType int

// Enum of all LinkTypes (see Evernote API documentation at https://dev.evernote.com/doc/articles/note_links.php)
const (
	AppLink       LinkType = iota // 'In-App Note Link'
	WebLink       LinkType = iota // 'Note Link'
	PublicLink    LinkType = iota // 'Public Link'
	ShortenedLink LinkType = iota // 'Evernote Shortened URLs'
)

func (lt LinkType) String() string {
	return [...]string{"AppLink", "WebLink", "PublicLink", "ShortenedLink"}[lt]
}

// NoteLink is an app, web, public, or shortened link that points to a Note (see Evernote API documentation at https://dev.evernote.com/doc/articles/note_links.php)
type NoteLink struct {
	Type     LinkType
	URL      url.URL
	Text     string
	NoteGUID string // is not set for ShortenedLinks
}

func (nl NoteLink) String() string {
	return fmt.Sprintf("{Type: %s, URL: %s, Text: %s, NoteGUID %s}", nl.Type.String(), nl.URL.String(), nl.Text, nl.NoteGUID)
}

// NoteGraph contains all Notes and NoteLinks and keeps track of which Notes are linked to other Notes
type NoteGraph struct {
	Notes           []Note          // all Notes
	NoteLinks       []NoteLink      // all NoteLinks
	NoteGUIDs       map[string]bool // GUIDs of all notes
	LinkedNoteGUIDs map[string]bool // GUIDs of notes linked to other notes
}

// NewNoteGraph creates a new instance of NoteGraph
func NewNoteGraph() *NoteGraph {
	return &NoteGraph{
		Notes:           []Note{},
		NoteLinks:       []NoteLink{},
		NoteGUIDs:       map[string]bool{},
		LinkedNoteGUIDs: map[string]bool{},
	}
}

// AddNote adds a Note with all its NoteLinks to the NoteGraph, returns true if the Note is a linked note, otherwise false
func (ng *NoteGraph) AddNote(note Note, noteLinks []NoteLink) bool {
	ng.Notes = append(ng.Notes, note)
	ng.NoteGUIDs[note.GUID] = true

	ng.NoteLinks = append(ng.NoteLinks, noteLinks...)
	if len(noteLinks) != 0 {
		// if the Note contains NoteLinks we remember the GUID of the Note and the GUIDs of all notes this Note links to
		ng.LinkedNoteGUIDs[note.GUID] = true
		for _, noteLink := range noteLinks {
			ng.LinkedNoteGUIDs[noteLink.NoteGUID] = true
		}

		return true
	}

	return false
}

// Validate checks whether all NoteLinks point to an existing Note
func (ng *NoteGraph) Validate() bool {
	for linkedNoteGUID := range ng.LinkedNoteGUIDs {
		if !ng.NoteGUIDs[linkedNoteGUID] {
			return false
		}
	}

	return true
}

// GetNotes returns all Notes added to the NoteGraph
func (ng *NoteGraph) GetNotes() *[]Note {
	return &ng.Notes
}

// GetLinkedNotes returns all linked Notes of the NoteGraph
func (ng *NoteGraph) GetLinkedNotes() *[]Note {
	linkedNotes := []Note{}
	for _, note := range ng.Notes {
		if ng.LinkedNoteGUIDs[note.GUID] {
			linkedNotes = append(linkedNotes, note)
		}
	}

	return &linkedNotes
}

// GetNoteLinks returns all NoteLinks added to the NoteGraph
func (ng *NoteGraph) GetNoteLinks() *[]NoteLink {
	return &ng.NoteLinks
}
