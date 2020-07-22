package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/dreampuf/evernote-sdk-golang/edam"
	"github.com/shafreeck/retry"
	"github.com/sirupsen/logrus"
)

var (
	no  = false
	yes = true
)

// Timeout specifies the maximum time calls to the Evernote API are allowed to take
const Timeout = time.Duration(5) * time.Second

// Retries specifies how many times remote calls should be tried before returning an error
const Retries = 3

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
	GetUser() (*edam.User, error)
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
		return nil, fmt.Errorf("Failed to create Thrift HttpClient with UserStoreURL [%v]: %w", ec.GetUserStoreURL(), err)
	}

	thriftClient := thrift.NewTStandardClient(thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(thriftTransport), thrift.NewTBinaryProtocolFactory(true, true).GetProtocol(thriftTransport))
	ec.UserStoreClient = edam.NewUserStoreClient(thriftClient)
	return ec.UserStoreClient, nil
}

// GetUser returns the Evernote User
func (ec *EvernoteClient) GetUser() (*edam.User, error) {
	userStoreClient, err := ec.GetUserStoreClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to create UserStoreClient: %w", err)
	}

	retriable := retry.New()
	context, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	user := &edam.User{}
	retriableErr := retriable.EnsureN(context, Retries, func() error {
		logrus.Debugf("Retrieving user information from Evernote API endpoint [%s]", ec.GetHost())
		user, err = userStoreClient.GetUser(context, ec.AuthToken)
		if err != nil {
			logrus.Warnf("Retrieving user information from Evernote API endpoint [%s] failed with error [%s] - retrying [%d] times", ec.GetHost(), err, Retries)
			return retry.Retriable(err)
		}

		return nil
	})

	if retriableErr != nil {
		return nil, fmt.Errorf("Failed to retrieve user information from Evernote API endpoint [%s] after [%d] retries: %w", ec.GetHost(), Retries, retriableErr)
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
		return nil, fmt.Errorf("Failed to create UserStoreClient: %w", err)
	}

	retriable := retry.New()
	context, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	userUrls := &edam.UserUrls{}
	retriableErr := retriable.EnsureN(context, Retries, func() error {
		logrus.Tracef("Retrieving user URLs from Evernote API endpoint [%s]", ec.GetHost())
		userUrls, err = userStoreClient.GetUserUrls(context, ec.AuthToken)
		if err != nil {
			logrus.Warnf("Retrieving user URLs from Evernote API endpoint [%s] failed with error [%s] - retrying [%d] times", ec.GetHost(), err, Retries)
			return retry.Retriable(err)
		}

		return nil
	})

	if retriableErr != nil {
		return nil, fmt.Errorf("Failed to retrieve user URLs from Evernote API endpoint [%s] after [%d] retries: %w", ec.GetHost(), Retries, retriableErr)
	}

	thriftTransport, err := thrift.NewTHttpClient(userUrls.GetNoteStoreUrl())
	if err != nil {
		return nil, fmt.Errorf("Failed to create Thrift HttpClient with NoteStoreURL [%v]: %w", userUrls.GetNoteStoreUrl(), err)
	}

	thriftClient := thrift.NewTStandardClient(thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(thriftTransport), thrift.NewTBinaryProtocolFactory(true, true).GetProtocol(thriftTransport))
	return edam.NewNoteStoreClient(thriftClient), nil
}

// FindAllNotesMetadata returns the metadata of up to maxNotes notes from the specified offset
// Returns the metadata for all notes (does not apply a note filter) including note title and note attributes in descending NoteSortOrder
func (ec *EvernoteClient) FindAllNotesMetadata(offset int32, maxNotes int32) (*edam.NotesMetadataList, error) {
	noteStoreClient, err := ec.GetNoteStoreClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to create NoteStoreClient: %w", err)
	}

	retriable := retry.New()
	context, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	filter := &edam.NoteFilter{Order: &NoteSortOrder, Ascending: &no}
	resultSpec := &edam.NotesMetadataResultSpec{IncludeTitle: &yes, IncludeAttributes: &yes}

	notesMetadataList := &edam.NotesMetadataList{}
	retriableErr := retriable.EnsureN(context, Retries, func() error {
		logrus.Debugf("Retrieving metadata for notes from offset [%d] with page size [%d] from Evernote API endpoint [%s]", offset, maxNotes, ec.GetHost())
		notesMetadataList, err = noteStoreClient.FindNotesMetadata(context, ec.AuthToken, filter, offset, maxNotes, resultSpec)
		if err != nil {
			logrus.Warnf("Retrieving metadata for notes from offset [%d] with page size [%d] from Evernote API endpoint [%s] failed with error [%s] - retrying [%d] times", offset, maxNotes, ec.GetHost(), err, Retries)
			return retry.Retriable(err)
		}

		return nil
	})

	if retriableErr != nil {
		return nil, fmt.Errorf("Failed to retrieve metadata for notes from offset [%d] with page size [%d] from Evernote API endpoint [%s] after [%d] retries: %w", offset, maxNotes, ec.GetHost(), Retries, retriableErr)
	}

	return notesMetadataList, nil
}

// GetNoteWithContent returns the note specified by the GUID including the note content (ENML)
func (ec *EvernoteClient) GetNoteWithContent(guid edam.GUID) (*edam.Note, error) {
	noteStoreClient, err := ec.GetNoteStoreClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to create NoteStoreClient: %w", err)
	}

	retriable := retry.New()
	context, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	resultSpec := &edam.NoteResultSpec{IncludeContent: &yes}

	note := &edam.Note{}
	retriableErr := retriable.EnsureN(context, Retries, func() error {
		logrus.Debugf("Retrieving note with GUID [%s] from Evernote API endpoint [%s]", guid, ec.GetHost())
		note, err = noteStoreClient.GetNoteWithResultSpec(context, ec.AuthToken, guid, resultSpec)
		if err != nil {
			logrus.Warnf("Retrieving note with GUID [%s] from Evernote API endpoint [%s] failed with error [%s] - retrying [%d] times", guid, ec.GetHost(), err, Retries)
			return retry.Retriable(err)
		}

		return nil
	})

	if retriableErr != nil {
		return nil, fmt.Errorf("Failed to retrieve note with GUID [%s] from Evernote API endpoint [%s] after [%d] retries: %w", guid, ec.GetHost(), Retries, retriableErr)
	}

	return note, nil
}
