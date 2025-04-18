# file2qr

`file2qr` is a command-line utility written in Go that converts files to QR codes.

## Features

- Convert any file (text or binary) to a QR code
- Read from files or standard input (stdin)
- Output to PNG files or display directly in the terminal
- Support for Base64 encoding of binary data
- Configurable QR code size and error correction level
- Follows standard Unix/Linux command-line interface conventions
- Licensed under GNU GPL v3

## Installation

### Prerequisites

- Go 1.16 or newer

### Building from Source

```bash
# Clone the repository
git clone https://github.com/username/file2qr.git
cd file2qr

# Install dependencies
go get -u github.com/skip2/go-qrcode

# Build the application
go build -o file2qr

# Optional: Install system-wide
sudo cp file2qr /usr/local/bin/
sudo cp file2qr.1 /usr/local/share/man/man1/
```

## Usage

```
file2qr [OPTIONS] [FILE]
```

If `FILE` is not specified, `file2qr` reads from standard input.

### Options

- `-o, --output FILE` : Write QR code to FILE (PNG format)
- `-s, --size PIXELS` : Set QR code image size in pixels (default: 256)
- `-r, --recovery LEVEL` : Set error correction level (low, medium, high, highest)
- `-t, --term-size SIZE` : Size of QR code in terminal (default: 40)
- `-b, --base64` : Encode content using Base64 (recommended for binary files)
- `-h, --help` : Display help message
- `-v, --version` : Display version information

### Examples

Display a QR code for a text file in the terminal:
```bash
file2qr message.txt
```

Generate a QR code PNG for a binary file:
```bash
file2qr -b -o archive-qr.png archive.zip
```

Use in a pipeline (stdin):
```bash
cat config.json | file2qr
```

Use different error correction level:
```bash
file2qr -r high -o secure-qr.png passwords.txt
```

## Terminal Requirements

To display QR codes in the terminal, you need:
- A terminal with support for 24-bit color (true color)
- A font that includes Unicode block characters

Most modern terminal emulators meet these requirements.

## Capacity Limits

QR codes have inherent capacity limitations:
- Regular text data: approximately 2900 characters maximum
- Binary data with Base64: approximately 2100 bytes maximum

For larger files, consider compression before encoding.

## License

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgements

- [go-qrcode](https://github.com/skip2/go-qrcode) - QR code generation library