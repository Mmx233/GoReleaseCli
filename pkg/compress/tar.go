package compress

import (
	"context"
	"fmt"
)

var (
	TarMakeTarGzip = _NewMakeFn("tar.gz", func(ctx context.Context, info *PathInfo) error {
		return _Tar{info: info}.MakeTarGzip(ctx)
	})
)

type _Tar struct {
	info *PathInfo
}

func (z _Tar) MakeTarGzip(ctx context.Context) error {
	return z.info.Exec(ctx, "tar", "--transform", fmt.Sprintf("flags=r;s|%s|%s|", z.info.BaseName, z.info.TargetName), "-zcf", z.info.OutputPath, z.info.BaseName).Run()
}
