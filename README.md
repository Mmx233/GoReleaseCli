# GoReleaseCli

[![Lisense](https://img.shields.io/github/license/Mmx233/GoReleaseCli)](https://github.com/Mmx233/GoReleaseCli/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/Mmx233/GoReleaseCli?color=blueviolet&include_prereleases)](https://github.com/Mmx233/GoReleaseCli/releases)
[![GoReport](https://goreportcard.com/badge/github.com/Mmx233/GoReleaseCli)](https://goreportcard.com/report/github.com/Mmx233/GoReleaseCli)

```shell
~$ release --help-long
usage: release [<flags>] <target>

Golang build production release helper.

Flags:
  -h, --[no-]help         Show context-sensitive help (also try --help-long and
                          --help-man).
  -v, --[no-]version      Show application version.
  -j, --thread=(NumCpu+1) How many threads to use for parallel compilation.
  -c, --compress=COMPRESS  Compress the binary into the specified format of
                           compressed file.
      --ldflags=LDFLAGS   Add custom ldflags.
      --[no-]soft-float   Enable soft float for mips.
      --[no-]cgo          Enable go cgo.
      --os=OS             Target os.
      --arch=ARCH         Target arch.
  -d, --output="build"    Output dir path.
  -o, --name=NAME         Output binary file name.

Args:
  <target>  target package
```

## :saxophone: 使用

CGO、软浮点、生成压缩文件默认关闭

默认编译所有架构类型，可以使用 flag `--os` 和 `--arch` 指定系统或架构，使用英文逗号分隔。程序会自动匹配出有效架构进行编译

```shell
~$ release ./cmd/release
~$ release ./cmd/release --os linux,windows
~$ release ./cmd/release --arch amd64,386
```

编译时已带有 `-extldflags "-static -fpic" -s -w` 以及 `trimpath` 的 ldflags，如果需要附加自定义 ldflags，可以用 flag 继续加

```shell
~$ release ./cmd/release --ldflags='-X main.Version=5.5.5'
```

当使用 --soft-float 时，会为所有 mips 架构添加软浮点版本

```shell
~$ release ./cmd/release --soft-float
```

压缩到指定格式，依赖 7z lib，目前支持 `zip` `tar.gz`

```shell
~$ release  ./cmd/release -c tar.gz
```

默认情况下，会使用 target 目录的目录名，编译结果放在 build 目录下，这也是可以通过 flag 修改的

```shell
~$ release ./cmd/release --output dist # 修改输出目录为 dist
~$ release ./cmd/release -d dist

~$ release ./cmd/release --name asd # 修改名称为 asd
~$ release ./cmd/release -o asd
```