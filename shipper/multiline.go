package shipper

import (
	"bytes"
	"regexp"
)

type MultiLineHandlerMode int

const (
	MultiLineHandlerModeBeginsWith MultiLineHandlerMode = iota
	MultiLineHandlerModeContinuesWith
)

type MultiLineHandler struct {
	LineHandler LineHandler
	Regex       *regexp.Regexp
	Mode        MultiLineHandlerMode

	buf bytes.Buffer
}

func (h *MultiLineHandler) HandleLine(line string) {
	isMatch := h.Regex.MatchString(line)

	switch h.Mode {
	case MultiLineHandlerModeBeginsWith:
		if isMatch {
			fullLine := h.buf.String()
			h.buf.Reset()
			h.LineHandler.HandleLine(fullLine)
		}

	case MultiLineHandlerModeContinuesWith:
		if !isMatch {
			fullLine := h.buf.String()
			h.buf.Reset()
			h.LineHandler.HandleLine(fullLine)
		}
	}
	h.buf.WriteString(line)
	h.buf.WriteString("\n")
}
