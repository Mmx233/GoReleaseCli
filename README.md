# GoReleaseCli

[![Lisense](https://img.shields.io/github/license/Mmx233/GoReleaseCli)](https://github.com/Mmx233/GoReleaseCli/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/Mmx233/GoReleaseCli?color=blueviolet&include_prereleases)](https://github.com/Mmx233/GoReleaseCli/releases)
[![GoReport](https://goreportcard.com/badge/github.com/Mmx233/GoReleaseCli)](https://goreportcard.com/report/github.com/Mmx233/GoReleaseCli)

English | [中文](./README_CN.md)

```shell
~$ release --help-long
usage: release [<flags>] <target>

Golang build production release helper.

Flags:
  -h, --[no-]help          Show context-sensitive help (also try --help-long and
                           --help-man).
  -v, --[no-]version       Show application version.
  -j, --thread=(NumCpu+1)  How many threads to use for parallel compilation.
  -c, --compress=COMPRESS  Compress the binary into the specified format of
                           compressed file.
      --[no-]disable-default-ldflags
                           Disable ldflags added by default.
      --ldflags=LDFLAGS    Add custom ldflags.
      --[no-]soft-float    Enable soft float for mips.
      --[no-]cgo           Enable go cgo.
      --os=OS              Target os
      --arch=ARCH          Target arch.
  -d, --output="build"     Output dir path.
  -o, --name=NAME          Output binary file name.

Args:
  <target>  Target package.
```

## :saxophone: Usage

CGO, soft-float, compression is disabled by default.

By default, compile for all architecture types. You can use the flags `--os` and `--arch` to specify the operating system or architecture, separated by commas. The program will automatically match valid architectures for compilation.

```shell
~$ release ./cmd/release
~$ release ./cmd/release --os linux,windows
~$ release ./cmd/release --arch amd64,386
```

During compilation, default ldflags include `-extldflags "-static -fpic" -s -w` as well as `trimpath`. If additional custom ldflags are needed, you can use an additional flag to append them.

```shell
~$ release ./cmd/release --ldflags='-X main.Version=5.5.5'

~$ release ./cmd/release --disable-default-ldflags # Remove default ldflags.
```

When using `--soft-float`, a soft floating-point version will be added for all MIPS architectures.

```shell
~$ release ./cmd/release --soft-float
```

Compress to the specified format, dependent on the 7z library. Currently supported formats include `zip` and `tar.gz`.

```shell
~$ release  ./cmd/release -c tar.gz
```

By default, the directory name of the target directory will be used, and the compilation result will be placed in the `build` directory. This can also be modified using flags.

```shell
~$ release ./cmd/release --output dist # Modify the output directory to be "dist"
~$ release ./cmd/release -d dist

~$ release ./cmd/release --name asd # Change the name to "asd".
~$ release ./cmd/release -o asd
```

## :factory: Use in GitHub Action

```yaml
name: Release

on:
  push:
    tags:
      - v**
jobs:
  release_docker:
    runs-on: ubuntu-latest
    steps:

  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3

      - name: Build
        uses: Mmx233/GoReleaseCli@v1.1.8
        with:
          target: ./cmd/derper
          compress: tar.gz
          soft-float: true

      - name: Upload assets
        uses: softprops/action-gh-release@v1
        with:
          files: build/*.tar.gz
          prerelease: false
```