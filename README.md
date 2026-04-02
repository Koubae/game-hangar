Game Hangar
============

![GameHangar](docs/img/game-hangar-pic.png "GameHangar")

todo

### Documentation

* [GameHangar -- Documentation](docs/README.MD)
  * [Identity (Account & Authorization)](docs/identity/README.MD)



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

### Admin / Management

* Root User: 
  * Username: `root`
  * Password: `Q7m!v2Zp#9tX`

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

#### Database 

* [Working with PostgreSQL in Go using pgx](https://donchev.is/post/working-with-postgresql-in-go-using-pgx/)
  * [Reddit post](https://www.reddit.com/r/golang/comments/1c8br5c/does_anyone_have_a_clear_example_of_how_to_use/)

* Use properly a transaction => [How to use jackc/pgx with connection pool, context, prepared statements etc](https://stackoverflow.com/a/76986702/13903942)
* Cool discussion [How to use pooling correctly? #1989](https://github.com/jackc/pgx/discussions/1989)

* Check out a "clean" db pool implementation: [clean-net-http)](https://github.com/Koubae/go-example/tree/master/workspace/web/clean-net-http)

