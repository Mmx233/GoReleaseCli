package compress

import (
	"context"
	"os"
	"os/exec"
	"path"
)

var (
	SevenZipMakeZip = _NewMakeFn("zip", func(ctx context.Context, info *PathInfo) error {
		return _SevenZip{info: info}.MakeZip(ctx)
	})
	SevenZipMakeTarGzip = _NewMakeFn("tar", func(ctx context.Context, info *PathInfo) error {
		return _SevenZip{info: info}.MakeTarGzip(ctx)
	})
)

type _SevenZip struct {
	info *PathInfo
}

func (z _SevenZip) Exec(ctx context.Context, args ...string) *exec.Cmd {
	return z.info.Exec(ctx, "7z", args...)
}

func (z _SevenZip) Add(ctx context.Context, method, dist, src string) error {
	return z.Exec(ctx, "a", "-t"+method, dist, src, "-mx9").Run()
}

func (z _SevenZip) Rename(ctx context.Context, dist, from, to string) error {
	return z.Exec(ctx, "rn", dist, from, to).Run()
}

func (z _SevenZip) MakeZip(ctx context.Context) error {
	err := z.Add(ctx, "zip", z.info.OutputPath, z.info.BaseName)
	if err != nil {
		return err
	}
	return z.Rename(ctx, z.info.OutputPath, z.info.BaseName, z.info.TargetName)
}

func (z _SevenZip) MakeTarGzip(ctx context.Context) error {
	err := z.Add(ctx, "tar", z.info.OutputPath, z.info.BaseName)
	if err != nil {
		return err
	}
	defer os.Remove(path.Join(z.info.Dir, z.info.OutputPath))

	if err = z.Rename(ctx, z.info.OutputPath, z.info.BaseName, z.info.TargetName); err != nil {
		return err
	}

	return z.Add(ctx, "gzip", z.info.OutputPath+".gz", z.info.OutputPath)
}
