package md2man

import (
	"github.com/govenue/encoding/markdown"
)

func Render(doc []byte) []byte {
	renderer := RoffRenderer(0)
	extensions := 0
	extensions |= markdown.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= markdown.EXTENSION_TABLES
	extensions |= markdown.EXTENSION_FENCED_CODE
	extensions |= markdown.EXTENSION_AUTOLINK
	extensions |= markdown.EXTENSION_SPACE_HEADERS
	extensions |= markdown.EXTENSION_FOOTNOTES
	extensions |= markdown.EXTENSION_TITLEBLOCK

	return markdown.Markdown(doc, renderer, extensions)
}
