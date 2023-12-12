package withnohead

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/nohead"
	"fyne.io/fyne/v2/test"
)

// AssertHeadlessObjectRendersToImage asserts that the given `CanvasObject` renders the same image as the one stored in the master file.
// The theme used is the standard test theme which may look different to how it shows on your device.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the given image is not equal to the loaded master image.
// In this case the given image is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
//
// Since 2.3
func AssertHeadlessObjectRendersToImage(t *testing.T, masterFilename string, o fyne.CanvasObject, msgAndArgs ...interface{}) bool {
	c := nohead.NewHeadlessCanvasWithPainter(nohead.NewHeadlessPainter())
	c.SetPadded(false)
	c.SetContent(o)

	return test.AssertRendersToImage(t, masterFilename, c, msgAndArgs...)
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
func AssertHeadlessObjectRendersToMarkup(t *testing.T, masterFilename string, o fyne.CanvasObject, msgAndArgs ...interface{}) bool {
	c := nohead.GetInMemoryHeadlessCanvas()
	c.SetContent(o)
	c.EnsureMinSize()

	return test.AssertRendersToMarkup(t, masterFilename, c, msgAndArgs...)
}
