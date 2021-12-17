package npc

import (
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/go-gl/mathgl/mgl64"
)

// Settings holds different NPC settings such as the NPC's name, skin, position, etc. These values may be changed at
// runtime by calling the respective methods on the *player.Player returned by Create, the values passed in a Settings
// struct to Create are merely the initial values.
type Settings struct {
	Name       string
	Skin       skin.Skin
	Position   mgl64.Vec3
	Yaw, Pitch float64
	Scale      float64
	Immobile   bool
	Vulnerable bool
}
