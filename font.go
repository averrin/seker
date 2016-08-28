package seker

import (
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

// Font abstraction with lock
type Font struct {
	sync.Mutex
	Name string
	Size int
	Font *ttf.Font
}

var cache map[string]map[int]*Font

// Draw text
func (F *Font) Draw(text string, color string) (*sdl.Surface, error) {
	F.Lock()
	defer F.Unlock()
	return F.Font.RenderUTF8_Blended(text, GetColor(color))
}

// GetFont with size
func GetFont(fname string, size int) *Font {
	if cache == nil {
		cache = map[string]map[int]*Font{}
	}
	cf := cache[fname]
	if cf != nil {
		f := cf[size]
		if f != nil {
			return f
		}
	} else {
		cache[fname] = map[int]*Font{}
	}
	fo := CreateFont(fname, size)
	cache[fname][size] = fo
	return fo
}

func CreateFont(fname string, size int) *Font {
	cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := filepath.Join(cwd, "fonts")
	font, _ := ttf.OpenFont(path.Join(dir, fname), size)
	fo := new(Font)
	fo.Name = fname
	fo.Size = size
	fo.Font = font
	return fo
}
