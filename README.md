# ature

[![CI](https://github.com/rpearce/ature/actions/workflows/ci.yml/badge.svg)](https://github.com/rpearce/ature/actions/workflows/ci.yml)

Temperature conversion CLI tool.

_This is a tool for me to learn Go, so you probably shouldn't use this._

## Install

```
go install github.com/rpearce/ature@latest
```

Prebuilt binaries for Linux, macOS, and Windows are attached to each
[release](https://github.com/rpearce/ature/releases).

`ature --version` prints the installed version.

## Example

```
λ ature ctof 16
60.8

λ ature ktof 273.15
32

λ ature ktof 273.15 | xargs ature ftoc | xargs ature ctof
32

λ ature ftoc -- -10
-23.333333
```

Results are rounded to six decimal places.

Note that because of unix patterns, if you provide a negative number, you must
first signify the end of command/flag arguments with `--` and then provide the
negative value.

## All commands

```
λ ature -h
Convert temperature values

Usage:
  ature [flags]
  ature [command]

Available Commands:
  ctof        Convert Celsius to Fahrenheit
  ctok        Convert Celsius to Kelvin
  ftoc        Convert Fahrenheit to Celsius
  ftok        Convert Fahrenheit to Kelvin
  help        Help about any command
  ktoc        Convert Kelvin to Celsius
  ktof        Convert Kelvin to Fahrenheit

Flags:
  -h, --help      help for ature
  -v, --version   version for ature

Use "ature [command] --help" for more information about a command.
```

## Development

[mise](https://mise.jdx.dev) pins the Go toolchain and tools. After
`mise install`:

```
mise run test      # go vet and go test -race
mise run lint      # golangci-lint
mise run fmt       # gofumpt and goimports
mise run snapshot  # local goreleaser build into dist/
```

Pushing a `v*` tag runs goreleaser in GitHub Actions and publishes a
release.
