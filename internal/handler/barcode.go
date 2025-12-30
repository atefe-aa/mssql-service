package handler

import "strings"

func decodeBarcode(barcode string) string {
	parts := strings.SplitN(barcode, "z", 2)

	if len(parts) != 2 {
		return barcode
	}

	prefix := parts[0]
	suffix := parts[1]

	hexMap := map[byte]string{
		'A': "10", 'a': "10",
		'B': "11", 'b': "11",
		'C': "12", 'c': "12",
		'D': "13", 'd': "13",
		'E': "14", 'e': "14",
		'F': "15", 'f': "15",
		'X': "99", 'x': "99",
	}

	var decodedPrefix strings.Builder
	for i := 0; i < len(prefix); i++ {
		char := prefix[i]
		if replacement, exists := hexMap[char]; exists {
			decodedPrefix.WriteString(replacement)
		} else {
			decodedPrefix.WriteByte(char)
		}
	}

	decodedPrefixStr := decodedPrefix.String()
	var formattedPrefix string

	if len(decodedPrefixStr) > 1 {
		firstPart := decodedPrefixStr[0:1]
		secondPart := decodedPrefixStr[1:]
		formattedPrefix = firstPart + "." + secondPart
	} else {
		formattedPrefix = decodedPrefixStr
	}

	return formattedPrefix + "-" + suffix
}
