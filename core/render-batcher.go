package core

import (
	"sort"

	usl "github.com/metaleap/go-util/slice"
)

type RenderBatchCriteria int

const (
	BatchByProgram RenderBatchCriteria = 0
	BatchByTexture RenderBatchCriteria = 1
	BatchByBuffer  RenderBatchCriteria = 2
)

type renderBatchEntry struct {
	node, mesh, prog, fx int
	texes                []int
	face                 int32
}

type renderBatchList struct {
	all   []renderBatchEntry
	n     int
	prios [3]RenderBatchCriteria
}

func (me *renderBatchList) compare(i, j int, crit RenderBatchCriteria) (less, equal bool) {
	switch crit {
	case BatchByProgram:
		less = me.all[i].prog < me.all[j].prog
		equal = me.all[i].prog == me.all[j].prog
	case BatchByBuffer:
		less = Core.Libs.Meshes[me.all[i].mesh].meshBuffer.glIbo.GlHandle < Core.Libs.Meshes[me.all[j].mesh].meshBuffer.glIbo.GlHandle
		equal = Core.Libs.Meshes[me.all[i].mesh].meshBuffer.glIbo.GlHandle == Core.Libs.Meshes[me.all[j].mesh].meshBuffer.glIbo.GlHandle
	case BatchByTexture:
		if len(me.all[i].texes) > 0 {
			less, equal = false, true
			for t := 0; t < len(me.all[i].texes) && (less || equal); t++ {
				less = less || me.all[i].texes[t] < me.all[j].texes[t]
				equal = equal && me.all[i].texes[t] == me.all[j].texes[t]
			}
		} else {
			less, equal = len(me.all[j].texes) > 0, len(me.all[j].texes) == 0
		}
	}
	return
}

func (me *renderBatchList) Len() int {
	return me.n
}

func (me *renderBatchList) Less(i, j int) (less bool) {
	var eq bool
	if less, eq = me.compare(i, j, me.prios[0]); eq {
		if less, eq = me.compare(i, j, me.prios[1]); eq {
			less, _ = me.compare(i, j, me.prios[2])
		}
	}
	return
}

func (me *renderBatchList) Swap(i, j int) {
	me.all[i], me.all[j] = me.all[j], me.all[i]
}

type RenderBatcher struct {
	Enabled  bool
	Priority [3]RenderBatchCriteria
}

func (me *RenderTechniqueScene) prepEntry(entry *renderBatchEntry, nid int, fi int32, mesh *Mesh, effect *FxEffect) {
	var ti int
	entry.mesh, entry.node, entry.fx, entry.face = mesh.ID, nid, effect.ID, fi
	entry.prog = ogl.progs.Index(effect.uberPnames[me.name()])
	usl.IntEnsureLen(&entry.texes, Stats.Programs.maxTexUnits)
	for ti = 0; ti < len(entry.texes); ti++ {
		entry.texes[ti] = -1
	}
	ti = 0
	for oi := 0; oi < len(effect.FxProcs); oi++ {
		if effect.FxProcs[oi].IsTex() {
			entry.texes[ti] = effect.FxProcs[oi].Tex.ImageID
			ti++
		}
	}
	me.thrPrep.Done()
}

func (me *RenderTechniqueScene) prepBatch(scene *Scene, size int) {
	var (
		entry  *renderBatchEntry
		mesh   *Mesh
		mat    *FxMaterial
		effect *FxEffect
		fi, fl int32
		nid    int
	)
	b := &me.thrPrep.batch
	b.n = 0
	if len(b.all) < size {
		b.all = make([]renderBatchEntry, size)
	}
	for nid = 1; nid < len(scene.allNodes); nid++ {
		if scene.allNodes.Ok(nid) && me.Camera.thrPrep.nodeRender[nid] {
			if mesh, mat = scene.allNodes[nid].meshMat(); mat.HasFaceEffects() {
				for fi, fl = 0, int32(len(mesh.raw.faces)); fi < fl; fi++ {
					if effect = mat.faceEffect(&mesh.raw.faces[fi]); effect != nil {
						entry = &b.all[b.n]
						me.thrPrep.Add(1)
						go me.prepEntry(entry, nid, fi, mesh, effect)
						b.n++
					}
				}
			} else if effect = Core.Libs.Effects.get(mat.DefaultEffectID); effect != nil {
				entry = &b.all[b.n]
				me.thrPrep.Add(1)
				go me.prepEntry(entry, nid, -1, mesh, effect)
				b.n++
			}
		}
	}
	b.prios = me.Batch.Priority
	me.thrPrep.Wait()
	sort.Sort(b)
}