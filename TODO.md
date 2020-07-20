# Features
* Introduce command line arguments to control graph generation
* Add LinkType parameter to EvernoteNoteGraph to make it possible to choose AppLink and WebLink for links in NoteGraph
* Add IncludeUnlinkedNotes parameter to NoteGraphUtil to control whether all or only linked Notes are included in the NoteGraph
* Add IncludeLinkTypes parameter to EvernoteNoteGraph to choose which LinkTypes to consider for links in the NoteGraph

# Improvements
* Review usage of pass-by-value vs pass-by-reference across the codebase for efficiency purpose
* Consider surrounding note and edge labels and descriptsion with CDATA in GraphML output
* Review codebase for missing parameter validation
* Consider CI
* Consider security Github Actions
* Tag/Publish version
* Consider fetching notes content in parallel

# Documentation
* Add inspiration section to README
* Add Prerequisites section to README
* Add Installation section to README
* Add Usage section to README (including yEd, Cytoscape, Gephi)
* Add Example section to README
* Add License section to README
