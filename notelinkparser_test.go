package main

import (
	"testing"

	"net/url"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

const (
	Host    = "www.evernote.com"
	UserID  = "76136038"
	ShardID = "s12"
)

var noteLinkParser = NewNoteLinkParser(Host, UserID, ShardID)

var testENML = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd">
<en-note>
	<div><a href="https://example.org/">NonNoteLink</a></div>
	<div><a href="https://www.evernote.com/shard/s12/nl/76136038/d72dfad0-7d58-41b5-b2c9-4ca434abd543/">WebLink</a></div>
	<div><a href="evernote:///view/76136038/s12/4d971333-8b65-45d6-857b-243c850cabf5/4d971333-8b65-45d6-857b-243c850cabf5/">AppLink</a></div>
	<div><a href="https://www.evernote.com/shard/s12/sh/4d971333-8b65-45d6-857b-243c850cabf5/25771cdb535e9183/">PublicLink</a></div>
	<div><a href="https://www.evernote.com/l/AAxNlxMzi2VF1oV7JDyFDKv1JXcc21NekYM">ShortenedLink</a></div>
</en-note>`

func TestCreateWebLinkURL(t *testing.T) {
	noteGUID := uuid.NewV4().String()
	webLinkURL := CreateWebLinkURL(noteGUID)
	assert.Equal(t, "https://"+Host+"/shard/"+ShardID+"/nl/"+UserID+"/"+noteGUID+"/", webLinkURL.String())
}

func TestCreateShortenedLinkURL(t *testing.T) {
	random := uuid.NewV4().String()
	shortenedLinkURL := CreateShortenedLinkURL(random)
	assert.Equal(t, "https://"+Host+"/l/"+random, shortenedLinkURL.String())
}

func TestCreateAppLinkURL(t *testing.T) {
	noteGUID := uuid.NewV4().String()
	appLinkURL := CreateAppLinkURL(noteGUID)
	assert.Equal(t, "evernote:///view/"+UserID+"/"+ShardID+"/"+noteGUID+"/"+noteGUID+"/", appLinkURL.String())
}

func TestCreatePublicLinkURL(t *testing.T) {
	noteGUID := uuid.NewV4().String()
	shareKey := uuid.NewV4().String()
	publicLinkURL := CreatePublicLinkURL(noteGUID, shareKey)
	assert.Equal(t, "https://"+Host+"/shard/"+ShardID+"/sh/"+noteGUID+"/"+shareKey+"/", publicLinkURL.String())
}

func TestParseNonNoteLink(t *testing.T) {
	noteGUID := uuid.NewV4().String()
	nonNoteLink := noteLinkParser.ParseNoteLink(noteGUID, *CreateURL("https://example.org/"), "example.org")
	assert.Nil(t, nonNoteLink)
}

func TestParseWebNoteLink(t *testing.T) {
	sourceNoteGUID := uuid.NewV4().String()
	targetNoteGUID := uuid.NewV4().String()
	webLinkURL := CreateWebLinkURL(targetNoteGUID)
	webLinkText := "WebLink"
	webLink := noteLinkParser.ParseNoteLink(sourceNoteGUID, *webLinkURL, webLinkText)
	assert.Equal(t, WebLink, webLink.URLType)
	assert.Equal(t, sourceNoteGUID, webLink.SourceNoteGUID)
	assert.Equal(t, targetNoteGUID, webLink.TargetNoteGUID)
	assert.Equal(t, *webLinkURL, webLink.URL)
	assert.Equal(t, webLinkText, webLink.Text)
}

func TestParseShortenedNoteLink(t *testing.T) {
	sourceNoteGUID := uuid.NewV4().String()
	random := uuid.NewV4().String()
	shortenedLinkURL := CreateShortenedLinkURL(random)
	shortenedLinkText := "ShortenedLink"
	shortenedLink := noteLinkParser.ParseNoteLink(sourceNoteGUID, *shortenedLinkURL, shortenedLinkText)
	assert.Equal(t, ShortenedLink, shortenedLink.URLType)
	assert.Equal(t, sourceNoteGUID, shortenedLink.SourceNoteGUID)
	assert.Empty(t, shortenedLink.TargetNoteGUID)
	assert.Equal(t, *shortenedLinkURL, shortenedLink.URL)
	assert.Equal(t, shortenedLinkText, shortenedLink.Text)
}

func TestParseAppNoteLink(t *testing.T) {
	sourceNoteGUID := uuid.NewV4().String()
	targetNoteGUID := uuid.NewV4().String()
	appLinkURL := CreateAppLinkURL(targetNoteGUID)
	appLinkText := "AppLink"
	appLink := noteLinkParser.ParseNoteLink(sourceNoteGUID, *appLinkURL, appLinkText)
	assert.Equal(t, AppLink, appLink.URLType)
	assert.Equal(t, sourceNoteGUID, appLink.SourceNoteGUID)
	assert.Equal(t, targetNoteGUID, appLink.TargetNoteGUID)
	assert.Equal(t, *appLinkURL, appLink.URL)
	assert.Equal(t, appLinkText, appLink.Text)
}

func TestParsePublicNoteLink(t *testing.T) {
	sourceNoteGUID := uuid.NewV4().String()
	targetNoteGUID := uuid.NewV4().String()
	shareKey := uuid.NewV4().String()
	appLinkURL := CreatePublicLinkURL(targetNoteGUID, shareKey)
	appLinkText := "PublicLink"
	appLink := noteLinkParser.ParseNoteLink(sourceNoteGUID, *appLinkURL, appLinkText)
	assert.Equal(t, PublicLink, appLink.URLType)
	assert.Equal(t, sourceNoteGUID, appLink.SourceNoteGUID)
	assert.Equal(t, targetNoteGUID, appLink.TargetNoteGUID)
	assert.Equal(t, *appLinkURL, appLink.URL)
	assert.Equal(t, appLinkText, appLink.Text)
}

func TestExtractNoteLinks(t *testing.T) {
	noteGUID := uuid.NewV4().String()
	noteLinks, err := noteLinkParser.ExtractNoteLinks(noteGUID, testENML)
	if err != nil {
		panic(err)
	}

	assert.Len(t, noteLinks, 4)
}

func CreateURL(link string) *url.URL {
	url, err := url.Parse(link)
	if err != nil {
		panic(err)
	}

	return url
}

func CreateWebLinkURL(noteGUID string) *url.URL {
	url, err := noteLinkParser.CreateWebLinkURL(noteGUID)
	if err != nil {
		panic(err)
	}

	return url
}

func CreateShortenedLinkURL(random string) *url.URL {
	url, err := noteLinkParser.CreateShortenedLinkURL(random)
	if err != nil {
		panic(err)
	}

	return url
}

func CreateAppLinkURL(noteGUID string) *url.URL {
	url, err := noteLinkParser.CreateAppLinkURL(noteGUID)
	if err != nil {
		panic(err)
	}

	return url
}

func CreatePublicLinkURL(noteGUID, shareKey string) *url.URL {
	url, err := noteLinkParser.CreatePublicLinkURL(noteGUID, shareKey)
	if err != nil {
		panic(err)
	}

	return url
}
