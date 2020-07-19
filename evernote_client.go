package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/dreampuf/evernote-sdk-golang/edam"
	"github.com/sirupsen/logrus"
)

var (
	no  = false
	yes = true
)

// Timeout specifies the maximum time calls to the Evernote API are allowed to take
const Timeout = time.Duration(30) * time.Second

// NoteSortOrder defines the order in which note metadata is fetched
var NoteSortOrder = int32(edam.NoteSortOrder_CREATED)

// SandboxEvernoteCom is the sandbox Evernote API endpoint URL
const SandboxEvernoteCom = "sandbox.evernote.com"

// EvernoteCom is the Evernote API endpoint URL
const EvernoteCom = "www.evernote.com"

// EvernoteClient is a thin wrapper around the Evernote SDK
type EvernoteClient struct {
	AuthToken       string
	Sandbox         bool
	UserStoreClient *edam.UserStoreClient
	NoteStoreClient *edam.NoteStoreClient
}

// IEvernoteClient is an interface that exposes all EvernoteClient functions required to contstruct a NoteGraph
type IEvernoteClient interface {
	GetHost() string
	GetUserStoreURL() string
	FindAllNotesMetadata(offset int32, maxNotes int32) (*edam.NotesMetadataList, error)
	GetNoteWithContent(guid edam.GUID) (*edam.Note, error)
}

// NewEvernoteClient creates a new instance of EvernoteClient
func NewEvernoteClient(authToken string, sandbox bool) *EvernoteClient {
	return &EvernoteClient{
		AuthToken: authToken,
		Sandbox:   sandbox}
}

// GetHost returns the hostname of the Evernote API used by EvernoteClient
func (ec *EvernoteClient) GetHost() string {
	if ec.Sandbox {
		return SandboxEvernoteCom
	}

	return EvernoteCom
}

// GetUserStoreURL returns the URL of the Evernote UserStore API
func (ec *EvernoteClient) GetUserStoreURL() string {
	return fmt.Sprintf("https://%s/edam/user", ec.GetHost())
}

// GetUserStoreClient returns the Evernote UserStoreClient
func (ec *EvernoteClient) GetUserStoreClient() (*edam.UserStoreClient, error) {
	if ec.UserStoreClient != nil {
		return ec.UserStoreClient, nil
	}

	thriftTransport, err := thrift.NewTHttpClient(ec.GetUserStoreURL())
	if err != nil {
		return nil, err
	}

	thriftClient := thrift.NewTStandardClient(thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(thriftTransport), thrift.NewTBinaryProtocolFactory(true, true).GetProtocol(thriftTransport))
	ec.UserStoreClient = edam.NewUserStoreClient(thriftClient)
	return ec.UserStoreClient, nil
}

// GetUser returns the Evernote User
func (ec *EvernoteClient) GetUser() (*edam.User, error) {
	userStoreClient, err := ec.GetUserStoreClient()
	if err != nil {
		return nil, err
	}

	context, cancelFunc := context.WithTimeout(context.Background(), Timeout)
	defer func() {
		logrus.Tracef("GetUser context completed")
		cancelFunc()
	}()

	logrus.Debugf("Retrieving user information from Evernote API endpoint [%s]", ec.GetHost())
	user, err := userStoreClient.GetUser(context, ec.AuthToken)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetNoteStoreClient returns the Evernote NoteStoreClient
func (ec *EvernoteClient) GetNoteStoreClient() (*edam.NoteStoreClient, error) {
	if ec.NoteStoreClient != nil {
		return ec.NoteStoreClient, nil
	}

	userStoreClient, err := ec.GetUserStoreClient()
	if err != nil {
		return nil, err
	}

	context, cancelFunc := context.WithTimeout(context.Background(), Timeout)
	defer func() {
		logrus.Tracef("GetNoteStoreClient context completed")
		cancelFunc()
	}()

	logrus.Tracef("Retrieving user URLs from Evernote API endpoint [%s]", ec.GetHost())
	userUrls, err := userStoreClient.GetUserUrls(context, ec.AuthToken)
	if err != nil {
		return nil, err
	}

	thriftTransport, err := thrift.NewTHttpClient(userUrls.GetNoteStoreUrl())
	if err != nil {
		return nil, err
	}

	thriftClient := thrift.NewTStandardClient(thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(thriftTransport), thrift.NewTBinaryProtocolFactory(true, true).GetProtocol(thriftTransport))
	return edam.NewNoteStoreClient(thriftClient), nil
}

// FindAllNotesMetadata returns the metadata of up to maxNotes notes from the specified offset
// Returns the metadata for all notes (does not apply a note filter) including note title and note attributes in descending NoteSortOrder
func (ec *EvernoteClient) FindAllNotesMetadata(offset int32, maxNotes int32) (*edam.NotesMetadataList, error) {
	noteStoreClient, err := ec.GetNoteStoreClient()
	if err != nil {
		return nil, err
	}

	context, cancelFunc := context.WithTimeout(context.Background(), Timeout)
	defer func() {
		logrus.Tracef("FindAllNotesMetadata context completed")
		cancelFunc()
	}()

	filter := &edam.NoteFilter{Order: &NoteSortOrder, Ascending: &no}
	resultSpec := &edam.NotesMetadataResultSpec{IncludeTitle: &yes, IncludeAttributes: &yes}
	logrus.Debugf("Retrieving metadata for notes from offset [%d] with page size [%d] from Evernote API endpoint [%s]", offset, maxNotes, ec.GetHost())
	return noteStoreClient.FindNotesMetadata(context, ec.AuthToken, filter, offset, maxNotes, resultSpec)
}

// GetNoteWithContent returns the note specified by the GUID including the note content (ENML)
func (ec *EvernoteClient) GetNoteWithContent(guid edam.GUID) (*edam.Note, error) {
	noteStoreClient, err := ec.GetNoteStoreClient()
	if err != nil {
		return nil, err
	}

	context, cancelFunc := context.WithTimeout(context.Background(), Timeout)
	defer func() {
		logrus.Tracef("GetNoteWithContent context completed")
		cancelFunc()
	}()

	resultSpec := &edam.NoteResultSpec{IncludeContent: &yes}
	logrus.Debugf("Retrieving note with GUID [%s] from Evernote API endpoint [%s]", guid, ec.GetHost())
	return noteStoreClient.GetNoteWithResultSpec(context, ec.AuthToken, guid, resultSpec)
}
