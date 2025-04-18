/*
 * file2qr - Convert files to QR codes
 * Copyright (C) 2025 file2qr contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"strings"

	"github.com/skip2/go-qrcode"
)

const VERSION = "1.0.0"

const (
	ProgramName = "file2qr"
)

// displayQRCodeInTerminal creates a temporary PNG file and displays it using terminal graphics
func displayQRCodeInTerminal(qrImage *qrcode.QRCode, size int) error {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "file2qr-*.png")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	// Clean up when done
	defer os.Remove(tmpPath)

	// Write QR code to the temporary file
	err = qrImage.WriteFile(size, tmpPath)
	if err != nil {
		return fmt.Errorf("failed to write QR code to temporary file: %v", err)
	}

	// Open the temporary image file
	f, err := os.Open(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to open temporary QR code file: %v", err)
	}
	defer f.Close()

	// Decode the image
	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode QR code image: %v", err)
	}

	// Get image dimensions
	bounds := img.Bounds()

	// Add some padding around the QR code
	fmt.Println()

	// Loop through the image two rows at a time
	for y := bounds.Min.Y; y < bounds.Max.Y-1; y += 2 {
		// Add some left padding
		fmt.Print("  ")

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the colors of the top and bottom pixels in this column
			top := img.At(x, y)
			bottom := img.At(x, y+1)

			// Convert both colors to 8-bit RGB values
			r1, g1, b1 := rgb(top)
			r2, g2, b2 := rgb(bottom)

			// Set foreground color (top pixel)
			fmt.Printf("\x1b[38;2;%d;%d;%dm", r1, g1, b1)

			// Set background color (bottom pixel)
			fmt.Printf("\x1b[48;2;%d;%d;%dm", r2, g2, b2)

			// Print the "upper half block" character
			fmt.Print("▀")
		}

		// Reset colors and move to the next terminal line
		fmt.Print("\x1b[0m\n")
	}

	// Add some bottom padding
	fmt.Println()

	return nil
}

// rgb converts a color.Color into standard 8-bit RGB values (0–255 per channel)
func rgb(c color.Color) (r, g, b uint8) {
	rr, gg, bb, _ := c.RGBA()                             // Returns 16-bit (0–65535) values
	return uint8(rr >> 8), uint8(gg >> 8), uint8(bb >> 8) // Convert to 8-bit
}

// readFromStdin reads all data from standard input
func readFromStdin() ([]byte, error) {
	return io.ReadAll(os.Stdin)
}

// showUsage prints a brief usage message
func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE]\n", ProgramName)
	fmt.Fprintf(os.Stderr, "Convert files to QR codes.\n\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nIf FILE is not specified, %s reads from standard input.\n", ProgramName)
	fmt.Fprintf(os.Stderr, "If -o/--output is not specified, displays QR code in terminal.\n")
}

// showVersion prints version information
func showVersion() {
	fmt.Printf("%s %s\n", ProgramName, VERSION)
	fmt.Println("Copyright (C) 2025 file2qr contributors")
	fmt.Println("License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>.")
	fmt.Println("This is free software: you are free to change and redistribute it.")
	fmt.Println("There is NO WARRANTY, to the extent permitted by law.")
}

func main() {
	// Define command line flags
	outputFile := flag.String("o", "", "Output QR code file path (PNG format)")
	size := flag.Int("s", 256, "QR code size in pixels")
	recovery := flag.String("r", "medium", "QR code recovery level: low, medium, high, highest")
	terminalSize := flag.Int("t", 40, "Size of QR code when displayed in terminal")
	base64Flag := flag.Bool("b", false, "Base64 encode content (recommended for binary files)")
	versionFlag := flag.Bool("v", false, "Display version information and exit")
	helpFlag := flag.Bool("h", false, "Display this help and exit")

	// Add long options
	flag.StringVar(outputFile, "output", *outputFile, "Output QR code file path (PNG format)")
	flag.IntVar(size, "size", *size, "QR code size in pixels")
	flag.StringVar(recovery, "recovery", *recovery, "QR code recovery level: low, medium, high, highest")
	flag.IntVar(terminalSize, "term-size", *terminalSize, "Size of QR code when displayed in terminal")
	flag.BoolVar(base64Flag, "base64", *base64Flag, "Base64 encode content (recommended for binary files)")
	flag.BoolVar(versionFlag, "version", *versionFlag, "Display version information and exit")
	flag.BoolVar(helpFlag, "help", *helpFlag, "Display this help and exit")

	// Custom usage message
	flag.Usage = showUsage

	// Parse flags
	flag.Parse()

	// Show version and exit if requested
	if *versionFlag {
		showVersion()
		os.Exit(0)
	}

	// Show help and exit if requested
	if *helpFlag {
		showUsage()
		os.Exit(0)
	}

	// Determine input source (file or stdin)
	var inputData []byte
	var err error

	if flag.NArg() > 0 {
		// Use the last positional argument as input file
		inputFile := flag.Arg(flag.NArg() - 1)
		inputData, err = os.ReadFile(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Read from stdin
		inputData, err = readFromStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
	}

	// Base64 encode if requested
	var qrContent string
	if *base64Flag {
		qrContent = base64.StdEncoding.EncodeToString(inputData)
		if isTerminal(os.Stderr.Fd()) {
			fmt.Fprintf(os.Stderr, "Data encoded as Base64 (length: %d characters)\n", len(qrContent))
		}
	} else {
		qrContent = string(inputData)
	}

	// Determine recovery level
	var recLevel qrcode.RecoveryLevel
	switch strings.ToLower(*recovery) {
	case "low":
		recLevel = qrcode.Low
	case "medium":
		recLevel = qrcode.Medium
	case "high":
		recLevel = qrcode.High
	case "highest":
		recLevel = qrcode.Highest
	default:
		recLevel = qrcode.Medium
	}

	// Generate QR code
	qrImage, err := qrcode.New(qrContent, recLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating QR code: %v\n", err)
		if len(qrContent) > 2900 {
			fmt.Fprintf(os.Stderr, "Content size (%d bytes) might be too large for a QR code.\n", len(qrContent))
			fmt.Fprintf(os.Stderr, "Try reducing file size or using the -b/--base64 option for binary files.\n")
		}
		os.Exit(1)
	}

	// Output QR code to file or terminal
	if *outputFile != "" {
		// Output to file
		err = qrImage.WriteFile(*size, *outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving QR code to file: %v\n", err)
			os.Exit(1)
		}
		if isTerminal(os.Stderr.Fd()) {
			fmt.Fprintf(os.Stderr, "QR code saved to: %s\n", *outputFile)
		}
	} else {
		// Output to terminal
		if isTerminal(os.Stdout.Fd()) {
			err := displayQRCodeInTerminal(qrImage, *terminalSize)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error displaying QR code: %v\n", err)
				os.Exit(1)
			}
		} else {
			// If stdout is not a terminal, we can't display the image visually
			fmt.Fprintf(os.Stderr, "Error: Output is not a terminal. Use -o/--output to specify an output file.\n")
			os.Exit(1)
		}
	}
}

// isTerminal determines if the given file descriptor is a terminal
// This is needed to handle both piped input/output and direct terminal usage
func isTerminal(fd uintptr) bool {
	// This is a simplified version. For a real program,
	// you might want to use a package like github.com/mattn/go-isatty
	// or implement proper terminal detection based on the OS
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
