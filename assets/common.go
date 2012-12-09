package assets

const (
	TRANSFORM_TYPE_LOOKAT    = 0
	TRANSFORM_TYPE_MATRIX    = iota
	TRANSFORM_TYPE_ROTATE    = iota
	TRANSFORM_TYPE_SKEW      = iota
	TRANSFORM_TYPE_SCALE     = iota
	TRANSFORM_TYPE_TRANSLATE = iota
)

type HasAsset struct {
	Asset *Asset
}

type HasExtras struct {
	Extras []*Extra
}

type HasID struct {
	//	The unique identifier of this *Def*, *Inst* or *Lib*.
	ID string
}

type HasName struct {
	//	The optional pretty-print name/title of this *Def*, *Inst* or *Lib*.
	Name string
}

type HasSid struct {
	//	Scoped ID
	Sid string
}

type HasTechniques struct {
	Techniques []*Technique
}

type Asset struct {
	HasExtras
	Created      string
	Modified     string
	Keywords     string
	Revision     string
	Subject      string
	Title        string
	UpAxis       string
	Contributors []*Contributor
	Coverage     *AssetGeographicLocation
}

type AssetGeographicLocation struct {
	Longitude        float64
	Latitude         float64
	Altitude         float64
	AltitudeAbsolute bool
}

type BindMaterial struct {
	HasExtras
	HasTechniques
	Params          []*Param
	TechniqueCommon struct {
		MaterialInsts []*FxMaterialInst
	}
}

type Contributor struct {
	Author        string
	AuthorEmail   string
	AuthorWebsite string
	AuthoringTool string
	Comments      string
	Copyright     string
	SourceData    string
}

type Extra struct {
	HasID
	HasName
	HasAsset
	HasTechniques
	Type string
}

type IndexedInputs struct {
	Count   uint64
	Inputs  []*InputShared
	Vcount  []int64
	Indices []int64
}

type Input struct {
	Semantic string
	Source   string
}

type InputShared struct {
	Input
	Offset uint64
	Set    uint64
}

type Layers map[string]bool

type Param struct {
	HasName
	HasSid
	Semantic string
	Type     string
}

type ParamDef struct {
	HasSid
	Value interface{}
}

type ParamDefs map[string]*ParamDef

type ParamInst struct {
	Ref   string
	Value interface{}
}

type Technique struct {
	Profile string
	Data    interface{}
}

type Transform struct {
	HasSid
	Type int
	F    []float64
}
