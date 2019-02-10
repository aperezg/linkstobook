# Links to Book
It's a sample tool with unique motive to use [go-docker](https://docker.io/go-docker)

I created this tool for learn how to work [go-docker](https://docker.io/go-docker) library and I used the [pandoc docker image](https://github.com/jagregory/pandoc-docker) to convert web pages into an epub.

I highly recommend that you use better the pandoc tool instead of my library for this purpose.

## Install

```sh
$ go get github.com/aperezg/linkstobook
```

## How to use

```
A simple tool to convert web pages into epub

Usage:
  linkstobook [command]

Available Commands:
  convert     Convert output into a epub file
  help        Help about any command

Flags:
  -h, --help   help for linkstobook

Use "linkstobook [command] --help" for more information about a command.
```

ex
```sh
$ linkstobook convert --web https://blog.friendsofgo.tech/posts/crear-tu-primer-cli-en-go/,https://blog.friendsofgo.tech/posts/dockerizando-tu-aplicacion-en-go/
```
