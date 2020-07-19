EvernoteNoteGraph
================
A small Go program that generates a graph of all linked notes a user has created in [Evernote](https://evernote.com/).

**EvernoteNoteGraph** uses the [Evernote API](http://dev.evernote.com/doc/) to retrieve all notes and generate a note graph in [GraphML](https://en.wikipedia.org/wiki/GraphML) based on the [Note Links](https://dev.evernote.com/doc/articles/note_links.php) detected across all notes.