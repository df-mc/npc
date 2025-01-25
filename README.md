# npc
NPC library for Dragonfly.

## Getting started
The NPC library may be imported using `go get`:
```
go get github.com/df-mc/npc
```

## Usage
[![Go Reference](https://pkg.go.dev/badge/github.com/df-mc/npc.svg)](https://pkg.go.dev/github.com/df-mc/npc)

Usage of the NPC library is simple. It relies on the `Create` method:

```go
// var tx *world.Tx

settings := npc.Settings{
    Name: "Example NPC",
    Scale: 2,
    Position: mgl64.Vec3{1, 2, 3},
    Skin   ...,
}
npc.Create(settings, tx, nil)
```
Instead of `nil`, an `npc.HandlerFunc` may be passed to handle the NPC being hit by other
players.

Note that the `npc.Settings` passed initially may be overwritten by calling methods on
the `*player.Player` returned by `npc.Create()`.

The NPC library also contains convenience functions for reading/parsing skin data from files.
The `npc.Skin()`, `npc.*Model()` and `npc.*Texture()` functions may be used to do so.