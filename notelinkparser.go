package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
)

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
func (elp *NoteLinkParser) ExtractNoteLinks(noteGUID, noteContent string) ([]NoteLink, error) {
	enmlDocument, err := htmlquery.Parse(strings.NewReader(noteContent))
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

		noteLink := elp.ParseNoteLink(noteGUID, *linkURL, linkText)
		if noteLink != nil {
			noteLinks = append(noteLinks, *noteLink)
		}
	}

	return noteLinks, nil
}

// ParseNoteLink parses the supplied URL and returns a NoteLink if the URL points to an Evernote note, otherwise returns nil
// For AppLinks and WebLinks method ParseNoteLink verifies that the link is for the user and shard provided when creating the NoteLinkParser
// (AppLinks and WebLinks for other users are not accessible) and if this is not the case returns nil instead of the AppLink / WebLink
func (elp *NoteLinkParser) ParseNoteLink(noteGUID string, linkURL url.URL, linkText string) *NoteLink {
	trimmedPath := strings.TrimRight(linkURL.Path, "/")
	pathElements := strings.Split(trimmedPath, "/")

	if linkURL.Scheme == "evernote" {
		if len(pathElements) == 6 && pathElements[1] == "view" && pathElements[2] == elp.UserID && pathElements[3] == elp.ShardID && pathElements[4] == pathElements[5] {
			// evernote:///view/[userId]/[shardId]/[noteGuid]/[noteGuid]/
			return &NoteLink{SourceNoteGUID: noteGUID, TargetNoteGUID: pathElements[4], Text: linkText, URL: linkURL, URLType: AppLink}
		}
	}

	if linkURL.Scheme == "https" && linkURL.Hostname() == elp.EvernoteHost {
		if len(pathElements) == 3 && pathElements[1] == "l" {
			// https://[evernoteHost]/l/[random string]/
			return &NoteLink{SourceNoteGUID: noteGUID, Text: linkText, URL: linkURL, URLType: ShortenedLink}
		} else if len(pathElements) == 6 && pathElements[1] == "shard" && pathElements[2] == elp.ShardID && pathElements[3] == "nl" && pathElements[4] == elp.UserID {
			// https://[evernoteHost]/shard/[shardId]/nl/[userId]/[noteGuid]/
			return &NoteLink{SourceNoteGUID: noteGUID, TargetNoteGUID: pathElements[5], Text: linkText, URL: linkURL, URLType: WebLink}
		} else if len(pathElements) == 6 && pathElements[1] == "shard" && pathElements[3] == "sh" {
			// https://[evernoteHost]/shard/[shardId]/sh/[noteGuid]/[shareKey]/
			return &NoteLink{SourceNoteGUID: noteGUID, TargetNoteGUID: pathElements[4], Text: linkText, URL: linkURL, URLType: PublicLink}
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
