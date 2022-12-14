# ature

Temperature conversion CLI tool.

_This is a tool for me to learn Go, so you probably shouldn't use this._

## Example

```
λ ature ctof 16
60.8

λ ature ktof 273.15
32

λ ature ftoc -- -10
-23.333333333333336
```

Note that because of unix patterns, if you provide a negative number, you must
first signify the end of command/flag arguments with `--` and then provide the
negative value.

## All commands

```
λ ature -h
Convert temperature values

Usage:
  ature [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  ctof        Convert Celsius to Fahrenheit
  ctok        Convert Celsius to Kelvin
  ftoc        Convert Fahrenheit to Celsius
  ftok        Convert Fahrenheit to Kelvin
  help        Help about any command
  ktoc        Convert Kelvin to Celsius
  ktof        Convert Kelvin to Fahrenheit

Flags:
  -h, --help   help for ature

Use "ature [command] --help" for more information about a command.
```
