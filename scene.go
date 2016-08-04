package seker

import (
	"errors"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Geometry struct {
	Width  int32
	Height int32
}

type Application interface {
	GetSurface() *sdl.Surface
	GetWindow() *sdl.Window
}

type Scene struct {
	sync.Mutex
	App         Application
	Rect        sdl.Rect
	LayersStack []*Layer
	Layers      map[string]*Layer
	Changed     bool
	Geometry
}

var DefaultFont *Font
var BoldFont *Font

//NewScene constructor
func NewScene(app Application, size Geometry) *Scene {
	scene := new(Scene)
	scene.App = app
	scene.Geometry = size
	scene.Rect = sdl.Rect{0, 0, size.Width, size.Height}
	scene.Layers = map[string]*Layer{}

	LoadFonts(16)

	scene.AddLayer("root")
	scene.Draw()
	return scene
}

func (S *Scene) Run() {
	for {
		changed := S.Draw()
		if changed {
			S.App.GetWindow().UpdateSurface()
			S.Changed = false
		}
		// sdl.Delay(5)
		time.Sleep(5)
	}
}

func LoadFonts(size int) {
	DefaultFont = GetFont("Fantasque Regular.ttf", size)
	BoldFont = GetFont("Fantasque Bold.ttf", size)
}

func (S *Scene) Reset() {
	S.Layers = map[string]*Layer{}
	S.LayersStack = []*Layer{}
}

func (S *Scene) UpLayer(layerName string) {
	layer := S.Layers[layerName]
	var n int
	var l *Layer
	for n, l = range S.LayersStack {
		if l == layer {
			break
		}
	}
	S.LayersStack = append(S.LayersStack[:n], S.LayersStack[n+1:]...)
	S.LayersStack = append(S.LayersStack, layer)
}

func (S *Scene) Draw() bool {
	changed := S.GetChanged()
	if !changed {
		return changed
	}
	// fmt.Println(len(S.LayersStack))
	S.Clear()
	w := sync.WaitGroup{}
	for _, layer := range S.LayersStack {
		w.Add(1)
		go func(L *Layer) {
			s := L.Draw()
			S.Lock()
			s.Blit(&L.Rect, S.App.GetSurface(), &L.Rect)
			S.Unlock()
			w.Done()
		}(layer)
	}
	w.Wait()
	return changed
}

func (S *Scene) GetChanged() bool {
	changed := S.Changed
	for _, l := range S.LayersStack {
		ch := l.GetChanged()
		if !changed && ch {
			return true
		}
	}
	return changed
}

func (S *Scene) Clear() {
	S.App.GetSurface().FillRect(&S.Rect, 0xff242424)
}

func (S *Scene) RemoveLayer(name string) {
	S.Lock()
	defer S.Unlock()
	_, ok := S.Layers[name]
	if !ok {
		return
	}
	delete(S.Layers, name)
	for i, l := range S.LayersStack {
		if l.Name == name {
			l.Destroy()
			S.LayersStack = append(S.LayersStack[:i], S.LayersStack[i+1:]...)
			break
		}
	}
	S.Changed = true
}

func (S *Scene) AddLayer(name string) (*Layer, error) {
	amask := uint32(0xff000000)
	rmask := uint32(0x00ff0000)
	gmask := uint32(0x0000ff00)
	bmask := uint32(0x000000ff)
	var layer *Layer

	S.Lock()
	defer S.Unlock()
	layer, ok := S.Layers[name]
	if ok {
		return layer, errors.New("Use another layer name")
	}
	layer = new(Layer)
	S.Layers[name] = layer
	S.LayersStack = append(S.LayersStack, layer)
	layer.Name = name
	layer.Rect = S.Rect
	layer.Surface, _ = sdl.CreateRGBSurface(sdl.SWSURFACE, S.Width, S.Height, 32, rmask, gmask, bmask, amask)
	layer.Surface.FillRect(&layer.Rect, 0x00000000)
	S.Changed = true
	return layer, nil
}
