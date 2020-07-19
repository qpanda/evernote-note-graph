package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
	"github.com/freddy33/graphml"
	"github.com/stretchr/testify/assert"
)

var graphMLUtil = GraphMLUtil{}

const (
	NodeAID          = "A"
	NodeALabel       = "Node A Label"
	NodeADescription = "Node A Description"
	NodeAURL         = "https://example.org/nodeA"

	NodeBID          = "B"
	NodeBLabel       = "Node B Label"
	NodeBDescription = "Node B Description"
	NodeBURL         = "https://example.org/nodeB"

	NodeCID          = "C"
	NodeCLabel       = "Node C Label"
	NodeCDescription = "Node C Description"
	NodeCURL         = "https://example.org/nodeC"

	NodeDID          = "D"
	NodeDLabel       = "Node D Label"
	NodeDDescription = "Node D Description"
	NodeDURL         = "https://example.org/nodeD"

	EdgeABID           = "AB"
	EdgeABSourceNodeID = "A"
	EdgeABTargetNodeID = "B"
	EdgeABLabel        = "A->B"
	EdgeABDescription  = "A->B"

	EdgeBC1ID           = "BC1"
	EdgeBC1SourceNodeID = "B"
	EdgeBC1TargetNodeID = "C"
	EdgeBC1Label        = "B->C"
	EdgeBC1Description  = "B->C(weight=1)"

	EdgeBC2ID           = "BC2"
	EdgeBC2SourceNodeID = "B"
	EdgeBC2TargetNodeID = "C"
	EdgeBC2Label        = "B->C"
	EdgeBC2Description  = "B->C(weight=2)"
)

func TestLogGraphMLDocument(t *testing.T) {
	graphMLDocument := CreateTestGraphMLDocument()
	encodedGraphMLDocument := EncodeGraphMLDocument(graphMLDocument)
	t.Log(encodedGraphMLDocument)
}

func TestEncodeDecodeGraphMLDocument(t *testing.T) {
	graphMLDocument := CreateTestGraphMLDocument()
	encodedGraphMLDocument := EncodeGraphMLDocument(graphMLDocument)
	decodedGraphMLDocument := DecodeGraphMLDocument(encodedGraphMLDocument)

	assert.Equal(t, graphMLDocument.Keys, decodedGraphMLDocument.Keys, "keys should be equal")
	assert.Equal(t, graphMLDocument.Graphs, decodedGraphMLDocument.Graphs, "graphs should be equal")
}

func TestCreateGraphMLDocument(t *testing.T) {
	graphMLDocument := CreateTestGraphMLDocument()
	encodedGraphMLDocument := EncodeGraphMLDocument(graphMLDocument)
	xmlDocument, err := xmlquery.Parse(strings.NewReader(encodedGraphMLDocument))
	if err != nil {
		panic(err)
	}

	graph := xmlquery.FindOne(xmlDocument, "/graphml/graph")
	assert.Equal(t, "TestGraph", graph.SelectAttr("id"))
	assert.Equal(t, string(graphml.EdgeDirected), graph.SelectAttr("edgedefault"))

	AssertKeyEqual(t, xmlDocument, NodeLabelID, "node", NodeLabelName, "string")
	AssertKeyEqual(t, xmlDocument, NodeDescriptionID, "node", NodeDescriptionName, "string")
	AssertKeyEqual(t, xmlDocument, NodeURLID, "node", NodeURLName, "string")
	AssertKeyEqual(t, xmlDocument, EdgeLabelID, "edge", EdgeLabelName, "string")
	AssertKeyEqual(t, xmlDocument, EdgeDescriptionID, "edge", EdgeDescriptionName, "string")

	AssertNodeEqual(t, xmlDocument, NodeAID, NodeALabel, NodeADescription, NodeAURL)
	AssertNodeEqual(t, xmlDocument, NodeBID, NodeBLabel, NodeBDescription, NodeBURL)
	AssertNodeEqual(t, xmlDocument, NodeCID, NodeCLabel, NodeCDescription, NodeCURL)
	AssertNodeEqual(t, xmlDocument, NodeDID, NodeDLabel, NodeDDescription, NodeDURL)

	AssertEdgeEqual(t, xmlDocument, EdgeABID, EdgeABSourceNodeID, EdgeABTargetNodeID, EdgeABLabel, EdgeABDescription)
	AssertEdgeEqual(t, xmlDocument, EdgeBC1ID, EdgeBC1SourceNodeID, EdgeBC1TargetNodeID, EdgeBC1Label, EdgeBC1Description)
	AssertEdgeEqual(t, xmlDocument, EdgeBC2ID, EdgeBC2SourceNodeID, EdgeBC2TargetNodeID, EdgeBC2Label, EdgeBC2Description)
}

