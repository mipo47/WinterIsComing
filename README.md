# Winter is coming

Multiplayer "tower defense" game, no dependencies, only standart Go packages

## Game rules
- There is a board of 10x30 cells (like a chess board), one side of the broad has a Zombie, another side has _The Wall_ with an Archer on it.
- Zombie is walking through the board each 1.5 seconds, aiming to reach _The Wall_
- Archer is trying to shoot the walking _Zombie_ from _The Wall_
- Zombie dies from single shot or reaches _The Wall_
- Players wins if all zombies will be killed until reaches _The Wall_

## Test
```
go run run.go test
```
## Start server, client, client with web
```
go run run.go server
go run run.go client
go run run.go web
```
## Build all
```
go run run.go build
```

Binaries can be found in /out directory (Windows/Ubuntu):
* Server - TCP server, hosts same game for multiple clients
* Client - TCP client, connects to server and kills zombies with other clients using simple AI
* Web (client_web) - Same as client, but adds HTML UI and allows to kill zombies using mouse clicks
