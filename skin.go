package npc

import (
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

// ParseSkin parses a skin.Skin from a texture and model path. The texture path must point to a valid image file, such
// as a PNG file, with correct dimensions (64x32, 64x64 or 128x128) and the model path must point to a JSON encoded file
// that holds a skin model.
// The parsed skin.Skin is returned if reading both files was successful and if the data contained in them was valid.
// ParseSkin ensures the dimensions specified in the texture and model match and are valid.
func ParseSkin(texturePath, modelPath string) (skin.Skin, error) {
	texture, err := os.Open(texturePath)
	if err != nil {
		return skin.Skin{}, fmt.Errorf("failed opening texture file: %w", err)
	}
	model, err := os.Open(modelPath)
	if err != nil {
		return skin.Skin{}, fmt.Errorf("failed opening model file: %w", err)
	}
	defer func() {
		_ = texture.Close()
		_ = model.Close()
	}()
	return ReadSkin(texture, model)
}

// ReadSkin reads a skin.Skin from a texture and model io.Reader. The texture reader must hold valid image data of a
// format such as PNG with correct dimensions (64x32, 64x64 or 128x128) and the model reader must hold JSON data that
// holds a skin model.
// The parsed skin.Skin is returned if reading both files was successful and if the data contained in them was valid.
// ReadSkin ensures the dimensions specified in the texture and model match and are valid.
func ReadSkin(texture, model io.Reader) (skin.Skin, error) {
	pix, rect, err := readTexture(texture)
	if err != nil {
		return skin.Skin{}, err
	}
	m, config, mrect, err := readModel(model)
	if err != nil {
		return skin.Skin{}, err
	}
	// Verify that the dimensions of the texture we read match those specified in the model. Clientside behaviour is
	// unreliable if this is not the case.
	if rect != mrect {
		return skin.Skin{}, fmt.Errorf("skin texture dimensions did not match those specified in model: %v specified but got %v", mrect, rect)
	}

	s := skin.New(rect.Dx(), rect.Dy())
	s.ModelConfig = config
	s.Model = m
	s.Pix = pix
	return s, nil
}

// readModel parses a JSON model from a file at a path passed and returns the raw JSON data. In addition to that, it
// returns a parsed skin.ModelConfig and the bounds of the skin texture as specified in the model.
// If the file could not be parsed or if the model data was invalid, an error is returned.
func readModel(r io.Reader) ([]byte, skin.ModelConfig, image.Rectangle, error) {
	model, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, skin.ModelConfig{}, image.Rectangle{}, fmt.Errorf("failed reading model: %w", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(model, &m); err != nil {
		return nil, skin.ModelConfig{}, image.Rectangle{}, fmt.Errorf("failed decoding model: %w", err)
	}

	data := m["minecraft:geometry"].([]interface{})[0].(map[string]interface{})["description"].(map[string]interface{})

	// The model contains the texture width and height too. We return these as an image.Rectangle and later verify if
	// this matches the dimensions of the actual texture.
	w, h := int(data["texture_width"].(float64)), int(data["texture_height"].(float64))
	return model, skin.ModelConfig{Default: data["identifier"].(string)}, image.Rect(0, 0, w, h), nil
}

// readTexture parses a skin texture from the path passed and returns an RGBA-ordered byte slice, where each pixel is
// represented by 4 bytes. An error is returned if the image could not be opened or was otherwise invalid as a skin.
// The bounds of the image parsed are returned as an image.Rectangle.
func readTexture(r io.Reader) ([]byte, image.Rectangle, error) {
	img, s, err := image.Decode(r)
	if err != nil {
		return nil, image.Rectangle{}, fmt.Errorf("failed decoding texture: %v: %v", err, s)
	}

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if !(w == 64 && h == 32) && !(w == 64 && h == 64) && !(w == 128 && h == 128) {
		return nil, img.Bounds(), fmt.Errorf("invalid skin texture dimensions: %vx%v", w, h)
	}

	data := make([]byte, 0, w*h*4)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			data = append(data, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}
	return data, img.Bounds(), nil
}
