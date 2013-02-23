package core

import (
	gl "github.com/go3d/go-opengl/core"
)

func (me *EngineCore) onRender() {
	for thrRend.curCanvIndex = len(Core.Rendering.Canvases) - 1; thrRend.curCanvIndex >= 0; thrRend.curCanvIndex-- {
		if thrRend.curCanv = Core.Rendering.Canvases[thrRend.curCanvIndex]; thrRend.curCanv.renderThisFrame() {
			thrRend.curCanv.render()
		}
	}
}

func (me *RenderCanvas) render() {
	if !me.isFinal {
		me.frameBuf.Bind()
	}
	Core.Rendering.states.SetFramebufferSrgb(me.Srgb)
	for _, thrRend.curCam = range me.Cameras {
		thrRend.curCam.render()
	}
	Core.Rendering.states.SetFramebufferSrgb(false)
	if !me.isFinal {
		me.frameBuf.Unbind()
		thrRend.tmpQuadTex = &me.frameBufTex.Texture2D
	}
}

func (me *Camera) render() {
	if me.Enabled {
		thrRend.curTech = me.Rendering.Technique
		Core.Rendering.states.Apply(&me.thrRend.states)
		if me.Rendering.Viewport.shouldScissor {
			Core.Rendering.states.ForceEnableScissorTest()
		}
		if me.Rendering.Viewport.shouldScissor {
			gl.Scissor(me.Rendering.Viewport.glVpX, me.Rendering.Viewport.glVpY, me.Rendering.Viewport.glVpW, me.Rendering.Viewport.glVpH)
		}
		gl.Viewport(me.Rendering.Viewport.glVpX, me.Rendering.Viewport.glVpY, me.Rendering.Viewport.glVpW, me.Rendering.Viewport.glVpH)
		gl.Clear(me.thrRend.states.Other.ClearBits)
		thrRend.curScene = me.scene
		me.Rendering.Technique.render()
		if me.Rendering.Viewport.shouldScissor {
			Core.Rendering.states.ForceDisableScissorTest()
		}
	}
}

func (me *RenderTechniqueQuad) render() {
	thrRend.curMat, thrRend.curMatId = nil, ""
	me.glVao.Bind()
	me.fxTex.glTex = thrRend.tmpQuadTex
	thrRend.tmpTech, thrRend.tmpEffect = me, &me.Effect
	Core.useTechFx()
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	me.glVao.Unbind()
}

func (me *RenderTechniqueScene) render() {
	thrRend.curMat, thrRend.curMatId, thrRend.tmpTech = nil, "", me
	if thrRend.curScene != nil {
		thrRend.curScene.RootNode.renderChildren()
		thrRend.curScene.RootNode.renderSelf() // might be a skybox so "render" the root last
	}
}

func (me *Node) renderChildren() {
	for me.thrRend.curId, me.thrRend.curSubNode = range me.ChildNodes.M {
		me.thrRend.curSubNode.renderSelf()
		me.thrRend.curSubNode.renderChildren()
	}
}

func (me *Node) renderSelf() {
	if thrRend.curNode = me; me.model != nil {
		if thrRend.tmpMat = me.EffectiveMaterial(); thrRend.tmpMat != thrRend.curMat {
			if thrRend.curMat = thrRend.tmpMat; thrRend.curMat != nil {
				thrRend.tmpEffect = Core.Libs.Effects[thrRend.curMat.DefaultEffectID]
				Core.useTechFx()
			}
		}
		if me.Rendering.skyMode {
			Core.Rendering.states.DisableFaceCulling()
			gl.DepthFunc(gl.LEQUAL)
			thrRend.curProg.Uniform1i("uni_int_Sky", 1)
		}
		thrRend.curProg.UniformMatrix4fv("uni_mat4_VertexMatrix", 1, gl.FALSE, &me.thrRend.camProjMats[thrRend.curCam][0])
		me.model.render(me)
		if me.Rendering.skyMode {
			thrRend.curProg.Uniform1i("uni_int_Sky", 0)
			gl.DepthFunc(gl.LESS)
		}
	}
}

func (me *Model) render(node *Node) {
	me.mesh.render(node)
}

func (me *Mesh) render(node *Node) {
	if thrRend.curMeshBuf != me.meshBuffer {
		me.meshBuffer.use()
	}
	if thrRend.curMat.HasFaceEffects() {
		for thrRend.tmpFidx, thrRend.tmpFace = range me.raw.faces {
			thrRend.tmpEffect = thrRend.curMat.faceEffect(thrRend.tmpFace)
			Core.useTechFx()
			me.meshBuffer.use()
			thrRend.curProg.UniformMatrix4fv("uni_mat4_VertexMatrix", 1, gl.FALSE, &node.thrRend.camProjMats[thrRend.curCam][0])
			gl.DrawElementsBaseVertex(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.Util.PtrOffset(nil, uintptr(me.meshBufOffsetIndices+(int32(thrRend.tmpFidx)*3*4))), gl.Int(me.meshBufOffsetBaseIndex))
		}
	} else {
		gl.DrawElementsBaseVertex(gl.TRIANGLES, gl.Sizei(len(me.raw.indices)), gl.UNSIGNED_INT, gl.Util.PtrOffset(nil, uintptr(me.meshBufOffsetIndices)), gl.Int(me.meshBufOffsetBaseIndex))
	}
}
