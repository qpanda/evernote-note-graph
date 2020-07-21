package main

import (
	"github.com/freddy33/graphml"
	uuid "github.com/satori/go.uuid"
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

// ConvertNoteGraph converts the NoteGraph into a GraphML document
func (ngu *NoteGraphUtil) ConvertNoteGraph(noteGraph *NoteGraph, allNotes bool) *graphml.Document {
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
