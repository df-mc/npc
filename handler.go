package npc

import (
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// handler ...
type handler struct {
	player.NopHandler

	l *world.Loader
	f func(attacker *player.Player)

	vulnerable bool
}

// HandleHurt ...
func (h *handler) HandleHurt(ctx *event.Context, _ *float64, src damage.Source) {
	if src, ok := src.(damage.SourceEntityAttack); ok {
		if attacker, ok := src.Attacker.(*player.Player); ok {
			h.f(attacker)
		}
	}
	if !h.vulnerable {
		ctx.Cancel()
	}
}

// HandleMove ...
func (h *handler) HandleMove(_ *event.Context, pos mgl64.Vec3, _, _ float64) {
	h.l.Move(pos)
}

// HandleTeleport ...
func (h *handler) HandleTeleport(_ *event.Context, pos mgl64.Vec3) {
	h.l.Move(pos)
}

// HandleQuit ...
func (h *handler) HandleQuit() {
	_ = h.l.Close()
}
