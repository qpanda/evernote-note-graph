package main

import (
	"fmt"

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
	fmt.Printf("NoteGraph Stats\n")
	fmt.Printf("   Notes: %d\n", len(*noteGraph.GetNotes()))
	fmt.Printf("   Linked Notes: %d\n", len(*noteGraph.GetLinkedNotes()))
	fmt.Printf("   Note Links: %d\n", len(*noteGraph.GetNoteLinks()))
	fmt.Printf("   Valid Note Links: %d\n", len(*noteGraph.GetValidNoteLinks()))
	fmt.Printf("   Broken Note Links: %d\n", len(*noteGraph.GetBrokenNoteLinks()))
}

// PrintBrokenNoteLinks prints all broken NoteLinks
func (ngu *NoteGraphUtil) PrintBrokenNoteLinks(noteGraph *NoteGraph) {
	brokenNoteLinks := *noteGraph.GetBrokenNoteLinks()
	if len(brokenNoteLinks) > 0 {
		fmt.Printf("Broken Note Links\n")
		for _, noteLink := range brokenNoteLinks {
			sourceNote := noteGraph.GetNote(noteLink.SourceNoteGUID)
			targetNote := noteGraph.GetNote(noteLink.SourceNoteGUID)
			fmt.Printf("   NoteLink [%v] from source Note [%v] to target Note [%v]\n", brokenNoteLinks, sourceNote, targetNote)
		}
	}
}

// ConvertNoteGraph converts the NoteGraph into a GraphML document
func (ngu *NoteGraphUtil) ConvertNoteGraph(noteGraph *NoteGraph, allNotes bool) *graphml.Document {
	logrus.Infof("Converting NoteGraph to GraphML")

	notes := ngu.GraphNotes(noteGraph, allNotes)
	noteLinks := ngu.GraphNoteLinks(noteGraph)

	nodes := ngu.CreateNodes(notes)
	edges := ngu.CreateEdges(noteLinks)

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
		node := ngu.GraphMLUtil.CreateNode(note.GUID, note.Title, note.Description, note.URL.String())
		nodes = append(nodes, *node)
	}

	return nodes
}

// CreateEdges creates a GraphML edge from the NoteLink
func (ngu *NoteGraphUtil) CreateEdges(noteLinks []NoteLink) []graphml.Edge {
	edges := []graphml.Edge{}
	for _, noteLink := range noteLinks {
		edge := ngu.GraphMLUtil.CreateEdge(uuid.NewV4().String(), noteLink.SourceNoteGUID, noteLink.TargetNoteGUID, noteLink.Text, noteLink.Text)
		edges = append(edges, *edge)
	}

	return edges
}
