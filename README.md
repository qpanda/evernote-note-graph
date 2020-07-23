EvernoteNoteGraph
================
A small Go program that generates a graph of all linked notes a user has created in [Evernote](https://evernote.com/).

**EvernoteNoteGraph** uses the [Evernote API](http://dev.evernote.com/doc/) to retrieve all notes and generates a note graph in [GraphML](https://en.wikipedia.org/wiki/GraphML) format based on the [Note Links](https://dev.evernote.com/doc/articles/note_links.php) detected across all notes. Technically speaking the generated note graph is a [directed](https://en.wikipedia.org/wiki/Directed_graph) [disconnected](https://en.wikipedia.org/wiki/Connectivity_%28graph_theory%29) [multigraph](https://en.wikipedia.org/wiki/Multigraph).

The note graph in GraphML format can then be loaded into a graph editor such as [yEd](https://www.yworks.com/products/yed), [Gephi](https://gephi.org/), or [Cytoscape](http://www.cytoscape.org/) for layouting and analysis / exploration.

## Inspiration
**EvernoteNoteGraph** is inspired by the way [Roam](https://roamresearch.com/) and [Obsidian](https://obsidian.md/) visually display linked notes as a graph.

## Prerequisites
In order to use **EvernoteNoteGraph** you have to request a [DeveloperToken](https://www.evernote.com/api/DeveloperToken.action) for your Evernote account. You may need to contact Evernote support to get the feature enabled for your account.

## Installation
TODO

## Using EvernoteTagCloud
Run ```evernote-note-graph -h``` to get usage information. All parameters except for -edamAuthToken (Evernote Developer Token / API Key) are optional.

    $ evernote-note-graph -h
    Usage of evernote-note-graph:
    -edamAuthToken string        
            Evernote API auth token
    -graphMLFilename string
            GraphML output filename (default "noteGraph.graphml")
    -linkedNotes
            Include only linked Notes (default true)
    -noteURL string
            WebLink or AppLink for Note URLs (default "WebLink")
    -sandbox
            Use sandbox.evernote.com
    -v    Verbose output

TODO

## Examples
TODO

## License
**EvernoteNoteGraph** is licensed under the MIT license.