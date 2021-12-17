# npc
NPC library for Dragonfly.

## Getting started
The NPC library may be imported using `go get`:
```
go get github.com/df-mc/npc
```

## Usage
Usage of the NPC library is simple. It relies on the `Create` method:

```go
// var w *world.World

settings := npc.Settings{
    Name: "Example NPC",
    Scale: 2,
    Position: mgl64.Vec3{1, 2, 3},
    Skin   ...,
}
p := npc.Create(settings, w, nil)
p.SwingArm()
```
Instead of `nil`, an `npc.HandlerFunc` may be passed to handle the NPC being hit by other
players.

Note that the `npc.Settings` passed initially may be overwritten by calling methods on
the `*player.Player` returned by `npc.Create`.