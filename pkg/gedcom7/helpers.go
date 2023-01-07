package gedcom7

import (
	"bytes"
)

//	func skipYAMLLine(s string) bool {
//		if len(s) > 2 {
//			if s[:3] == "---" || s[:3] == "..." {
//				return true
//			}
//		}
//
//		return false
//	}

func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		if data[i] == '\n' {
			// We have a line terminated by single newline.
			return i + 1, data[0:i], nil
		}
		// We have a line terminated by carriage return at the end of the buffer.
		if !atEOF && len(data) == i+1 {
			return 0, nil, nil
		}
		advance = i + 1
		if len(data) > i+1 && data[i+1] == '\n' {
			advance += 1
		}
		return advance, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
