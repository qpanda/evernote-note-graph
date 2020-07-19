package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
)

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

// NoteLink is an app, web, public, or shortened link that points to an Evernote note (see Evernote API documentation at https://dev.evernote.com/doc/articles/note_links.php)
type NoteLink struct {
	Type     LinkType
	URL      url.URL
	Text     string
	NoteGUID string // is not set for ShortenedLinks
}

func (nl NoteLink) String() string {
	return fmt.Sprintf("{Type: %s, URL: %s, Text: %s, NoteGUID %s}", nl.Type.String(), nl.URL.String(), nl.Text, nl.NoteGUID)
}

// NoteLinkParser can be used to create and parse NoteLinks
type NoteLinkParser struct {
	EvernoteHost string
	UserID       string
	ShardID      string
}

// NewNoteLinkParser creates a new instance of NoteLinkParser
func NewNoteLinkParser(evernoteHost, userID, shardID string) *NoteLinkParser {
	return &NoteLinkParser{
		EvernoteHost: evernoteHost,
		UserID:       userID,
		ShardID:      shardID}
}

// ExtractNoteLinks extracts all NoteLinks detected / found in the supplied note content (ENML)
func (elp *NoteLinkParser) ExtractNoteLinks(enml string) ([]NoteLink, error) {
	enmlDocument, err := htmlquery.Parse(strings.NewReader(enml))
	if err != nil {
		return nil, err
	}

	noteLinks := []NoteLink{}
	htmlLinks := htmlquery.Find(enmlDocument, "//a")
	for _, a := range htmlLinks {
		linkURL, err := url.Parse(htmlquery.SelectAttr(a, "href"))
		linkText := htmlquery.InnerText(a)
		if err != nil {
			return nil, err
		}

		noteLink := elp.ParseNoteLink(*linkURL, linkText)
		if noteLink != nil {
			noteLinks = append(noteLinks, *noteLink)
		}
	}

	return noteLinks, nil
}

// ParseNoteLink parses the supplied URL and returns a NoteLink if the URL points to an Evernote note, otherwise returns nil
// For AppLinks and WebLinks method ParseNoteLink verifies that the link is for the user and shard provided when creating the NoteLinkParser
// (AppLinks and WebLinks for other users are not accessible) and if this is not the case returns nil instead of the AppLink / WebLink
func (elp *NoteLinkParser) ParseNoteLink(linkURL url.URL, linkText string) *NoteLink {
	trimmedPath := strings.TrimRight(linkURL.Path, "/")
	pathElements := strings.Split(trimmedPath, "/")

	if linkURL.Scheme == "evernote" {
		if len(pathElements) == 6 && pathElements[1] == "view" && pathElements[2] == elp.UserID && pathElements[3] == elp.ShardID && pathElements[4] == pathElements[5] {
			// evernote:///view/[userId]/[shardId]/[noteGuid]/[noteGuid]/
			return &NoteLink{Type: AppLink, URL: linkURL, Text: linkText, NoteGUID: pathElements[4]}
		}
	}

	if linkURL.Scheme == "https" && linkURL.Hostname() == elp.EvernoteHost {
		if len(pathElements) == 3 && pathElements[1] == "l" {
			// https://[evernoteHost]/l/[random string]/
			return &NoteLink{Type: ShortenedLink, URL: linkURL, Text: linkText}
		} else if len(pathElements) == 6 && pathElements[1] == "shard" && pathElements[2] == elp.ShardID && pathElements[3] == "nl" && pathElements[4] == elp.UserID {
			// https://[evernoteHost]/shard/[shardId]/nl/[userId]/[noteGuid]/
			return &NoteLink{Type: WebLink, URL: linkURL, Text: linkText, NoteGUID: pathElements[5]}
		} else if len(pathElements) == 6 && pathElements[1] == "shard" && pathElements[3] == "sh" {
			// https://[evernoteHost]/shard/[shardId]/sh/[noteGuid]/[shareKey]/
			return &NoteLink{Type: PublicLink, URL: linkURL, Text: linkText, NoteGUID: pathElements[4]}
		}
	}

	return nil
}

// CreateAppLinkURL creates a NoteLink of type AppLink that points to the note with the provided GUID
func (elp *NoteLinkParser) CreateAppLinkURL(noteGUID string) (*url.URL, error) {
	// evernote:///view/[userId]/[shardId]/[noteGuid]/[noteGuid]/
	return url.Parse(fmt.Sprintf("evernote:///view/%s/%s/%s/%s/", elp.UserID, elp.ShardID, noteGUID, noteGUID))
}

// CreateShortenedLinkURL creates a NoteLink of type ShortenedLink
func (elp *NoteLinkParser) CreateShortenedLinkURL(random string) (*url.URL, error) {
	// https://[evernoteHost]/l/[random string]
	return url.Parse(fmt.Sprintf("https://%s/l/%s", elp.EvernoteHost, random))
}

// CreateWebLinkURL creates a NoteLink of type WebLink that points to the note with the provided GUID
func (elp *NoteLinkParser) CreateWebLinkURL(noteGUID string) (*url.URL, error) {
	// https://[evernoteHost]/shard/[shardId]/nl/[userId]/[noteGuid]/
	return url.Parse(fmt.Sprintf("https://%s/shard/%s/nl/%s/%s/", elp.EvernoteHost, elp.ShardID, elp.UserID, noteGUID))
}

// CreatePublicLinkURL creates a NoteLink of type PublicLink
func (elp *NoteLinkParser) CreatePublicLinkURL(noteGUID, shareKey string) (*url.URL, error) {
	// https://[evernoteHost]/shard/[shardId]/sh/[noteGuid]/[shareKey]/
	return url.Parse(fmt.Sprintf("https://%s/shard/%s/sh/%s/%s/", elp.EvernoteHost, elp.ShardID, noteGUID, shareKey))
}
