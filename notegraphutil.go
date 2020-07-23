package main

import (
	"strings"

	"github.com/freddy33/graphml"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// NoteGraphUtil converts a NoteGraph to GraphML and saves the the GraphML document to a file
type NoteGraphUtil struct {
	GraphMLUtil GraphMLUtil
}

// NoteGraphID is the ID used for the note graph of the GraphML document
const NoteGraphID = "NoteGraph"

// NewNoteGraphUtil creates a new instance of NoteGraphUtil
func NewNoteGraphUtil() *NoteGraphUtil {
	return &NoteGraphUtil{GraphMLUtil{}}
}

// PrintNoteGraphStats prints NoteGraph stats to stdout
func (ngu *NoteGraphUtil) PrintNoteGraphStats(noteGraph *NoteGraph) {
	logrus.Infof("NoteGraph Stats")
	logrus.Infof("   Notes: %d", len(*noteGraph.GetNotes()))
	logrus.Infof("   Linked Notes: %d", len(*noteGraph.GetLinkedNotes()))
	logrus.Infof("   Note Links: %d", len(*noteGraph.GetNoteLinks()))
	logrus.Infof("   Valid Note Links: %d", len(*noteGraph.GetValidNoteLinks()))
	logrus.Infof("   Broken Note Links: %d", len(*noteGraph.GetBrokenNoteLinks()))
}

// PrintBrokenNoteLinks prints all broken NoteLinks
func (ngu *NoteGraphUtil) PrintBrokenNoteLinks(noteGraph *NoteGraph) {
	brokenNoteLinks := *noteGraph.GetBrokenNoteLinks()
	if len(brokenNoteLinks) > 0 {
		logrus.Infof("Broken Note Links")
		for _, noteLink := range brokenNoteLinks {
			sourceNote := noteGraph.GetNote(noteLink.SourceNoteGUID)
			targetNote := noteGraph.GetNote(noteLink.SourceNoteGUID)
			logrus.Infof("   NoteLink [%v] from source Note [%v] to target Note [%v]", noteLink, sourceNote, targetNote)
		}
	}
}

// ConvertNoteGraph converts the NoteGraph into a GraphML document
func (ngu *NoteGraphUtil) ConvertNoteGraph(noteGraph *NoteGraph, allNotes bool) *graphml.Document {
	notes := ngu.GraphNotes(noteGraph, allNotes)
	noteLinks := ngu.GraphNoteLinks(noteGraph)

	nodes := ngu.CreateNodes(notes)
	edges := ngu.CreateEdges(noteLinks)

	logrus.Infof("Converting NoteGraph with [%d|%d] Notes|nodes and [%d|%d] NoteLinks|edges to GraphML", len(notes), len(nodes), len(noteLinks), len(edges))

	graph := ngu.GraphMLUtil.CreateGraph(NoteGraphID, graphml.EdgeDirected, nodes, edges)
	return ngu.GraphMLUtil.CreateGraphMLDocument([]graphml.Graph{*graph})
}

// GraphNotes returns all Notes to include in the GraphML graph
func (ngu *NoteGraphUtil) GraphNotes(noteGraph *NoteGraph, allNotes bool) []Note {
	if allNotes {
		return *noteGraph.GetNotes()
	}

	return *noteGraph.GetLinkedNotes()
}

// GraphNoteLinks returns all NoteLinks to include in the GraphML graph
func (ngu *NoteGraphUtil) GraphNoteLinks(noteGraph *NoteGraph) []NoteLink {
	return *noteGraph.GetValidNoteLinks()
}

// CreateNodes creates a GraphML node from the Note
func (ngu *NoteGraphUtil) CreateNodes(notes []Note) []graphml.Node {
	nodes := []graphml.Node{}
	for _, note := range notes {
		node := ngu.GraphMLUtil.CreateNode(note.GUID, strings.ReplaceAll(note.Title, " ", "‧"), strings.ReplaceAll(note.Description, " ", "‧"), note.URL.String())
		nodes = append(nodes, *node)
	}

	return nodes
}

// CreateEdges creates a GraphML edge from the NoteLink
func (ngu *NoteGraphUtil) CreateEdges(noteLinks []NoteLink) []graphml.Edge {
	edges := []graphml.Edge{}
	for _, noteLink := range noteLinks {
		edge := ngu.GraphMLUtil.CreateEdge(uuid.NewV4().String(), noteLink.SourceNoteGUID, noteLink.TargetNoteGUID, strings.ReplaceAll(noteLink.Text, " ", "‧"), strings.ReplaceAll(noteLink.Text, " ", "‧"))
		edges = append(edges, *edge)
	}

	return edges
}
