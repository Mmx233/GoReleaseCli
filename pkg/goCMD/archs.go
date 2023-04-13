package goCMD

var Arch = map[string][]string{
	"aix":       {"ppc64"},
	"android":   { /*"386", "amd64", "arm", */ "arm64"},
	"darwin":    {"amd64", "arm64"},
	"dragonfly": {"amd64"},
	"freebsd":   {"386", "amd64", "arm", "arm64"},
	"illumos":   {"amd64"},
	/*"ios":       {"amd64" , "arm64"},*/
	"js":      {"wasm"},
	"linux":   {"386", "amd64", "arm", "arm64", "mips", "mips64", "mips64le", "mipsle", "ppc64", "ppc64le", "riscv64", "s390x"},
	"netbsd":  {"386", "amd64", "arm", "arm64"},
	"openbsd": {"386", "amd64", "arm", "arm64", "mips64"},
	"plan9":   {"386", "amd64", "arm"},
	"solaris": {"amd64"},
	"windows": {"386", "amd64", "arm", "arm64"},
}
