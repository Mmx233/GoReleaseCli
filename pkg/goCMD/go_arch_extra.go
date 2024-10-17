package goCMD

// https://pkg.go.dev/cmd/go#hdr-Print_Go_environment_information

type ArchExtra struct {
	EnvKey string
	Values []ArchExtraValue
}

type ArchExtraValue struct {
	// IsDefault describes whether this extra arch is set by default.
	IsDefault  bool
	ExtraFloat string

	// Name maybe empty string, use for replace value of Value in result file name
	Name string
	// Value value of env.
	Value string
}

func (v ArchExtraValue) Names(showDefault bool) []string {
	var names []string
	if !v.IsDefault || showDefault {
		if v.Name != "" {
			names = []string{v.Name}
		} else if v.Value != "" {
			names = []string{v.Value}
		}
	}
	return names
}
func (v ArchExtraValue) NamesExtraFloat(showDefault bool) []string {
	names := v.Names(showDefault)
	return append(names, v.ExtraFloat)
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
				{Value: "5", Name: "v5", ExtraFloat: HardFloat, IsDefault: true},
				{Value: "6", Name: "v6", ExtraFloat: SoftFloat},
				{Value: "7", Name: "v7", ExtraFloat: SoftFloat},
			},
		},
	},
	"arm64": {
		{
			EnvKey: "GOARM64",
			Values: []ArchExtraValue{
				{Value: "v8.0", Name: "v8", IsDefault: true},
				// {Value: "v8.{1-9}"},
				{Value: "v9.0", Name: "v9"},
				// {Value: "v9.{1-5}"}
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
	"riscv64": {
		{
			EnvKey: "GORISCV64",
			Values: []ArchExtraValue{
				{Value: "rva20u64", IsDefault: true},
				{Value: "rva22u64"},
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
