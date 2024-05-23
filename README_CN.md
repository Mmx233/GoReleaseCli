# GoReleaseCli

[![Lisense](https://img.shields.io/github/license/Mmx233/GoReleaseCli)](https://github.com/Mmx233/GoReleaseCli/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/Mmx233/GoReleaseCli?color=blueviolet&include_prereleases)](https://github.com/Mmx233/GoReleaseCli/releases)
[![GoReport](https://goreportcard.com/badge/github.com/Mmx233/GoReleaseCli)](https://goreportcard.com/report/github.com/Mmx233/GoReleaseCli)

[English](./README.md) | 中文

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
      --[no-]extra-arches  Build all extra arches.
      --[no-]extra-arches-show-default
                           Show default extra arch name.
  -d, --output="build"     Output dir path.
  -o, --name=NAME          Output binary file name.

Args:
  <target>  Target package.
```

## :saxophone: 使用

CGO、软浮点、生成压缩文件默认关闭

默认编译所有架构类型，可以使用 flag `--os` 和 `--arch` 指定系统或架构，使用英文逗号分隔。程序会自动匹配出有效架构进行编译

```shell
~$ release ./cmd/release
~$ release ./cmd/release --os linux,windows
~$ release ./cmd/release --arch amd64,386
```

编译时默认已带有 `-extldflags "-static -fpic" -s -w` 以及 `-trimpath` 的 ldflags，如果需要附加自定义 ldflags，可以用 flag 继续加

```shell
~$ release ./cmd/release --ldflags='-X main.Version=5.5.5'

~$ release ./cmd/release --disable-default-ldflags # 移除默认 ldflags
```

当使用 `--extra-arches` 时，会编译出所有子架构如 arm/v6 arm/v7

默认情况下，默认架构编译结果的名称不会添加额外的子架构后缀。但你可以通过 `--extra-arches-show-default` 启用它

```shell
~$ release ./cmd/release --extra-arches

~$ release ./cmd/release --extra-arches --extra-arches-show-default # 为默认架构添加子架构后缀
```

压缩到指定格式，依赖 `7z` lib，没有 `7z` 时会尝试使用 `zip` + `zipnote` / `tar` 分别为不同压缩类型压缩，目前支持 `zip` `tar.gz`

```shell
~$ release  ./cmd/release -c tar.gz
```

默认情况下，会使用 target 目录的目录名，编译结果放在 `build` 目录下，这也是可以通过 flag 修改的

```shell
~$ release ./cmd/release --output dist # 修改输出目录为 dist
~$ release ./cmd/release -d dist

~$ release ./cmd/release --name asd # 修改名称为 asd
~$ release ./cmd/release -o asd
```

## :factory: 在 GitHub Action 中使用

### 在容器中构建

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
        uses: actions/checkout@v4

      - name: Build
        uses: Mmx233/GoReleaseCli@v1.2.2
        with:
          target: ./cmd/derper
          compress: tar.gz

      - name: Upload assets
        uses: softprops/action-gh-release@v1
        with:
          files: build/*.tar.gz
          prerelease: false
```

### 在 Action Runner 环境中构建

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
        uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: 'stable'

      - name: Setup Release Cli
        uses: robinraju/release-downloader@v1.10
        with:
          repository: "Mmx233/GoReleaseCli"
          latest: true
          fileName: 'release_linux_amd64.tar.gz'
          extract: true
          out-file-path: './build/'

      - name: Build
        run: ./build/release ./cmd/derper --perm 777 -c tar.gz --extra-arches --output build/output

      - name: Upload assets
        uses: softprops/action-gh-release@v1
        with:
          files: build/output/*.tar.gz
          prerelease: false
```