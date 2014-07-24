package core

import (
	ugl "github.com/go3d/go-opengl/util"
)

type RenderTechniqueQuad struct {
	renderTechniqueBase

	Effect FxEffect

	glVao ugl.VertexArray
}

func newRenderTechniqueQuad(view *RenderView) RenderTechnique {
	me := &RenderTechniqueQuad{}
	me.init("Quad", view)
	me.Effect.init()
	me.Effect.FxProcs.EnableTex2D(-1).Tex.Sampler = Core.Render.Fx.Samplers.NoFilteringClamp
	me.Effect.UpdateRoutine()
	me.glVao.Create()
	me.glVao.Setup(nil, nil)
	return me
}

func (me *RenderTechniqueQuad) dispose() {
	me.glVao.Dispose()
	me.Effect.dispose()
	me.renderTechniqueBase.dispose()
}
