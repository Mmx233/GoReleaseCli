package compress

import (
	"bytes"
	"fmt"
)

var (
	ZipMakeZip = _NewMakeFn("zip", func(info *PathInfo) error {
		return _Zip{info: info}.MakeZip()
	})
)

type _Zip struct {
	info *PathInfo
}

func (z _Zip) Rename(dist, from, to string) error {
	cmd := z.info.Exec("zipnote", "-w", dist)
	cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("@ %s\n@=%s\n", from, to))
	return cmd.Run()
}

func (z _Zip) MakeZip() error {
	if err := z.info.Exec("zip", "-q", "-9", z.info.OutputPath, z.info.BaseName).Run(); err != nil {
		return err
	}
	return z.Rename(z.info.OutputPath, z.info.BaseName, z.info.TargetName)
}
