# on

Command line interface for [fsnotify](https://github.com/fsnotify/fsnotify).

It watches the files and writes the event right into _stdout_. Pretty useful to complex bash pipe's manipulation.

### Usage

Install: `go get github.com/olebedev/on` or compiled binaries at [releases](https://github.com/olebedev/on/releases) page.

```
$ on --help
NAME:
   on - cli for fsnotify

USAGE:
   on [options] [path]

   Path could be mix of folders and filepaths, default is '.'.
   Regular usecase is watch the file's changes and execute some
   bash script or command line tool. It could be done in this way:

   on | xargs -n1 -I{} <program>

VERSION:
   0.1.0

AUTHOR(S):
   olebedev <ole6edev@gmail.com>

COMMANDS:
GLOBAL OPTIONS:
   --template, -t "{{.Name}}" output template to render received event, see: https://godoc.org/gopkg.in/fsnotify.v1#Event
   --mask, -m "15"    event's bitwise mask, see: https://godoc.org/gopkg.in/fsnotify.v1#Op
   --help, -h     show help
   --version, -v    print the version
```

### License
MIT
