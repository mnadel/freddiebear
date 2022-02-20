# bearfred

An Alfred Workflow for [Bear](https://bear.app) inspired by [alfred-bear](https://github.com/chrisbro/alfred-bear).

Implemented in Golang for faster searching, with daily note journaling capabilities.

# Searching

To search for a note by title, use the `bt` keyword.

To search a note by its full contents, use the `bs` keyword.

# Journaling

The `captainslog` keyword will either open today's log (a note with the title `YYYY-MM-DD`) or create a new note with that title.

Run `bearfred help journal` for details on how to tweak the tag it attaches to new notes.
