package main

import (
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
	"github.com/freddy33/graphml"
	"github.com/stretchr/testify/assert"
)

func TestCreateEdges(t *testing.T) {
	noteLinkParser = NewNoteLinkParser("www.evernote.com", "76136038", "s12")
	webNoteLink := noteLinkParser.ParseNoteLink("sourceNoteGUID", *CreateWebLinkURL("targetNoteGUID"), "WebLink")
	edges := NewNoteGraphUtil().CreateEdges([]NoteLink{*webNoteLink})

	assert.Equal(t, webNoteLink.SourceNoteGUID, edges[0].Source)
	assert.Equal(t, webNoteLink.TargetNoteGUID, edges[0].Target)
}

func TestCreateNodes(t *testing.T) {
	webLinkURL := CreateWebLinkURL("GUID")
	note := Note{GUID: "GUID", Title: "Title", Description: "Title", URL: *webLinkURL, URLType: WebLink}
	nodes := NewNoteGraphUtil().CreateNodes([]Note{note})

	assert.Equal(t, note.GUID, nodes[0].ID)
}

func TestConvertNoteGraph(t *testing.T) {
	// four Notes, valid NoteLinks, disconnected graph
	noteA := Note{GUID: "A", Title: "TitleA", Description: "DescriptionA", URL: *CreateWebLinkURL("A"), URLType: WebLink}
	noteB := Note{GUID: "B", Title: "TitleB", Description: "DescriptionB", URL: *CreateWebLinkURL("B"), URLType: WebLink}
	noteC := Note{GUID: "C", Title: "TitleC", Description: "DescriptionC", URL: *CreateWebLinkURL("C"), URLType: WebLink}
	noteD := Note{GUID: "D", Title: "TitleD", Description: "DescriptionD", URL: *CreateWebLinkURL("D"), URLType: WebLink}

	noteLinkAB := NoteLink{SourceNoteGUID: "A", TargetNoteGUID: "B", URL: *CreateWebLinkURL("B"), URLType: WebLink}
	noteLinkCD := NoteLink{SourceNoteGUID: "C", TargetNoteGUID: "D", URL: *CreateWebLinkURL("D"), URLType: WebLink}

	noteGraph := NewNoteGraph()
	noteGraph.Add(noteA, []NoteLink{noteLinkAB})
	noteGraph.Add(noteB, []NoteLink{})
	noteGraph.Add(noteC, []NoteLink{noteLinkCD})
	noteGraph.Add(noteD, []NoteLink{})

	graphMLDocument := NewNoteGraphUtil().ConvertNoteGraph(noteGraph)
	encodedGraphMLDocument := EncodeGraphMLDocument(graphMLDocument)
	xmlDocument, err := xmlquery.Parse(strings.NewReader(encodedGraphMLDocument))
	if err != nil {
		panic(err)
	}

	graph := xmlquery.FindOne(xmlDocument, "/graphml/graph")
	assert.Equal(t, NoteGraphID, graph.SelectAttr("id"))
	assert.Equal(t, string(graphml.EdgeDirected), graph.SelectAttr("edgedefault"))

	AssertNoteEqualNode(t, xmlDocument, noteA.GUID, noteA.Title, noteA.Description, noteA.URL.String())
	AssertNoteEqualNode(t, xmlDocument, noteB.GUID, noteB.Title, noteB.Description, noteB.URL.String())
	AssertNoteEqualNode(t, xmlDocument, noteC.GUID, noteC.Title, noteC.Description, noteC.URL.String())
	AssertNoteEqualNode(t, xmlDocument, noteD.GUID, noteD.Title, noteD.Description, noteD.URL.String())

	AssertNoteLinkEqualEdge(t, xmlDocument, "1", noteLinkAB.SourceNoteGUID, noteLinkAB.TargetNoteGUID, noteLinkAB.Text, noteLinkAB.Text)
	AssertNoteLinkEqualEdge(t, xmlDocument, "2", noteLinkCD.SourceNoteGUID, noteLinkCD.TargetNoteGUID, noteLinkCD.Text, noteLinkCD.Text)
}

func AssertNoteEqualNode(t *testing.T, xmlNode *xmlquery.Node, nodeID, nodeLabel, nodeDescription, nodeURL string) {
	node := xmlquery.FindOne(xmlNode, "/graphml/graph/node[@id='"+nodeID+"']")
	assert.Equal(t, nodeLabel, xmlquery.FindOne(node, "/data[@key='"+NodeLabelID+"']").InnerText())
	assert.Equal(t, nodeDescription, xmlquery.FindOne(node, "/data[@key='"+NodeDescriptionID+"']").InnerText())
	assert.Equal(t, nodeURL, xmlquery.FindOne(node, "/data[@key='"+NodeURLID+"']").InnerText())
}

func AssertNoteLinkEqualEdge(t *testing.T, xmlNode *xmlquery.Node, edgeIndex, sourceNodeID, targetNodeID, edgeLabel, edgeDescription string) {
	edge := xmlquery.FindOne(xmlNode, "/graphml/graph/edge["+edgeIndex+"]")
	assert.Equal(t, sourceNodeID, edge.SelectAttr("source"))
	assert.Equal(t, targetNodeID, edge.SelectAttr("target"))
	assert.Equal(t, edgeLabel, xmlquery.FindOne(edge, "/data[@key='"+EdgeLabelID+"']").InnerText())
	assert.Equal(t, edgeDescription, xmlquery.FindOne(edge, "/data[@key='"+EdgeDescriptionID+"']").InnerText())
}
