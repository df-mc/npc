package npc

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/df-mc/dragonfly/server/player/skin"
	"image"
	"io"
	"io/ioutil"
	"os"
)

// Always import image/png so that image.Decode can always decode PNGs. By far most of the skins are stored as PNGs so
// it seems reasonable enough to do this.
import _ "image/png"

var (
	//go:embed default_model.json
	defaultModel []byte
	// DefaultModel is a Model for a default custom skin. It is a model with a 64x64 texture.
	DefaultModel, _ = ReadModel(bytes.NewReader(defaultModel))
)

// MustSkin creates a new skin.Skin using the Texture and Model passed. It panics if the dimensions as specified in the
// model don't match those of the Texture.
func MustSkin(tex Texture, mod Model) skin.Skin {
	s, err := Skin(tex, mod)
	if err != nil {
		panic(err)
	}
	return s
}

// Skin creates a new skin.Skin using the Texture and Model passed. It returns an error if the dimensions as specified
// in the Model did not match those of the Texture.
func Skin(tex Texture, mod Model) (skin.Skin, error) {
	// Verify that the dimensions of the texture we read match those specified in the model. Clientside behaviour is
	// unreliable if this is not the case.
	if tex.rect != mod.rect {
		return skin.Skin{}, fmt.Errorf("skin texture dimensions did not match those specified in model: %v specified but got %v", mod.rect, tex.rect)
	}

	s := skin.New(tex.rect.Dx(), tex.rect.Dy())
	s.ModelConfig = mod.conf
	s.Model = mod.json
	s.Pix = tex.pix
	return s, nil
}

// Model is the model of a skin.Skin parsed with ParseModel or ReadModel.
type Model struct {
	json []byte
	conf skin.ModelConfig
	rect image.Rectangle
}

// Texture is the texture of a skin.Skin parsed with ParseTexture or ReadTexture.
type Texture struct {
	pix  []byte
	rect image.Rectangle
}

// MustParseModel parses a Model from a JSON file at a specific path. It panics if the file could not be opened or if
// the JSON contained was invalid or not a proper model.
func MustParseModel(path string) Model {
	mod, err := ParseModel(path)
	if err != nil {
		panic(err)
	}
	return mod
}

// MustParseTexture parses a Texture from an image file at a specific path. It panics if the file could not be opened or
// if it did not contain a valid skin image.
func MustParseTexture(path string) Texture {
	tex, err := ParseTexture(path)
	if err != nil {
		panic(err)
	}
	return tex
}

// ParseModel parses a Model from a JSON file at a specific path. An error is returned if the file could not be opened
// or if the JSON contained was invalid or not a proper model.
func ParseModel(path string) (Model, error) {
	model, err := os.Open(path)
	if err != nil {
		return Model{}, fmt.Errorf("failed opening model file: %w", err)
	}
	defer model.Close()
	return ReadModel(model)
}

// ParseTexture parses a Texture from an image file at a specific path. An error is returned if the file could not be
// opened or if it did not contain a valid skin image.
func ParseTexture(path string) (Texture, error) {
	texture, err := os.Open(path)
	if err != nil {
		return Texture{}, fmt.Errorf("failed opening texture file: %w", err)
	}
	defer texture.Close()
	return ReadTexture(texture)
}

// ReadModel parses a JSON model from a file at a path passed and returns a Model. It contains a parsed skin.ModelConfig
// and the bounds of the skin texture as specified in the model. If the file could not be parsed or if the model data
// was invalid, an error is returned.
func ReadModel(r io.Reader) (Model, error) {
	model, err := ioutil.ReadAll(r)
	if err != nil {
		return Model{}, fmt.Errorf("failed reading model: %w", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(model, &m); err != nil {
		return Model{}, fmt.Errorf("failed decoding model: %w", err)
	}

	data := m["minecraft:geometry"].([]interface{})[0].(map[string]interface{})["description"].(map[string]interface{})

	// The model contains the texture width and height too. We return these as an image.Rectangle and later verify if
	// this matches the dimensions of the actual texture.
	w, h := int(data["texture_width"].(float64)), int(data["texture_height"].(float64))
	return Model{json: model, conf: skin.ModelConfig{Default: data["identifier"].(string)}, rect: image.Rect(0, 0, w, h)}, nil
}

// ReadTexture parses a skin texture from the path passed and returns a Texture struct with data, where each pixel is
// represented by 4 bytes. An error is returned if the image could not be opened or was otherwise invalid as a skin.
func ReadTexture(r io.Reader) (Texture, error) {
	img, s, err := image.Decode(r)
	if err != nil {
		return Texture{}, fmt.Errorf("failed decoding texture: %v: %v", err, s)
	}

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if !(w == 64 && h == 32) && !(w == 64 && h == 64) && !(w == 128 && h == 128) {
		return Texture{}, fmt.Errorf("invalid skin texture dimensions: %vx%v", w, h)
	}

	data := make([]byte, 0, w*h*4)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			data = append(data, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}
	return Texture{pix: data, rect: img.Bounds()}, nil
}
