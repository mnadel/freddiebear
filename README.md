# freddiebear

An Alfred Workflow for [Bear](https://bear.app) inspired by [alfred-bear](https://github.com/chrisbro/alfred-bear).

Implemented in Golang for faster searching, with daily note journaling capabilities.

# Searching

To search for a note by title, use the `bt` keyword.

To search a note by its full contents, use the `bs` keyword.

# Journaling

The `captainslog` keyword will either open today's log (a note with the title `YYYY-MM-DD`) or create a new note with that title.

Run `freddiebear help journal` for details on how to tweak the tag it attaches to new notes.

# Exporting

You can `export` the text contents of your notes to Markdown files. Specify the directory and we'll create files in the form of `<title> (<sha>).md`.

Titles aren't unique, so we append a unique ID for each note. It also allows us to track renamed notes if you re-export into an existing directory.

# Implementation

This Golang implementation is pretty snappy on my current 4MB database. Most of the performance gains of this implementaion over a Python implementation appears to be reduced startup cost. That said, [db.go](https://github.com/mnadel/freddiebear/blob/main/db/db.go) includes some SQLite3 pragmas that will hopefully keep it snappy as it grows. Sample timing that returns about half the records in the database:

```
freddiebear search --all drip  0.00s user 0.00s system 66% cpu 0.012 total
```

## --show-tags

Show tags will generate a list longest-path tags to show as an Alfred item's subtitle.

For example, a note with a tags `q` and `a/b/c` will have four tags in the database:
1. `q`
1. `a`
1. `a/b`
1. `a/b/c`

And we'll only return the terminal/non-intermediate tags (`a/b/c` and `q` in this example).

The current implementation uses a O(n^2) algorithm, but in practice is quite fast for small sets of tags. Compare its implementation to one that uses a prefix trie:

```
→ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/mnadel/freddiebear/util
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkRemoveIntermediatePrefixes-16   4701206               233.7 ns/op
BenchmarkPrefixTrie-16                    611623              1865 ns/op
```
