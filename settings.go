package npc

import (
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/go-gl/mathgl/mgl64"
)

// Settings holds different NPC settings such as the NPC's name, skin, position, etc.
type Settings struct {
	Name       string
	Skin       skin.Skin
	Position   mgl64.Vec3
	Yaw, Pitch float64
	Scale      float64
	Immobile   bool
	Vulnerable bool
}
