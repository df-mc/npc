package npc

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

// HandlerFunc may be passed to Create to handle a *player.Player attacking an NPC.
type HandlerFunc func(p *player.Player)

// Create creates a new NPC with the Settings passed. A world.Loader is spawned in the background which follows the
// NPC to prevent it from despawning. Create panics if the world passed is nil.
// The HandlerFunc passed handles a player interacting with the NPC. Nil may be passed to avoid calling any function
// when the entity is interacted with.
// Create returns the *player.Player created. This entity has been added to the world passed. It may be removed from
// the world like any other entity by calling (*player.Player).Close.
func Create(s Settings, tx *world.Tx, f HandlerFunc) *player.Player {
	if tx == nil {
		panic("tx passed to npc.create must not be nil")
	}
	if f == nil {
		f = func(*player.Player) {}
	}
	opts := world.EntitySpawnOpts{
		Position: s.Position,
	}

	npc := opts.New(player.Type,
		player.Config{
			Name: s.Name,
		},
	)

	tx.AddEntity(npc)
	l := world.NewLoader(1, tx.World(), world.NopViewer{})
	h := &handler{f: f, l: l, vulnerable: s.Vulnerable}

	npc.ExecWorld(func(tx *world.Tx, e world.Entity) {
		pl := e.(*player.Player)
		pl.Move(mgl64.Vec3{}, s.Yaw, s.Pitch)
		pl.SetScale(s.Scale)
		pl.SetHeldItems(s.MainHand, s.OffHand)
		pl.Armour().Set(s.Helmet, s.Chestplate, s.Leggings, s.Boots)

		if s.Immobile {
			pl.SetImmobile()
		}

		pl.Handle(h)
	})

	h.syncPosition(tx, s.Position)
	go syncWorld(npc, l)

	e, ok := npc.Entity(tx)
	if !ok {
		panic("npc is not in a world")
	}
	return e.(*player.Player)
}

// syncWorld periodically synchronises the world of the world.Loader passed with a player.Player's world. It stops doing
// so once the world returned by (*player.Player).World is nil.
func syncWorld(npc *world.EntityHandle, l *world.Loader) {
	t := time.NewTicker(time.Second / 20)
	defer t.Stop()

	for range t.C {
		npc.ExecWorld(func(tx *world.Tx, e world.Entity) {
			if w := tx.World(); w != l.World() {
				if w == nil {
					// The NPC was closed in the meantime, stop synchronising the world.
					return
				}
				l.ChangeWorld(tx, w)
			}
		})
	}
}
