package tools

import (
	"archive/zip"
	"compress/flate"
	"io"
	"os"
	"path"
	"strings"
)

type CompressInfo struct {
	// src info
	BaseName string // filename without parent path
	Ext      string
	Name     string // base name without ext
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

		OutputPath: path.Join(dir, name+"."+targetFormat),
		TargetName: targetName + ext,
	}
}

func MakeZip(filePath, targetName string) error {
	info := DecodeCompressPath(filePath, targetName, "zip")

	file, err := os.Create(info.OutputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	compressFileInfo, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	compressFileInfo.Name = info.TargetName

	archive := zip.NewWriter(file)
	defer archive.Close()
	archive.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	writer, err := archive.CreateHeader(compressFileInfo)
	if err != nil {
		return err
	}

	src, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(writer, src)
	return err
}
