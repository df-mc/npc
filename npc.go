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
func Create(s Settings, w *world.World, f HandlerFunc) *player.Player {
	if w == nil {
		panic("world passed to npc.create must not be nil")
	}
	if f == nil {
		f = func(*player.Player) {}
	}
	npc := player.New(s.Name, s.Skin, s.Position)
	npc.Move(mgl64.Vec3{}, s.Yaw, s.Pitch)
	npc.SetScale(s.Scale)
	npc.SetHeldItems(s.MainHand, s.OffHand)
	npc.Armour().Set(s.Helmet, s.Chestplate, s.Leggings, s.Boots)
	if s.Immobile {
		npc.SetImmobile()
	}
	l := world.NewLoader(1, w, world.NopViewer{})
	h := &handler{f: f, l: l, vulnerable: s.Vulnerable}
	npc.Handle(h)
	w.AddEntity(npc)

	h.syncPosition(s.Position)
	go syncWorld(npc, l)
	return npc
}

// syncWorld periodically synchronises the world of the world.Loader passed with a player.Player's world. It stops doing
// so once the world returned by (*player.Player).World is nil.
func syncWorld(npc *player.Player, l *world.Loader) {
	t := time.NewTicker(time.Second / 20)
	defer t.Stop()

	for range t.C {
		if w := npc.World(); w != l.World() {
			if w == nil {
				// The NPC was closed in the meantime, stop synchronising the world.
				return
			}
			l.ChangeWorld(w)
		}
	}
}
