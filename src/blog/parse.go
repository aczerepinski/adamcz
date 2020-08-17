package blog

import (
	"fmt"
	"strings"
)

// DELIMITER should be placed immediately before and offer a key,
// such as *** Title ***
const DELIMITER string = "***"

// Parser parses blog & music posts and is heavily influenced &
// informed by Thorsten Ball's wonderful book, Writing an
// Interpreter in Go.
type Parser struct {
	input    string
	position int
	// readPosition int
	char       byte
	currentKey string
	parsed     map[string]string
}

// Parse converts a file of a particular format into
// a map of key value pairs. An example:
//
// *** Title ***
// This is the title
//
// The above will be confirted into {"title": "This is the title"}
func (p *Parser) Parse(file []byte) map[string]string {
	p.input = string(file)
	p.parsed = map[string]string{}
	for p.position < len(p.input) {
		p.readKey()
		p.readValue()
	}
	return p.parsed
}

func (p *Parser) readKey() {
	p.readDelimiter()
	p.currentKey = p.readText()
	p.readDelimiter()
}

func (p *Parser) readValue() {
	p.skipWhitespace()
	p.parsed[p.currentKey] = strings.TrimSpace(p.readText())
	fmt.Printf("parsed: \n\n%+v\n\n", p.parsed)
}

func (p *Parser) readText() string {
	p.skipWhitespace()
	position := p.position
	for !p.isDelimiter() {
		p.readChar()
	}
	return p.input[position:p.position]
}

func (p *Parser) isDelimiter() bool {
	if p.position+len(DELIMITER) > len(p.input) {
		fmt.Printf("\n%+v\n\n", p)
		panic("out of bounds bug")
	}

	peekUntil := p.position + len(DELIMITER)
	return string(p.input[p.position:peekUntil]) == DELIMITER
}

func (p *Parser) readDelimiter() {
	p.skipWhitespace()
	// for !p.isDelimiter() {
	// 	// if string(p.char) == "*" {
	// 	// 	fmt.Printf("previous: %s, current: %s, next: %s\n",
	// 	// 		string(p.input[p.position-1]), string(p.input[p.position]), string(p.input[p.position+1]))
	// 	// 	// fmt.Printf("%+v\n", p)
	// 	// 	panic("bug")
	// 	// }
	// 	p.position++
	// 	// p.readPosition++
	// }
	p.position = p.position + len(DELIMITER)
	// p.readPosition = p.readPosition + len(DELIMITER) + 1
}

func (p *Parser) readChar() {
	if p.position >= len(p.input) {
		p.char = 0
	} else {
		p.char = p.input[p.position]
	}
	if len(p.input)-p.position > len(DELIMITER) || !p.isDelimiter() {
		p.position++
	} else {
		p.position = p.position + len(DELIMITER)
	}

}

func (p *Parser) skipWhitespace() {
	for p.char == ' ' || p.char == '\t' || p.char == '\n' || p.char == '\r' {
		p.readChar()
	}
}
