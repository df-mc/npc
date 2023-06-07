package npc

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

// handler implements the handler for an NPC entity. It manages the execution of the HandlerFunc assigned to the NPC
// and makes sure the *world.Loader's position remains synchronised with that of the NPC.
type handler struct {
	player.NopHandler

	l *world.Loader
	f HandlerFunc

	vulnerable bool
}

// HandleHurt ...
func (h *handler) HandleHurt(ctx *event.Context, _ *float64, _ *time.Duration, src world.DamageSource) {
	if src, ok := src.(entity.AttackDamageSource); ok {
		if attacker, ok := src.Attacker.(*player.Player); ok {
			h.f(attacker)
		}
	}
	if !h.vulnerable {
		ctx.Cancel()
	}
}

// HandleAttackEntity ...
func (h *handler) HandleAttackEntity(ctx *event.Context, e world.Entity, _, _ *float64, _ *bool) {
	if attacker, ok := e.(*player.Player); ok {
		h.f(attacker)
	}
	if !h.vulnerable {
		ctx.Cancel()
	}
}

// HandleMove ...
func (h *handler) HandleMove(_ *event.Context, pos mgl64.Vec3, _, _ float64) {
	h.syncPosition(pos)
}

// HandleTeleport ...
func (h *handler) HandleTeleport(_ *event.Context, pos mgl64.Vec3) {
	h.syncPosition(pos)
}

// syncPosition synchronises the position passed with the one in the world.Loader held by the handler. It ensures the
// chunk at this new position is loaded.
func (h *handler) syncPosition(pos mgl64.Vec3) {
	h.l.Move(pos)
	h.l.Load(1)
}

// HandleQuit ...
func (h *handler) HandleQuit() {
	_ = h.l.Close()
}
