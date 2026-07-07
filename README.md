# zpl2png

A small pure-Go command-line utility that reads ZPL (Zebra Programming Language) from stdin and writes the rendered label as a PNG to stdout. Rendering is done by [ingridhq/zebrash](https://github.com/ingridhq/zebrash).

## Installation

`zpl2png` will run on most Linux and MacOS systems.

To install it, just `cd` into the directory in which you wish to install it and then copy-paste the appropriate one-liner from below.

### Linux (32-bit)

```
curl -s -L -o zpl2png https://github.com/alasdairmorris/zpl2png/releases/latest/download/zpl2png-linux-386 && chmod +x zpl2png
```

### Linux (64-bit)

```
curl -s -L -o zpl2png https://github.com/alasdairmorris/zpl2png/releases/latest/download/zpl2png-linux-amd64 && chmod +x zpl2png
```

### Mac OS X (Intel)

```
curl -s -L -o zpl2png https://github.com/alasdairmorris/zpl2png/releases/latest/download/zpl2png-darwin-amd64 && chmod +x zpl2png
```

### Mac OS X (Apple Silicon)

```
curl -s -L -o zpl2png https://github.com/alasdairmorris/zpl2png/releases/latest/download/zpl2png-darwin-arm64 && chmod +x zpl2png
```

## Usage

```sh
zpl2png [flags] < input.zpl > output.png
```

A ZPL stream may contain multiple labels (each delimited by `^XA ... ^XZ`). By default the first label is rendered.

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-index` | `0` | Index of the label to render (0-based). `-1` prints the label count to stderr without rendering. |
| `-dpmm` | `8` | Print density in dots per millimetre (6=152dpi, 8=203dpi, 12=300dpi, 24=600dpi). |
| `-width` | `0` | Label width in millimetres (`0` = library default, 101.6mm / 4in). |
| `-height` | `0` | Label height in millimetres (`0` = library default, 203.2mm / 8in). |
| `-grayscale` | `false` | Output 8-bit grayscale PNG preserving anti-aliasing instead of binary monochrome. |
| `-inverted` | `false` | Render inverted-orientation labels upside-down. |
| `-o` | `-` | Output file path (`-` for stdout). |

### Examples

```sh
# Render a label to a file
printf '^XA^FO50,50^A0N,50,50^FDHello World^FS^XZ' | zpl2png > label.png

# 300 dpi, custom 100x150mm label size
zpl2png -dpmm=12 -width=100 -height=150 < label.zpl > label.png

# Count how many labels are in a stream
zpl2png -index=-1 < batch.zpl

# Render the second label in a multi-label stream
zpl2png -index=1 < batch.zpl > second.png
```

### Build From Source

If you have Go installed and would prefer to build the app yourself, you can do:

```
go install github.com/alasdairmorris/zpl2png@latest
```

## License

[MIT](LICENSE)
