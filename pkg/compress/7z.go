package compress

import (
	"os"
	"os/exec"
	"path"
)

var (
	SevenZipMakeZip = _NewMakeFn("zip", func(info *PathInfo) error {
		return _SevenZip{info: info}.MakeZip()
	})
	SevenZipMakeTarGz = _NewMakeFn("tar", func(info *PathInfo) error {
		return _SevenZip{info: info}.MakeTarGzip()
	})
)

type _SevenZip struct {
	info *PathInfo
}

func (z _SevenZip) Exec(args ...string) *exec.Cmd {
	return z.info.Exec("7z", args...)
}

func (z _SevenZip) Add(method, dist, src string) error {
	return z.Exec("a", "-t"+method, dist, src, "-mx9").Run()
}

func (z _SevenZip) Rename(dist, from, to string) error {
	return z.Exec("rn", dist, from, to).Run()
}

func (z _SevenZip) MakeZip() error {
	err := z.Add("zip", z.info.OutputPath, z.info.BaseName)
	if err != nil {
		return err
	}
	return z.Rename(z.info.OutputPath, z.info.BaseName, z.info.TargetName)
}

func (z _SevenZip) MakeTarGzip() error {
	err := z.Add("tar", z.info.OutputPath, z.info.BaseName)
	if err != nil {
		return err
	}
	defer os.Remove(path.Join(z.info.Dir, z.info.OutputPath))

	if err = z.Rename(z.info.OutputPath, z.info.BaseName, z.info.TargetName); err != nil {
		return err
	}

	return z.Add("gzip", z.info.OutputPath+".gz", z.info.OutputPath)
}
