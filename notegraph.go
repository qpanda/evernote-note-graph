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

// NoteLink is an app, web, public, or shortened link that points from source Note to target Note (see Evernote API documentation at https://dev.evernote.com/doc/articles/note_links.php)
type NoteLink struct {
	Type           LinkType
	URL            url.URL
	Text           string
	SourceNoteGUID string
	TargetNoteGUID string // not set for ShortenedLinks
}

func (nl NoteLink) String() string {
	return fmt.Sprintf("{Type: %s, URL: %s, Text: %s, SourceNoteGUID: %s, TargetNoteGUID: %s}", nl.Type.String(), nl.URL.String(), nl.Text, nl.SourceNoteGUID, nl.TargetNoteGUID)
}

// NoteGraph contains all Notes and NoteLinks and keeps track of which Notes are linked to other Notes
type NoteGraph struct {
	Notes     map[string]Note // all Notes
	NoteLinks []NoteLink      // all NoteLinks
}

// NewNoteGraph creates a new instance of NoteGraph
func NewNoteGraph() *NoteGraph {
	return &NoteGraph{
		Notes:     map[string]Note{},
		NoteLinks: []NoteLink{},
	}
}

// Add adds a Note with all its NoteLinks to the NoteGraph, returns true if the Note is a linked note, otherwise false
func (ng *NoteGraph) Add(note Note, noteLinks []NoteLink) bool {
	ng.Notes[note.GUID] = note
	ng.NoteLinks = append(ng.NoteLinks, noteLinks...)
	return len(noteLinks) != 0
}

// GetNotes returns all Notes added to the NoteGraph
func (ng *NoteGraph) GetNotes() *[]Note {
	notes := []Note{}
	for _, note := range ng.Notes {
		notes = append(notes, note)
	}

	return &notes
}

// GetNoteLinks returns all NoteLinks added to the NoteGraph
func (ng *NoteGraph) GetNoteLinks() *[]NoteLink {
	return &ng.NoteLinks
}

// GetLinkedNotes returns source and target Notes of all NoteLinks where both source and target Note exists
func (ng *NoteGraph) GetLinkedNotes() *[]Note {
	linkedNotes := map[string]Note{}
	for _, noteLink := range ng.NoteLinks {
		sourceNote, sourceNoteFound := ng.Notes[noteLink.SourceNoteGUID]
		targetNote, targetNoteFound := ng.Notes[noteLink.TargetNoteGUID]

		if sourceNoteFound && targetNoteFound {
			linkedNotes[sourceNote.GUID] = sourceNote
			linkedNotes[targetNote.GUID] = targetNote
		}
	}

	notes := []Note{}
	for _, note := range linkedNotes {
		notes = append(notes, note)
	}

	return &notes
}

// GetValidNoteLinks returns all valid NoteLinks (both source and target Note exist)
func (ng *NoteGraph) GetValidNoteLinks() *[]NoteLink {
	validNoteLinks := []NoteLink{}
	for _, noteLink := range ng.NoteLinks {
		_, sourceNoteFound := ng.Notes[noteLink.SourceNoteGUID]
		_, targetNoteFound := ng.Notes[noteLink.TargetNoteGUID]

		if sourceNoteFound && targetNoteFound {
			validNoteLinks = append(validNoteLinks, noteLink)
		}
	}

	return &validNoteLinks
}

// GetBrokenNoteLinks returns all broken NoteLinks, either source, or target, or both notes missing
func (ng *NoteGraph) GetBrokenNoteLinks() *[]NoteLink {
	brokenNoteLinks := []NoteLink{}
	for _, noteLink := range ng.NoteLinks {
		_, sourceNoteFound := ng.Notes[noteLink.SourceNoteGUID]
		_, targetNoteFound := ng.Notes[noteLink.TargetNoteGUID]

		if !(sourceNoteFound && targetNoteFound) {
			brokenNoteLinks = append(brokenNoteLinks, noteLink)
		}
	}

	return &brokenNoteLinks
}
