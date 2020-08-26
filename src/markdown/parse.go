package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

const codeChars string = "[\\sa-zA-Z0-9/=~<>\\\\%\\-\\+\\*\\^`&?\\$\\|!'\"#:;_\\.,@\\[\\]\\(\\)\\{\\}]+"

type parser struct {
	char      byte
	input     string
	position  int
	currentEl string
	parsed    []string
}

func (p *parser) parse(md string) []string {
	p.input = md
	p.char = 1
	for p.position <= len(p.input) {
		p.readElement()
	}
	return p.parsed
}

func (p *parser) readElement() {
	start := p.position

	if start+2 < len(p.input) && p.isCodeDelimiter() {
		p.position = p.position + 3
		for !p.isCodeDelimiter() && p.char != 0 {
			p.readChar()
		}
		p.position = p.position + 3

		p.parsed = append(p.parsed, strings.TrimSpace(p.input[start:p.position]))
	} else {
		for !p.isParagraphDelimiter() && p.char != 0 {
			p.readChar()
		}
		if p.position+2 >= len(p.input) {
			p.parsed = append(p.parsed, strings.TrimSpace(p.input[start:]))
		} else {
			p.parsed = append(p.parsed, strings.TrimSpace(p.input[start:p.position]))
		}
		p.position = p.position + 2
	}
}

func (p *parser) isCodeDelimiter() bool {
	return p.input[p.position] == '`' &&
		p.input[p.position+1] == '`' &&
		p.input[p.position+2] == '`'
}

func (p *parser) isParagraphDelimiter() bool {
	if p.position+1 >= len(p.input) {
		return true
	}
	return p.input[p.position] == '\n' && p.input[p.position+1] == '\n'
}

func (p *parser) skipWhitespace() {
	for p.char == ' ' || p.char == '\t' || p.char == '\n' || p.char == '\r' {
		p.readChar()
	}
}

func (p *parser) readChar() {
	if p.position > len(p.input) {
		p.char = 0
	} else {
		p.char = p.input[p.position]
	}
	p.position++
}

// ToHTML accepts markdown and returns HTML, implementing
// only a subset of the common markdown spec. Links
// are currently supported.
func ToHTML(md string) string {
	var output []string
	p := &parser{}
	elements := p.parse(md)

	for _, e := range elements {
		asTag := convertElement(e)
		if strings.HasPrefix(e, "```") {
			output = append(output, asTag)
		} else {
			output = append(output, convertInlineCode(convertLinks(asTag)))
		}
	}

	return strings.Join(output, "")
}

func convertLinks(md string) string {
	r := regexp.MustCompile(`(\[)([\w\s\-\.]+)(\])(\()([\.\/:a-zA-Z0-9#_\-]+)(\))`)
	return r.ReplaceAllString(md, `<a href="$5">$2</a>`)
}

func convertInlineCode(md string) string {
	inlineChars := strings.Replace(codeChars, "`", "", -1)
	r := regexp.MustCompile(fmt.Sprintf("`(%s)`", inlineChars))
	return r.ReplaceAllString(md, `<code class="inline">$1</code>`)
}

func convertElement(md string) string {
	hTag := regexp.MustCompile(`(#+\s*)([\w\s\(\)\.!:?\\+&]+)`)
	if strings.HasPrefix(md, "###") {
		return hTag.ReplaceAllString(md, "<h3>$2</h3>")
	}
	if strings.HasPrefix(md, "##") {
		return hTag.ReplaceAllString(md, "<h2>$2</h2>")
	}
	if strings.HasPrefix(md, "#") {
		return hTag.ReplaceAllString(md, "<h1>$2</h1>")
	}
	if strings.HasPrefix(md, "```") {
		code := regexp.MustCompile(fmt.Sprintf("(`{3})([a-z]+)\\n(%s)(`{3})", codeChars))
		matches := code.FindStringSubmatch(md)
		if len(matches) < 4 {
			return md
		}
		content := strings.ReplaceAll(matches[3], "<", "&lt;")
		return `<pre><code class="` + matches[2] + `">` + content + "</code></pre>"
	}
	return fmt.Sprintf("<p>%s</p>", md)
}
