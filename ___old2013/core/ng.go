package core

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/metaleap/go-util-fs"
	ugl "github.com/metaleap/go-opengl/util"

	ngctx "github.com/metaleap/go-ngine/glctx"
)

//	Call this to "un-init" go:ngine and to release any and all GPU or RAM resources still allocated.
func Dispose() {
	Core.dispose()
	ogl.dispose()
	UserIO.dispose()
}

//	Initializes go:ngine; this first attempts to initialize OpenGL and then open a window to your supplied specifications with a GL 3.3-or-higher profile.
func Init(fullscreen bool, ctx ngctx.CtxProvider) (err error) {
	var (
		glVerIndex         = len(ugl.KnownVersions) - 1
		badVer, tmpDirPath string
		glVer              float64
	)
	UserIO.ctx = ctx
	if err = ufs.EnsureDirExists(Core.fileIO.resolveLocalFilePath(Options.AppDir.Temp.BaseName)); err != nil {
		return
	}
	defer runtime.GC()
	if len(Options.AppDir.Temp.BaseName) > 0 {
		for _, diagTmpDirName := range []string{Options.AppDir.Temp.ShaderSources, Options.AppDir.Temp.CachedTextures} {
			tmpDirPath = Core.fileIO.resolveLocalFilePath(filepath.Join(Options.AppDir.Temp.BaseName, diagTmpDirName))
			if err = ufs.EnsureDirExists(tmpDirPath); err != nil {
				return
			}
		}
		for _, diagTmpDirName := range []string{Options.AppDir.Temp.ShaderSources} {
			tmpDirPath = Core.fileIO.resolveLocalFilePath(filepath.Join(Options.AppDir.Temp.BaseName, diagTmpDirName))
			if err = ufs.ClearDirectory(tmpDirPath); err != nil {
				return
			}
		}
	}
	if Options.Initialization.GlContext.CoreProfile.ForceFirst {
		for i, v := range ugl.KnownVersions {
			if v == Options.Initialization.GlContext.CoreProfile.VersionHint {
				glVerIndex = i
				break
			}
		}
	}
tryInit:
	glVer, UserIO.Window.fullscreen = ugl.KnownVersions[glVerIndex], fullscreen
	if err = UserIO.init(glVer); err == nil {
		if err, badVer = ogl.init(); err == nil && len(badVer) == 0 {
			Stats.reset()
			Loop.init()
			if err = Core.init(); err != nil {
				return
			}
			Diag.LogIfGlErr("INIT")
		} else if len(badVer) > 0 && !Options.Initialization.GlContext.CoreProfile.ForceFirst {
			Options.Initialization.GlContext.CoreProfile.ForceFirst = true
			UserIO.isCtxInit, UserIO.Window.isCreated = false, false
			goto tryInit
		}
	} else if Options.Initialization.GlContext.CoreProfile.ForceFirst && (glVerIndex > 0) {
		glVerIndex--
		UserIO.isCtxInit, UserIO.Window.isCreated = false, false
		goto tryInit
	} else {
		badVer = ogl.lastBadVer
	}
	if len(badVer) > 0 {
		err = errf(ogl.versionErrorMessage(glMinVerStr, badVer))
	}
	return
}

func errf(format string, fmtArgs ...interface{}) error {
	return fmt.Errorf(format, fmtArgs...)
}

func strf(format string, fmtArgs ...interface{}) string {
	return fmt.Sprintf(format, fmtArgs...)
}
