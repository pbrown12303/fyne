package nohead

import (
	"fyne.io/fyne/v2"
)

// Declare conformity with Clipboard interface
var _ fyne.Clipboard = (*headlessClipboard)(nil)

// headlessClipboard represents the system headlessClipboard
type headlessClipboard struct {
	headlessWindow *headlessWindow
	content        string
}

// Content returns the clipboard content
func (c *headlessClipboard) Content() string {
	return c.content
}

// SetContent sets the clipboard content
func (c *headlessClipboard) SetContent(content string) {
	c.content = content
}