func TestSaveGraphMLDocument(t *testing.T) {
	testGraphFile := filepath.Join(os.TempDir(), "testGraph.graphml")
	defer os.Remove(testGraphFile)

	graphMLDocument := CreateTestGraphMLDocument()
	graphMLUtil.SaveGraphMLDocument(testGraphFile, graphMLDocument)
}

func CreateTestGraphMLDocument() *graphml.Document {
	nodeA := graphMLUtil.CreateNode(NodeAID, NodeALabel, NodeADescription, NodeAURL)
	nodeB := graphMLUtil.CreateNode(NodeBID, NodeBLabel, NodeBDescription, NodeBURL)
	nodeC := graphMLUtil.CreateNode(NodeCID, NodeCLabel, NodeCDescription, NodeCURL)
	nodeD := graphMLUtil.CreateNode(NodeDID, NodeDLabel, NodeDDescription, NodeDURL)
	nodes := []graphml.Node{*nodeA, *nodeB, *nodeC, *nodeD}

	edgeAB := graphMLUtil.CreateEdge(EdgeABID, EdgeABSourceNodeID, EdgeABTargetNodeID, EdgeABLabel, EdgeABDescription)
	edgeBC1 := graphMLUtil.CreateEdge(EdgeBC1ID, EdgeBC1SourceNodeID, EdgeBC1TargetNodeID, EdgeBC1Label, EdgeBC1Description)
	edgeBC2 := graphMLUtil.CreateEdge(EdgeBC2ID, EdgeBC2SourceNodeID, EdgeBC2TargetNodeID, EdgeBC2Label, EdgeBC2Description)
	edges := []graphml.Edge{*edgeAB, *edgeBC1, *edgeBC2}

	graph := graphMLUtil.CreateGraph("TestGraph", graphml.EdgeDirected, nodes, edges)
	graphs := []graphml.Graph{*graph}

	return graphMLUtil.CreateGraphMLDocument(graphs)
}

func EncodeGraphMLDocument(graphMLDocument *graphml.Document) string {
	encodedGraphMLDocument := bytes.Buffer{}
	err := graphml.Encode(&encodedGraphMLDocument, graphMLDocument)
	if err != nil {
		panic(err)
	}

	return encodedGraphMLDocument.String()
}

func DecodeGraphMLDocument(graphMLDocument string) *graphml.Document {
	decodedGraphMLDocument, err := graphml.Decode(strings.NewReader(graphMLDocument))
	if err != nil {
		panic(err)
	}

	return decodedGraphMLDocument
}

func AssertKeyEqual(t *testing.T, xmlNode *xmlquery.Node, keyID, keyFor, keyAttrName, keyAttrType string) {
	key := xmlquery.FindOne(xmlNode, "/graphml/key[@id='"+keyID+"']")
	assert.Equal(t, keyFor, key.SelectAttr("for"))
	assert.Equal(t, keyAttrName, key.SelectAttr("attr.name"))
	assert.Equal(t, keyAttrType, key.SelectAttr("attr.type"))
}

func AssertNodeEqual(t *testing.T, xmlNode *xmlquery.Node, nodeID, nodeLabel, nodeDescription, nodeURL string) {
	node := xmlquery.FindOne(xmlNode, "/graphml/graph/node[@id='"+nodeID+"']")
	assert.Equal(t, nodeLabel, xmlquery.FindOne(node, "/data[@key='"+NodeLabelID+"']").InnerText())
	assert.Equal(t, nodeDescription, xmlquery.FindOne(node, "/data[@key='"+NodeDescriptionID+"']").InnerText())
	assert.Equal(t, nodeURL, xmlquery.FindOne(node, "/data[@key='"+NodeURLID+"']").InnerText())
}

func AssertEdgeEqual(t *testing.T, xmlNode *xmlquery.Node, edgeID, sourceNodeID, targetNodeID, edgeLabel, edgeDescription string) {
	edge := xmlquery.FindOne(xmlNode, "/graphml/graph/edge[@id='"+edgeID+"']")
	assert.Equal(t, sourceNodeID, edge.SelectAttr("source"))
	assert.Equal(t, targetNodeID, edge.SelectAttr("target"))
	assert.Equal(t, edgeLabel, xmlquery.FindOne(edge, "/data[@key='"+EdgeLabelID+"']").InnerText())
	assert.Equal(t, edgeDescription, xmlquery.FindOne(edge, "/data[@key='"+EdgeDescriptionID+"']").InnerText())
}
