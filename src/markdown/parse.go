package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

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
	// p.skipWhitespace()
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
	// for string(p.char) == " " || string(p.char) == "\t" {
	fmt.Println("in skip func, p.char is", string(p.char))
	for p.char == ' ' || p.char == '\t' || p.char == '\n' || p.char == '\r' {
		fmt.Println("skipping")
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
		output = append(output, convertInlineCode(convertLinks(asTag)))
	}

	return strings.Join(output, "")
}

func convertLinks(md string) string {
	r := regexp.MustCompile(`(\[)([\w\s\.]+)(\])(\()([\.\/:a-zA-Z0-9#_]+)(\))`)
	return r.ReplaceAllString(md, `<a href="$5">$2</a>`)
}

func convertInlineCode(md string) string {
	r := regexp.MustCompile("`([a-zA-Z0-9_@#:\\.\\[\\]\\(\\)\\{\\}]+)`")
	return r.ReplaceAllString(md, `<code class="inline">$1</code>`)
}

func convertElement(md string) string {
	hTag := regexp.MustCompile(`(#+\s*)([\w\s\(\)\.]+)`)
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
		code := regexp.MustCompile("(`{3})([a-z]+)\\n([\\sa-zA-Z/=\"#:_\\.,@\\[\\]\\(\\)\\{\\}]+)(`{3})")
		return code.ReplaceAllString(md, `<pre><code class="$2">$3</code></pre>`)
	}
	// ```elixir def change do create_if_not_exists table(:avocados) do add :date_picked, :utc_datetime # etc... ```
	return fmt.Sprintf("<p>%s</p>", md)
}
