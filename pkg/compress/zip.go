package compress

import (
	"bytes"
	"context"
	"fmt"
)

var (
	ZipMakeZip = _NewMakeFn("zip", func(ctx context.Context, info *PathInfo) error {
		return _Zip{info: info}.MakeZip(ctx)
	})
)

type _Zip struct {
	info *PathInfo
}

func (z _Zip) Rename(ctx context.Context, dist, from, to string) error {
	cmd := z.info.Exec(ctx, "zipnote", "-w", dist)
	cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("@ %s\n@=%s\n", from, to))
	return cmd.Run()
}

func (z _Zip) MakeZip(ctx context.Context) error {
	if err := z.info.Exec(ctx, "zip", "-q", "-9", z.info.OutputPath, z.info.BaseName).Run(); err != nil {
		return err
	}
	return z.Rename(ctx, z.info.OutputPath, z.info.BaseName, z.info.TargetName)
}
