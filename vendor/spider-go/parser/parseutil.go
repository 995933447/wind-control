package parser

import "regexp"

func ExtractString(re *regexp.Regexp, content []byte) string {
	matches := re.FindSubmatch(content)
	if len(matches) > 1 {
		return string(matches[1])
	} else {
		return ""
	}
}

func ParseMatchesToString(matches [][][]byte) []string {
	var result []string
	for _, match := range matches {
		result = append(result, string(match[len(match) - 1]))
	}
	return result
}

func ParseMatchesToHtml(matches [][][]byte) [][]byte {
	var result [][]byte
	for _, match := range matches {
		result = append(result, match[0])
	}
	return result
}
