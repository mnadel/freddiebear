# bearfred

An Alfred Workflow for [Bear](https://bear.app) inspired by [alfred-bear](https://github.com/chrisbro/alfred-bear).

Implemented in Golang for faster searching, with daily note journaling capabilities.

# Searching

To search for a note by title, use the `bt` keyword.

To search a note by its full contents, use the `bs` keyword.

# Journaling

The `captainslog` keyword will either open today's log (a note with the title `YYYY-MM-DD`) or create a new note with that title.

Run `bearfred help journal` for details on how to tweak the tag it attaches to new notes.

# Implementation

This Golang implementation is pretty snappy on my current 4MB database. Most of the performance gains of this implementaion over a Python implementation appears to be reduced startup cost. That said, [db.go](https://github.com/mnadel/bearfred/blob/main/db/db.go) includes some SQLite3 pragmas that will hopefully keep it snappy as it grows. Sample timing that returns about half the records in the database:

```
bearfred search --all drip  0.00s user 0.00s system 73% cpu 0.010 total
```
