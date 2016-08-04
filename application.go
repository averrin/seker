package seker

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

type ApplicationInterface interface {
	GetSurface() *sdl.Surface
	GetWindow() *sdl.Window
}

type Application struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Surface  *sdl.Surface
	Scene    *Scene
}

func (app *Application) GetSurface() *sdl.Surface {
	return app.Surface
}

func (app *Application) GetWindow() *sdl.Window {
	return app.Window
}

func (app *Application) Init(title string, flag uint32, geometry Geometry) {
	sdl.Init(sdl.INIT_EVERYTHING)
	ttf.Init()

	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int(geometry.Width), int(geometry.Height), flag)
	if err != nil {
		panic(err)
	}
	app.Window = window
	renderer, err := sdl.CreateRenderer(app.Window, -1, sdl.RENDERER_ACCELERATED)
	surface, err := app.Window.GetSurface()
	if err != nil {
		panic(err)
	}
	app.Renderer = renderer
	app.Surface = surface
	renderer.Clear()
	app.Scene = NewScene(*app, geometry)
	renderer.Present()
	time.Sleep(5)
	app.Window.UpdateSurface()
}

func (app *Application) Close() {
	app.Window.Destroy()
}
