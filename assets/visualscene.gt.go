package assets

//	Embodies the entire set of information that can be visualized from the contents of the *Lib* libraries.
//	The hierarchical structure of the Visual Scene is organized into a scene graph. A scene graph is a
//	directed acyclic graph or tree data structure that contains nodes of visual information and related data.
type VisualSceneDef struct {
	BaseDef
}

func (me *VisualSceneDef) init() {
}

//	An instance referencing a Visual Scene definition.
type VisualSceneInst struct {
	BaseInst

	//	The Visual Scene definition referenced by this instance.
	Def *VisualSceneDef
}

func (me *VisualSceneInst) init() {
}

//#begin-gt _definstlib.gt T:VisualScene

func newVisualSceneDef(id string) (me *VisualSceneDef) {
	me = &VisualSceneDef{}
	me.BaseDef.init(id)
	me.init()
	return
}

//	Creates and returns a new *VisualSceneInst* instance referencing this *VisualSceneDef* definition.
func (me *VisualSceneDef) NewInst(id string) (inst *VisualSceneInst) {
	inst = &VisualSceneInst{Def: me}
	inst.Base.init(id)
	inst.init()
	return
}

var (
	//	A *map* collection that contains *LibVisualSceneDefs* libraries associated by their *ID*.
	AllVisualSceneDefLibs = LibsVisualSceneDef{}

	//	The "default" *LibVisualSceneDefs* library for *VisualSceneDef*s.
	VisualSceneDefs = AllVisualSceneDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllVisualSceneDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllVisualSceneDefLibs* variable: a *map* collection that contains
//	*LibVisualSceneDefs* libraries associated by their *ID*.
type LibsVisualSceneDef map[string]*LibVisualSceneDefs

//	Creates a new *LibVisualSceneDefs* library with the specified *ID*, adds it to this *LibsVisualSceneDef*, and returns it.
//	
//	If this *LibsVisualSceneDef* already contains a *LibVisualSceneDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsVisualSceneDef) AddNew(id string) (lib *LibVisualSceneDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsVisualSceneDef) new(id string) (lib *LibVisualSceneDefs) {
	lib = newLibVisualSceneDefs(id)
	return
}

//	A library that contains *VisualSceneDef*s associated by their *ID*. To create a new *LibVisualSceneDefs* library, ONLY
//	use the *LibsVisualSceneDef.New()* or *LibsVisualSceneDef.AddNew()* methods.
type LibVisualSceneDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*VisualSceneDef
}

func newLibVisualSceneDefs(id string) (me *LibVisualSceneDefs) {
	me = &LibVisualSceneDefs{M: map[string]*VisualSceneDef{}}
	me.Base.init(id)
	return
}

//	Adds the specified *VisualSceneDef* definition to this *LibVisualSceneDefs*, and returns it.
//	
//	If this *LibVisualSceneDefs* already contains a *VisualSceneDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibVisualSceneDefs) Add(d *VisualSceneDef) (n *VisualSceneDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *VisualSceneDef* definition with the specified *ID*, adds it to this *LibVisualSceneDefs*, and returns it.
//	
//	If this *LibVisualSceneDefs* already contains a *VisualSceneDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibVisualSceneDefs) AddNew(id string) *VisualSceneDef { return me.Add(me.New(id)) }

//	Creates a new *VisualSceneDef* definition with the specified *ID* and returns it, but does not add it to this *LibVisualSceneDefs*.
func (me *LibVisualSceneDefs) New(id string) (def *VisualSceneDef) { def = newVisualSceneDef(id); return }

//	Removes the *VisualSceneDef* with the specified *ID* from this *LibVisualSceneDefs*.
func (me *LibVisualSceneDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibVisualSceneDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibVisualSceneDefs* library or its *VisualSceneDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibVisualSceneDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
