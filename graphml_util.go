package main

import (
	"encoding/xml"
	"os"

	"github.com/freddy33/graphml"
	"github.com/sirupsen/logrus"
)

// NodeLabelID is the ID of the GraphML attribute used for the label of nodes in the graph
const NodeLabelID = "node-label"

// NodeLabelName is the name of the GraphML attribute used for the label of nodes in the graph
const NodeLabelName = "label"

// NodeDescriptionID is the ID of the GraphML attribute used for the description of nodes in the graph
const NodeDescriptionID = "node-description"

// NodeDescriptionName is the name of the GraphML attribute used for the description of nodes in the graph
const NodeDescriptionName = "description"

// NodeURLID is the ID of the GraphML attribute used for the URL of nodes in the graph
const NodeURLID = "node-url"

// NodeURLName is the name of the GraphML attribute used for the URL of nodes in the graph
const NodeURLName = "url"

// EdgeLabelID is the ID of the GraphML attribute used for the label of edges in the graph
const EdgeLabelID = "edge-label"

// EdgeLabelName is the name of the GraphML attribute used for the label of edges in the graph
const EdgeLabelName = "label"

// EdgeDescriptionID is the ID of the GraphML attribute used for the description of edges in the graph
const EdgeDescriptionID = "edge-description"

// EdgeDescriptionName is the name of the GraphML attribute used for the description of edges in the graph
const EdgeDescriptionName = "description"

// GraphMLUtil provides a number of util methods to create GraphML documents with standardized node and edge data
type GraphMLUtil struct{}

// NewGraphMLUtil creates a new instance of GraphMLUtil
func NewGraphMLUtil() *GraphMLUtil {
	return &GraphMLUtil{}
}

// SaveGraphMLDocument saves the provided graphMLDocument with the specified filename on the file system
func (gu *GraphMLUtil) SaveGraphMLDocument(filename string, graphMLDocument *graphml.Document) error {
	logrus.Infof("Saving GraphML to file [%s]", filename)

	file, fileErr := os.Create(filename)
	defer file.Close()
	if fileErr != nil {
		return fileErr
	}

	encodeErr := graphml.Encode(file, graphMLDocument)
	if encodeErr != nil {
		return encodeErr
	}

	return nil
}

// CreateGraphMLDocument creates the GraphML document based on the supplied GraphML graphs with the standardised GraphML attributes definition
func (gu *GraphMLUtil) CreateGraphMLDocument(graphs []graphml.Graph) *graphml.Document {
	return &graphml.Document{
		Instr: xml.ProcInst{
			Target: "xml",
			Inst:   []byte("version=\"1.0\" encoding=\"UTF-8\"")},
		Attrs: []xml.Attr{
			{Name: xml.Name{Local: "xmlns"}, Value: "http://graphml.graphdrawing.org/xmlns"},
			{Name: xml.Name{Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
			{Name: xml.Name{Local: "xsi:schemaLocation"}, Value: "http://graphml.graphdrawing.org/xmlns http://www.yworks.com/xml/schema/graphml/1.1/ygraphml.xsd"}},
		Graphs: graphs,
		Keys: []graphml.Key{
			graphml.NewKey(graphml.KindNode, NodeLabelID, NodeLabelName, "string"),
			graphml.NewKey(graphml.KindNode, NodeDescriptionID, NodeDescriptionName, "string"),
			graphml.NewKey(graphml.KindNode, NodeURLID, NodeURLName, "string"),
			graphml.NewKey(graphml.KindEdge, EdgeLabelID, EdgeLabelName, "string"),
			graphml.NewKey(graphml.KindEdge, EdgeDescriptionID, EdgeDescriptionName, "string")}}
}

// CreateGraph creates a GraphML graph with the specified id, nodes, edges, and edge direction
func (gu *GraphMLUtil) CreateGraph(id string, edgeDefault graphml.EdgeDir, nodes []graphml.Node, edges []graphml.Edge) *graphml.Graph {
	return &graphml.Graph{
		ExtObject:   graphml.ExtObject{Object: graphml.Object{ID: id}},
		EdgeDefault: edgeDefault,
		Nodes:       nodes,
		Edges:       edges}
}

// CreateNode creates a GraphML node with the specified id and the supplied label, description, and URL as GraphML attributes
func (gu *GraphMLUtil) CreateNode(nodeID, nodeLabel, nodeDescription, nodeURL string) *graphml.Node {
	return &graphml.Node{
		ExtObject: graphml.ExtObject{
			Object: graphml.Object{ID: nodeID},
			Data: []graphml.Data{
				graphml.NewData(NodeLabelID, nodeLabel),
				graphml.NewData(NodeDescriptionID, nodeDescription),
				graphml.NewData(NodeURLID, nodeURL)}}}
}

// CreateEdge creates a GraphML edge with the specified id, source and target nodes, and the supplied label and description as GraphML attributes
func (gu *GraphMLUtil) CreateEdge(edgeID, sourceNodeID, targetNodeID, edgeLabel, edgeDescription string) *graphml.Edge {
	return &graphml.Edge{
		ExtObject: graphml.ExtObject{
			Object: graphml.Object{ID: edgeID},
			Data: []graphml.Data{
				graphml.NewData(EdgeLabelID, edgeLabel),
				graphml.NewData(EdgeDescriptionID, edgeDescription)}},
		Source: sourceNodeID,
		Target: targetNodeID}
}
