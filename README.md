# GoReleaseCli

```shell
~$ release --help-long
usage: release [<flags>] <target>

Golang build production release helper.

Flags:
  -h, --[no-]help        Show context-sensitive help (also try --help-long and
                         --help-man).
  -v, --[no-]version     Show application version.
      --ldflags=LDFLAGS  add custom ldflags
      --[no-]soft-float  enable soft float for mips
      --os=OS            target os
      --arch=ARCH        target arch
  -d, --output="build"   output dir path
  -o, --name=NAME        output binary file name

Args:
  <target>  target package
```

## :saxophone: 使用

默认编译所有架构类型，可以使用 flag `--os` 和 `--arch` 指定系统或架构，使用英文逗号分隔。程序会自动匹配出有效架构进行编译

```shell
~$ release ./cmd/release
~$ release ./cmd/release --os linux,windows
~$ release ./cmd/release --arch amd64,386
```

编译时已带有 `-extldflags "-static" -s -w` 的 ldflags，如果需要附加自定义 ldflags，可以用 flag 继续加

```shell
~$ release ./cmd/release --ldflags='-X main.Version=5.5.5'
```

当使用 --soft-float 时，会为所有 mips 架构添加软浮点版本

```shell
~$ release ./cmd/release --soft-float
```

默认情况下，会使用 target 目录的目录名，编译结果放在 build 目录下，这也是可以通过 flag 修改的

```shell
~$ release ./cmd/release --output dist # 修改输出目录为 dist
~$ release ./cmd/release -d dist

~$ release ./cmd/release --name asd # 修改名称为 asd
~$ release ./cmd/release -o asd
```