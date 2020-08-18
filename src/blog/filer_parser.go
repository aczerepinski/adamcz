package blog

import (
	"strings"
)

// DELIMITER should be placed immediately before and offer a key,
// such as *** title ***
const DELIMITER string = "***"

type parser struct {
	input      string
	position   int
	char       byte
	currentKey string
	parsed     map[string]string
}

func newParser(file []byte) *parser {
	return &parser{
		input:  string(file),
		parsed: map[string]string{},
	}
}

// ParseFile converts a file of a particular format into
// a map of key value pairs. An example:
//
// *** title ***
// This is the title
//
// The above will be converted into {"title": "This is the title"}
func ParseFile(file []byte) map[string]string {
	p := newParser(file)
	for p.position < len(p.input) {
		p.readKey()
		p.readValue()
	}
	return p.parsed
}

func (p *parser) readKey() {
	p.readDelimiter()
	p.currentKey = p.readText()
	p.readDelimiter()
}

func (p *parser) readValue() {
	p.parsed[p.currentKey] = p.readText()
}

func (p *parser) readText() string {
	position := p.position
	for p.position < len(p.input) && (p.almostDone() || !p.isDelimiter()) {
		p.readChar()
	}
	return strings.TrimSpace(p.input[position:p.position])
}

func (p *parser) isDelimiter() bool {
	peekUntil := p.position + len(DELIMITER)
	return string(p.input[p.position:peekUntil]) == DELIMITER
}

func (p *parser) almostDone() bool {
	return len(p.input)-p.position < len(DELIMITER)
}

func (p *parser) readDelimiter() {
	p.position = p.position + len(DELIMITER)
}

func (p *parser) readChar() {
	if p.position == len(p.input) {
		p.char = 0
	} else {
		p.char = p.input[p.position]
	}

	if p.almostDone() || !p.isDelimiter() {
		p.position++
	} else {
		p.position = p.position + len(DELIMITER)
	}

}
