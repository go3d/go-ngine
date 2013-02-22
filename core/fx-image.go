package core

import (
	"io"

	ugl "github.com/go3d/go-opengl/util"
)

type FxImage interface {
	Load() error
	Loaded() bool
	GpuSync() error
}

type FxImageStorage struct {
	DiskCache struct {
		Enabled      bool
		Compressor   func(w io.WriteCloser) io.WriteCloser
		Decompressor func(r io.ReadCloser) io.ReadCloser
	}
	Gpu struct {
		Bgra    bool
		UintRev bool
	}
}

type FxImageBase struct {
	PreProcess struct {
		FlipY    bool
		ToLinear bool
		ToBgra   bool
	}
	Storage FxImageStorage

	glTex    *ugl.TextureBase
	glSynced bool
}

func (me *FxImageBase) dispose() {
	me.GpuDelete()
}

func (me *FxImageBase) init(glTex *ugl.TextureBase) {
	me.glTex = glTex
	me.Storage = Options.Textures.Storage
	me.PreProcess.ToLinear, me.PreProcess.FlipY, me.PreProcess.ToBgra = true, true, me.Storage.Gpu.Bgra
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

func (me *FxImageBase) needPreproc() bool {
	return me.PreProcess.FlipY || me.PreProcess.ToBgra || me.PreProcess.ToLinear
}

func (me *FxImageBase) NoAutoMips() {
	me.glTex.MipMap.AutoGen = false
}
