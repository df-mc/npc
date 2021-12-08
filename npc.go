package npc

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

// New creates a new NPC with the Settings passed. A world loader is spawned in the background which follows the
// NPC to prevent it from despawning. The function passed handles a player attacking the NPC.
func New(s Settings, w *world.World, f func(*player.Player)) *player.Player {
	if f == nil {
		f = func(*player.Player) {}
	}
	npc := player.New(s.Name, s.Skin, s.Position)
	npc.Move(mgl64.Vec3{}, s.Yaw, s.Pitch)
	npc.SetScale(s.Scale)
	if s.Immobile {
		npc.SetImmobile()
	}

	l := world.NewLoader(1, w, world.NopViewer{})
	npc.Handle(&handler{f: f, l: l, vulnerable: s.Vulnerable})
	w.AddEntity(npc)

	go func() {
		t := time.NewTicker(time.Second / 20)
		defer t.Stop()

		for range t.C {
			if w := npc.World(); w != l.World() {
				l.ChangeWorld(w)
			}
		}
	}()
	return npc
}
