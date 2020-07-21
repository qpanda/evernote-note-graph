package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// InitLogger initializes the Logrus logger
func InitLogger() {
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	logrus.SetFormatter(Formatter)
	logrus.SetLevel(logrus.DebugLevel)
}

// ParseArgs parses command line arguments
func ParseArgs() (string, bool, URLType, bool, string) {
	edamAuthToken := flag.String("edamAuthToken", "", "Evernote API auth token")
	sandbox := flag.Bool("sandbox", false, "Use sandbox.evernote.com")
	noteURL := flag.String("noteURL", "WebLink", "WebLink or AppLink for Note URLs")
	linkedNotes := flag.Bool("linkedNotes", true, "Include only linked Notes")
	graphMLFilename := flag.String("graphMLFilename", "noteGraph.graphml", "GraphML output filename")
	flag.Parse()

	if *edamAuthToken == "" {
		flag.Usage()
		os.Exit(2)
	}

	urlType, err := NewURLType(*noteURL)
	if err != nil || (*urlType != WebLink && *urlType != AppLink) {
		flag.Usage()
		os.Exit(2)
	}

	return *edamAuthToken, *sandbox, *urlType, *linkedNotes, *graphMLFilename
}

// InitEvernoteClient initializes the EvernoteClient
func InitEvernoteClient(edamAuthToken string) IEvernoteClient {
	return NewEvernoteClient(edamAuthToken, true)
}

// InitNoteLinkParser initializes the NoteLinkParser
func InitNoteLinkParser(evernoteClient IEvernoteClient) *NoteLinkParser {
	user, err := evernoteClient.GetUser()
	if err != nil {
		panic(err)
	}

	evernoteHost := evernoteClient.GetHost()
	userID := fmt.Sprint(user.GetID())
	shardID := user.GetShardId()
	return NewNoteLinkParser(evernoteHost, userID, shardID)
}

// InitEvernoteNoteGraph initializes the EvernoteNoteGraph
func InitEvernoteNoteGraph(edamAuthToken string, noteURLType URLType) *EvernoteNoteGraph {
	evernoteClient := InitEvernoteClient(edamAuthToken)
	noteLinkParser := InitNoteLinkParser(evernoteClient)
	return NewEvernoteNoteGraph(evernoteClient, noteLinkParser, noteURLType)
}

// CreateEvernoteNoteGraph creates and saves the NoteGraph as a GraphML document
func CreateEvernoteNoteGraph(edamAuthToken string, sandbox bool, noteURLType URLType, linkedNotes bool, graphMLFilename string) {
	evernoteNoteGraph := InitEvernoteNoteGraph(edamAuthToken, noteURLType)
	noteGraph, noteGraphErr := evernoteNoteGraph.CreateNoteGraph()
	if noteGraphErr != nil {
		panic(noteGraphErr)
	}

	graphMLDocument := NewNoteGraphUtil().ConvertNoteGraph(noteGraph, !linkedNotes)
	saveGraphMLErr := NewGraphMLUtil().SaveGraphMLDocument(graphMLFilename, graphMLDocument)
	if saveGraphMLErr != nil {
		panic(saveGraphMLErr)
	}

	NewNoteGraphUtil().PrintNoteGraphStats(noteGraph)
	NewNoteGraphUtil().PrintBrokenNoteLinks(noteGraph)
}

func main() {
	InitLogger()
	CreateEvernoteNoteGraph(ParseArgs())
}
