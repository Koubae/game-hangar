Game Hangar
============

![GameHangar](docs/img/game-hangar-pic.png "GameHangar")

todo


### QuickStart


* 1) Install [air-verse/air](https://github.com/air-verse/air) globally

```bash
go install github.com/air-verse/air@latest

# Make sure that GOPATH and GOROOT is in your PATH
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin
```
* 2) Initialize (this will perform some installations)

```bash
make init
```

* 3) Start

```bash
# Run chat-identity
make run-server-local
```
