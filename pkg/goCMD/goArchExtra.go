package goCMD

type ArchExtra struct {
	EnvKey string
	Values []ArchExtraValue
}

type ArchExtraValue struct {
	// IsDefault describes whether this extra arch is set by default.
	IsDefault  bool
	ExtraFloat string
	// Value value of env.
	Value string
}

const (
	HardFloat = "hardfloat"
	SoftFloat = "softfloat"
)

var ExtraArches = map[string][]ArchExtra{
	"arm": {
		{
			EnvKey: "GOARM",
			Values: []ArchExtraValue{
				{Value: "5", ExtraFloat: HardFloat, IsDefault: true},
				{Value: "6", ExtraFloat: SoftFloat},
				{Value: "7", ExtraFloat: SoftFloat},
			},
		},
	},
	"386": {
		{
			EnvKey: "GO386",
			Values: []ArchExtraValue{
				{Value: "sse2", IsDefault: true},
				{Value: SoftFloat},
			},
		},
	},
	"amd64": {
		{
			EnvKey: "GOAMD64",
			Values: []ArchExtraValue{
				{Value: "v1", IsDefault: true},
				{Value: "v2"},
				{Value: "v3"},
				{Value: "v4"},
			},
		},
	},
	"wasm": {
		{
			EnvKey: "GOWASM",
			Values: []ArchExtraValue{
				{Value: "", IsDefault: true},
				{Value: "satconv"},
				{Value: "signext"},
			},
		},
	},
	"mips":     {_Mips},
	"mipsle":   {_Mips},
	"mips64":   {_Mips64},
	"mips64le": {_Mips64},
	"ppc64":    {_Ppc64},
	"ppc64le":  {_Ppc64},
}

var (
	_Mips = ArchExtra{
		EnvKey: "GOMIPS",
		Values: []ArchExtraValue{
			{Value: HardFloat, IsDefault: true},
			{Value: SoftFloat},
		},
	}
	_Mips64 = ArchExtra{
		EnvKey: "GOMIPS64",
		Values: []ArchExtraValue{
			{Value: HardFloat, IsDefault: true},
			{Value: SoftFloat},
		},
	}
	_Ppc64 = ArchExtra{
		EnvKey: "GOPPC64",
		Values: []ArchExtraValue{
			{Value: "power8", IsDefault: true},
			{Value: "power9"},
			{Value: "power10"},
		},
	}
)
