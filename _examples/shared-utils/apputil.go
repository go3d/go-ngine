package exampleutils

import (
	"fmt"
	"runtime"
	"time"

	ng "github.com/go3d/go-ngine/core"
	unum "github.com/metaleap/go-util/num"
)

var (
	ArtificialSplashScreenDelay = 0 * time.Second

	//	The RenderCanvas the example scene is initially being rendered to. This is an off-screen "render-to-texture" RenderCanvas.
	SceneCanvas *ng.RenderCanvas

	//	The primary scene-rendering camera, rendering to SceneCanvas.
	SceneCam *ng.Camera

	//	Unlike the off-screen (render-to-texture) SceneCanvas above, this RenderCanvas epresents the actual screen/window.
	PostFxCanvas *ng.RenderCanvas

	//	Takes the image rendered to SceneCanvas, may post-process it or not, and blits it to PostFxCanvas.
	PostFxCam *ng.Camera

	//	Do not set this field directly, only use PauseResume() to toggle it and effect the associated render-state changes.
	Paused bool

	retro  bool
	numCgo struct {
		preLoop  int64
		postLoop int64
	}
	winTitle struct {
		cw, ch         int
		camPos, camDir unum.Vec3
	}

	curKeyHint = 0
	sec        = 0
)

//	Returns the window title to be set by OnSec().
func appWindowTitle() string {
	winTitle.cw, winTitle.ch = ng.UserIO.Window.Width(), ng.UserIO.Window.Height()
	if SceneCanvas != nil {
		winTitle.cw, winTitle.ch = SceneCanvas.CurrentAbsoluteSize()
	}
	if SceneCam != nil {
		winTitle.camPos, winTitle.camDir = SceneCam.Controller.Pos, *SceneCam.Controller.Dir()
	}
	return fmt.Sprintf("%v FPS @ %vx%v   |   %s   |   Cam: P=%v D=%v", ng.Stats.FpsLastSec, winTitle.cw, winTitle.ch, KeyHints[curKeyHint], winTitle.camPos.String(), winTitle.camDir.String())
}

//	Refreshes the window title every second, showing the next one entry in KeyHints every 3 seconds.
func OnSec() {
	if sec++; sec == 3 {
		sec = 0
		if curKeyHint++; (curKeyHint > MaxKeyHint) || (curKeyHint >= (len(KeyHints))) {
			curKeyHint = 0
		}
	}
	ng.UserIO.Window.SetTitle(appWindowTitle())
}

//	Called by each example-app's func main(). Initializes go:ngine, sets SceneCam/SceneCanvas/PostFxCam/PostFxCanvas etc., calls the specified assetLoader function, then enters The Loop.
func Main(setupExampleScene, onAppThread, onWinThread func()) {
	//	go:ngine doesn't do this for you by design:
	runtime.LockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())

	width, height, fullscreen := 1280, 720, false
	// width, height, fullscreen := 1920, 1080, true
	win := ng.NewWindowOptions(fmt.Sprintf("Loading example app... (%v CPU cores)", runtime.GOMAXPROCS(0)), width, height, fullscreen)

	// release apps shouldn't do this, but during dev/test we want to observe max fps:
	win.SetSwapInterval(0)

	opt := ng.NewEngineOptions(AssetRootDirPath(), win)

	//	While the default for this (true on Macs only) is reasonable for release apps at present,
	//	here we force core profile to verify all of go:ngine's GL code is fully core-profile compliant
	opt.Initialization.GlContext.CoreProfile.ForceFirst = true

	//	Release apps shouldn't do this, but here we're verifying everything runs in the oldest-supported GL version:
	opt.Initialization.GlContext.CoreProfile.VersionHint = 3.3

	//	Worth toggling this every once in a while just to see whether it makes a perf diff at all...
	realThreads := false
	opt.Loop.ForceThreads.App, opt.Loop.ForceThreads.Prep = realThreads, realThreads

	//	STEP 1: init go:ngine
	if err := ng.Init(opt); err != nil {
		fmt.Printf("ABORT:\n%v\n", err)
	} else {
		defer ng.Dispose()
		defer PrintPostLoopSummary()

		//	STEP 2: post-init, pre-loop setup
		ng.Loop.On.EverySec, ng.Loop.On.AppThread, ng.Loop.On.WinThread = OnSec, onAppThread, onWinThread

		PostFxCanvas = ng.Core.Rendering.Canvases.Final()
		PostFxCam = PostFxCanvas.Cameras[0]

		if setupExampleScene != nil {
			SceneCanvas = ng.Core.Rendering.Canvases.AddNew(true, 1, 1)
			SceneCam = SceneCanvas.AddNewCamera3D()
			SceneCam.Rendering.States.ClearColor.Set(0.5, 0.6, 0.85, 1)
			setupExampleScene()
			ng.Core.SyncUpdates()
		}
		time.Sleep(ArtificialSplashScreenDelay) // change to higher value to check out your splash-screen
		numCgo.preLoop = runtime.NumCgoCall()

		//	STEP 3: enter... Da Loop.
		ng.Loop.Loop()
		numCgo.postLoop = runtime.NumCgoCall()
	}
}