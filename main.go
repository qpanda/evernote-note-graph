package main

import (
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

// InitEvernoteClient initializes the EvernoteClient
func InitEvernoteClient() IEvernoteClient {
	if os.Getenv("EDAM_AUTHTOKEN") == "" {
		panic("EDAM_AUTHTOKEN env variable is not set")
	}

	return NewEvernoteClient(os.Getenv("EDAM_AUTHTOKEN"), true)
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
func InitEvernoteNoteGraph() *EvernoteNoteGraph {
	evernoteClient := InitEvernoteClient()
	noteLinkParser := InitNoteLinkParser(evernoteClient)
	return NewEvernoteNoteGraph(evernoteClient, noteLinkParser, WebLink)
}

// CreateEvernoteNoteGraph creates and saves the NoteGraph as a GraphML document
func CreateEvernoteNoteGraph() {
	evernoteNoteGraph := InitEvernoteNoteGraph()
	noteGraph, noteGraphErr := evernoteNoteGraph.CreateNoteGraph()
	if noteGraphErr != nil {
		panic(noteGraphErr)
	}

	graphMLDocument := NewNoteGraphUtil().ConvertNoteGraph(noteGraph, false)
	saveGraphMLErr := NewGraphMLUtil().SaveGraphMLDocument("noteGraph.graphml", graphMLDocument)
	if saveGraphMLErr != nil {
		panic(saveGraphMLErr)
	}

	NewNoteGraphUtil().PrintNoteGraphStats(noteGraph)
	NewNoteGraphUtil().PrintBrokenNoteLinks(noteGraph)
}

func main() {
	InitLogger()
	CreateEvernoteNoteGraph()
}
