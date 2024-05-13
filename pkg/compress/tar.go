package compress

import "fmt"

var (
	TarMakeTarGzip = _NewMakeFn("tar.gz", func(info *PathInfo) error {
		return _Tar{info: info}.MakeTarGzip()
	})
)

type _Tar struct {
	info *PathInfo
}

func (z _Tar) MakeTarGzip() error {
	return z.info.Exec("tar", "--transform", fmt.Sprintf("flags=r;s|%s|%s|", z.info.BaseName, z.info.TargetName), "-zcf", z.info.OutputPath, z.info.BaseName).Run()
}
