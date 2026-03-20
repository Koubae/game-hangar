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
make quickstart
```


### Services

* **Identity:** http://localhost:10000
* **GameHangar:** http://localhost:10010


#### Step by Step 


* Prepare Go env

```bash
make init
```

* Run Migrations

```bash
make migrate-identity-up
``` 


* Start

```bash
# Run chat-identity
make run-server-local
```

### Documentation

* [Working with PostgreSQL in Go using pgx](https://donchev.is/post/working-with-postgresql-in-go-using-pgx/)
  * [Reddit post](https://www.reddit.com/r/golang/comments/1c8br5c/does_anyone_have_a_clear_example_of_how_to_use/)


