package withnohead

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/nohead"
	"fyne.io/fyne/v2/test"
)

var (
	dummyHeadlessCanvas *nohead.HeadlessCanvas
)

// AssertHeadlessObjectRendersToImage asserts that the given `CanvasObject` renders the same image as the one stored in the master file.
// The theme used is the standard test theme which may look different to how it shows on your device.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the given image is not equal to the loaded master image.
// In this case the given image is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
//
// Since 2.3
func AssertHeadlessObjectRendersToImage(t *testing.T, masterFilename string, o fyne.CanvasObject, c fyne.Canvas, msgAndArgs ...interface{}) bool {
	c.SetContent(o)
	switch typedCanvas := c.(type) {
	case *nohead.HeadlessCanvas:
		typedCanvas.EnsureMinSize()
	}

	return test.AssertRendersToImage(t, masterFilename, c, msgAndArgs...)
}

// AssertRendersToImage calls c.EnsureMinSize then calls the standart test.AssertRendersToImage
func AssertRendersToImage(t *testing.T, masterFilename string, c fyne.Canvas, msgAndArgs ...interface{}) bool {
	switch typedCanvas := c.(type) {
	case *nohead.HeadlessCanvas:
		typedCanvas.EnsureMinSize()
	}
	return test.AssertRendersToImage(t, masterFilename, c, msgAndArgs...)
}

// AssertRendersToMarkup calls c.EnsureMinSize then calls the standard test.AssertRendersToMarkup
func AssertRendersToMarkup(t *testing.T, masterFilename string, c fyne.Canvas, msgAndArgs ...interface{}) bool {
	switch typedCanvas := c.(type) {
	case *nohead.HeadlessCanvas:
		typedCanvas.EnsureMinSize()
	}
	return test.AssertRendersToMarkup(t, masterFilename, c, msgAndArgs...)
}

// AssertHeadlessObjectRendersToMarkup asserts that the given `CanvasObject` renders the same markup as the one stored in the master file.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the rendered markup is not equal to the loaded master markup.
// In this case the rendered markup is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
//
// Be aware, that the indentation has to use tab characters ('\t') instead of spaces.
// Every element starts on a new line indented one more than its parent.
// Closing elements stand on their own line, too, using the same indentation as the opening element.
// The only exception to this are text elements which do not contain line breaks unless the text includes them.
//
// Since 2.3
func AssertHeadlessObjectRendersToMarkup(t *testing.T, masterFilename string, o fyne.CanvasObject, c fyne.Canvas, msgAndArgs ...interface{}) bool {
	c.SetContent(o)
	switch typedCanvas := c.(type) {
	case *nohead.HeadlessCanvas:
		typedCanvas.EnsureMinSize()
	}

	return test.AssertRendersToMarkup(t, masterFilename, c, msgAndArgs...)
}

// GetInMemoryHeadlessCanvas returns a reusable in-memory canvas used for testing
func GetInMemoryHeadlessCanvas() *nohead.HeadlessCanvas {
	if dummyHeadlessCanvas == nil {
		dummyHeadlessCanvas = nohead.NewHeadlessCanvas()
	}

	return dummyHeadlessCanvas
}

// NewHeadlessApp returns an initialized headless app
func NewHeadlessApp() *nohead.HeadlessApp {
	return nohead.NewHeadlessApp()
}

// NewHeadlessCanvas returns an initialized headless canvas
func NewHeadlessCanvas() *nohead.HeadlessCanvas {
	return nohead.NewHeadlessCanvas()
}
