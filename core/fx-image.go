package core

import (
	ugl "github.com/go3d/go-opengl/util"
)

type FxImageBase struct {
	OnLoad   FxImageOnLoad
	InitFrom struct {
		RawData []byte
		RefUrl  string
	}

	glTex    *ugl.TextureBase
	glSynced bool
}

func (me *FxImageBase) dispose() {
	me.GpuDelete()
}

func (me *FxImageBase) init(glTex *ugl.TextureBase) {
	me.glTex = glTex
}

func (me *FxImageBase) gpuSync(tex ugl.Texture) (err error) {
	if err = tex.Recreate(); err == nil {
		me.glSynced = true
	}
	return
}

func (me *FxImageBase) GpuDelete() {
	me.glTex.Dispose()
	me.glSynced = false
}

func (me *FxImageBase) GpuSynced() bool {
	return me.glSynced
}

func (me *FxImageBase) NoAutoMips() {
	me.glTex.MipMap.AutoGen = false
}

type FxImageOnLoad func(img interface{}, err error)
