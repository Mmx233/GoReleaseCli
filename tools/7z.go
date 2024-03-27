package tools

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
)

type CompressInfo struct {
	// src info
	BaseName string // filename without parent path
	Ext      string
	Name     string // basename without ext
	Dir      string

	// dist info
	OutputPath string
	TargetName string // rename src file
}

func DecodeCompressPath(filePath, targetName, targetFormat string) *CompressInfo {
	baseName := path.Base(filePath)
	ext := path.Ext(baseName)
	name := strings.TrimSuffix(baseName, ext)
	dir := strings.TrimSuffix(filePath, baseName)

	return &CompressInfo{
		BaseName: baseName,
		Ext:      ext,
		Name:     name,
		Dir:      dir,

		OutputPath: name + "." + targetFormat,
		TargetName: targetName + ext,
	}
}

func MustSevenZip() {
	_, err := exec.LookPath("7z")
	if err != nil {
		log.Fatalln("7z not found in PATH")
	}
}

type SevenZip struct {
	Dir string
}

func (z SevenZip) cmd(args ...string) *exec.Cmd {
	cmd := exec.Command("7z", args...)
	cmd.Dir = z.Dir
	return cmd
}

func (z SevenZip) Add(method, dist, src string) error {
	return z.cmd("a", "-t"+method, dist, src, "-mx9").Run()
}

func (z SevenZip) Rename(dist, from, to string) error {
	return z.cmd("rn", dist, from, to).Run()
}

type MakeCompress func(filePath, targetName string) error

func MakeZip(filePath, targetName string) error {
	info := DecodeCompressPath(filePath, targetName, "zip")
	sevenZip := &SevenZip{Dir: info.Dir}

	err := sevenZip.Add("zip", info.OutputPath, info.BaseName)
	if err != nil {
		return err
	}
	return sevenZip.Rename(info.OutputPath, info.BaseName, info.TargetName)
}

func MakeTarGzip(filePath, targetName string) error {
	info := DecodeCompressPath(filePath, targetName, "tar")
	sevenZip := &SevenZip{Dir: info.Dir}

	err := sevenZip.Add("tar", info.OutputPath, info.BaseName)
	if err != nil {
		return err
	}
	defer os.Remove(path.Join(info.Dir, info.OutputPath))

	if err = sevenZip.Rename(info.OutputPath, info.BaseName, info.TargetName); err != nil {
		return err
	}

	return sevenZip.Add("gzip", info.OutputPath+".gz", info.OutputPath)
}
