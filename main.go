package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// ParseArgs parses command line arguments
func ParseArgs() (string, bool, URLType, bool, string, bool) {
	edamAuthToken := flag.String("edamAuthToken", "", "Evernote API auth token")
	sandbox := flag.Bool("sandbox", false, "Use sandbox.evernote.com")
	noteURL := flag.String("noteURL", "WebLink", "WebLink or AppLink for Note URLs")
	linkedNotes := flag.Bool("linkedNotes", true, "Include only linked Notes")
	graphMLFilename := flag.String("graphMLFilename", "notegraph.graphml", "GraphML output filename")
	verbose := flag.Bool("v", false, "Verbose output")

	flag.Parse()

	if *edamAuthToken == "" {
		flag.Usage()
		os.Exit(2)
	}

	noteURLType, err := NewURLType(*noteURL)
	if err != nil || (*noteURLType != WebLink && *noteURLType != AppLink) {
		flag.Usage()
		os.Exit(2)
	}

	return *edamAuthToken, *sandbox, *noteURLType, *linkedNotes, *graphMLFilename, *verbose
}

// PlainFormatter is a simple logrus Formatter
type PlainFormatter struct{}

// Format print log level and message only
func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\t - %s\n", strings.ToUpper(entry.Level.String()), entry.Message)), nil
}

// InitLogger initializes the Logrus logger
func InitLogger(verbose bool) {
	logrus.SetFormatter(&PlainFormatter{})
	logrus.SetOutput(os.Stdout)

	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

// InitEvernoteClient initializes the EvernoteClient
func InitEvernoteClient(edamAuthToken string, sandbox bool) IEvernoteClient {
	return NewEvernoteClient(edamAuthToken, sandbox)
}

// InitNoteLinkParser initializes the NoteLinkParser
func InitNoteLinkParser(evernoteClient IEvernoteClient) *NoteLinkParser {
	user, err := evernoteClient.GetUser()
	if err != nil {
		logrus.Errorf("Failed to retrieve user from Evernote API at [%s]: %v", evernoteClient.GetHost(), err)
		panic(err)
	}

	evernoteHost := evernoteClient.GetHost()
	userID := fmt.Sprint(user.GetID())
	shardID := user.GetShardId()

	logrus.Infof("Using Evernote API endpoint at [%s] with user [%s]", evernoteClient.GetHost(), user.GetUsername())
	return NewNoteLinkParser(evernoteHost, userID, shardID)
}

// InitEvernoteNoteGraph initializes the EvernoteNoteGraph
func InitEvernoteNoteGraph(edamAuthToken string, sandbox bool, noteURLType URLType) *EvernoteNoteGraph {
	evernoteClient := InitEvernoteClient(edamAuthToken, sandbox)
	noteLinkParser := InitNoteLinkParser(evernoteClient)
	return NewEvernoteNoteGraph(evernoteClient, noteLinkParser, noteURLType)
}

// CreateNoteGraph creates the NoteGraph from Evernote notes
func CreateNoteGraph(evernoteNoteGraph *EvernoteNoteGraph) *NoteGraph {
	noteGraph, noteGraphErr := evernoteNoteGraph.CreateNoteGraph()
	if noteGraphErr != nil {
		logrus.Errorf("Failed to create NoteGraph from Evernote API at [%s]: %v", evernoteNoteGraph.EvernoteClient.GetHost(), noteGraphErr)
		panic(noteGraphErr)
	}

	return noteGraph
}

// SaveNoteGraph saves the NoteGraph as GraphML
func SaveNoteGraph(noteGraph *NoteGraph, linkedNotes bool, graphMLFilename string) {
	graphMLDocument := NewNoteGraphUtil().ConvertNoteGraph(noteGraph, !linkedNotes)
	saveGraphMLErr := NewGraphMLUtil().SaveGraphMLDocument(graphMLFilename, graphMLDocument)
	if saveGraphMLErr != nil {
		logrus.Errorf("Failed to save NoteGraph to GraphML file [%s]: %v", graphMLFilename, saveGraphMLErr)
		panic(saveGraphMLErr)
	}
}

func main() {
	edamAuthToken, sandbox, noteURLType, linkedNotes, graphMLFilename, verbose := ParseArgs()

	InitLogger(verbose)

	evernoteNoteGraph := InitEvernoteNoteGraph(edamAuthToken, sandbox, noteURLType)
	noteGraph := CreateNoteGraph(evernoteNoteGraph)
	SaveNoteGraph(noteGraph, linkedNotes, graphMLFilename)

	NewNoteGraphUtil().PrintNoteGraphStats(noteGraph)
	NewNoteGraphUtil().PrintBrokenNoteLinks(noteGraph)
}
