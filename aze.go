package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"
)

// AZ encoder map
const az = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Map characters to smileys
var smileyMap = map[byte]string{
	'A': "😀", 'B': "😁", 'C': "😂", 'D': "🤣", 'E': "😃",
	'F': "😄", 'G': "😅", 'H': "😆", 'I': "😇", 'J': "😮",
	'K': "😉", 'L': "😊", 'M': "😋", 'N': "😌", 'O': "😍",
	'P': "😎", 'Q': "😏", 'R': "😐", 'S': "😑", 'T': "😒",
	'U': "😓", 'V': "😔", 'W': "😕", 'X': "😖", 'Y': "😗",
	'Z': "😘", 
}

// Reverse mapping for decoding
var reverseSmileyMap = make(map[string]byte)

func init() {
	for k, v := range smileyMap {
		reverseSmileyMap[v] = k
	}
}

func encodeAZ(data string, lineLength int) string {
	encodedData := ""
	smileysOnLine := 0

	for i := 0; i < len(data); {
		char := data[i]
		if strings.ContainsRune(az, rune(char)) {
			encodedData += smileyMap[char]
			smileysOnLine++
			if smileysOnLine == lineLength {
				encodedData += "\n"
				smileysOnLine = 0
			}
			i++
		} else {
			// Handle multi-byte Unicode characters
			r, size := utf8.DecodeRuneInString(data[i:])
			if r != utf8.RuneError && size > 0 {
				encodedData += data[i : i+size]
				i += size
			} else {
				encodedData += string(char)
				i++
			}
		}
	}
	return encodedData
}

func decodeAZ(encodedData string) string {
	decodedData := ""
	currentLine := ""

	for i := 0; i < len(encodedData); {
		char := encodedData[i]
		if strings.Contains(smileyMapString(), string(char)) {
			currentLine += string(char)
			i++
		} else {
			// Handle multi-byte Unicode characters
			r, size := utf8.DecodeRuneInString(encodedData[i:])
			if r != utf8.RuneError && size > 0 {
				currentLine += encodedData[i : i+size]
				i += size
			} else {
				currentLine += string(char)
				i++
			}
		}

        // Decode any remaining characters
        decodedData += decodeSmileyLine(currentLine)
        currentLine = ""
    }

    return decodedData
}

func decodeSmileyLine(line string) string {
    decodedLine := ""
    for _, char := range line {
        if decodedChar, found := reverseSmileyMap[string(char)]; found {
            decodedLine += string(decodedChar)
        } else {
            decodedLine += string(char)
        }
    }
    return decodedLine
}

func smileyMapString() string {
    var s strings.Builder
    for char := range smileyMap {
        s.WriteString(string(char))
    }
    return s.String()
}

func main() {
    decodeFlag := flag.Bool("d", false, "Decode using smiley encoding")
    lineLengthFlag := flag.Int("l", 32, "Set the line length for encoding")
    flag.Parse()

    if *decodeFlag {
        // Decoding mode
        decodedData, err := ioutil.ReadAll(os.Stdin)
        if err != nil {
            fmt.Println("Error reading from stdin:", err)
            os.Exit(1)
        }

        decodedText := decodeAZ(string(decodedData))
        fmt.Print(decodedText)
    } else {
        // Encoding mode
        inputData, err := ioutil.ReadAll(os.Stdin)
        if err != nil {
            fmt.Println("Error reading from stdin:", err)
            os.Exit(1)
        }

        encodedText := encodeAZ(string(inputData), *lineLengthFlag)
        fmt.Print(encodedText)
    }
}

