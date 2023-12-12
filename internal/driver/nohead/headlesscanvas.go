package nohead

import (
	"image"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	dummyHeadlessCanvas *HeadlessCanvas
)

// Declare conformity with Canvas interface
var _ fyne.Canvas = (*HeadlessCanvas)(nil)

// HeadlessCanvas is a canvas in a headless application - it
type HeadlessCanvas struct {
	common.Canvas
	renderedImage image.Image

	content  fyne.CanvasObject
	menu     fyne.CanvasObject
	debug    bool
	padded   bool
	size     fyne.Size
	scale    float32
	texScale float32

	transparent bool

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)
}

// GetInMemoryHeadlessCanvas returns a reusable in-memory canvas used for testing
func GetInMemoryHeadlessCanvas() *HeadlessCanvas {
	if dummyHeadlessCanvas == nil {
		dummyHeadlessCanvas = NewHeadlessCanvas()
	}

	return dummyHeadlessCanvas
}

// NewHeadlessCanvas returns a single use in-memory canvas used for testing.
// This canvas has no painter so calls to Capture() will return a blank image.
func NewHeadlessCanvas() *HeadlessCanvas {
	c := &HeadlessCanvas{
		scale: 1.0,
		size:  fyne.NewSize(10, 10),
	}
	c.SetPadded(true)
	c.Initialize(c, nil)
	c.setContent(&canvas.Rectangle{FillColor: theme.BackgroundColor()})
	return c
}

// NewHeadlessCanvasWithPainter allows creation of an in-memory canvas with a specific painter.
// The painter will be used to render in the Capture() call.
func NewHeadlessCanvasWithPainter(painter *HeadlessPainter) *HeadlessCanvas {
	canvas := NewHeadlessCanvas()
	canvas.SetPainter(painter)

	return canvas
}

// NewHeadlessTransparentCanvasWithPainter allows creation of an in-memory canvas with a specific painter without a background color.
// The painter will be used to render in the Capture() call.
//
// Since: 2.2
func NewHeadlessTransparentCanvasWithPainter(painter *HeadlessPainter) *HeadlessCanvas {
	canvas := NewHeadlessCanvasWithPainter(painter)
	canvas.transparent = true

	return canvas
}

// Capture renders the content and returns the image
func (c *HeadlessCanvas) Capture() image.Image {
	var img image.Image
	img = c.Painter().Capture(c)
	return img
}

// Content returns the canvas object that is the content of the canvas
func (c *HeadlessCanvas) Content() fyne.CanvasObject {
	c.RLock()
	defer c.RUnlock()

	return c.content
}

// DismissMenu dismisses the menu if it is active
func (c *HeadlessCanvas) DismissMenu() bool {
	// c.RLock()
	// menu := c.menu
	// c.RUnlock()
	// if menu != nil && menu.(*MenuBar).IsActive() {
	// 	menu.(*MenuBar).Toggle()
	// 	return true
	// }
	return false
}

// InteractiveArea returns the position and size of the interactive area
func (c *HeadlessCanvas) InteractiveArea() (fyne.Position, fyne.Size) {
	return fyne.Position{}, c.Size()
}

// MinSize returns the minimum size for the canvas
func (c *HeadlessCanvas) MinSize() fyne.Size {
	c.RLock()
	defer c.RUnlock()
	return c.canvasSize(c.content.MinSize())
}

// OnTypedKey returns the function to be called on a typed key
func (c *HeadlessCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	c.RLock()
	defer c.RUnlock()

	return c.onTypedKey
}

// OnTypedRune returns the function to be called on typed rune
func (c *HeadlessCanvas) OnTypedRune() func(rune) {
	c.RLock()
	defer c.RUnlock()

	return c.onTypedRune
}

// Padded indicates whether the canvas is presently padded
func (c *HeadlessCanvas) Padded() bool {
	c.RLock()
	defer c.RUnlock()

	return c.padded
}

// PixelCoordinateForPosition returns the canvas scaled pixel coordinates for the indicated position
func (c *HeadlessCanvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	c.RLock()
	texScale := c.texScale
	c.RUnlock()
	multiple := c.Scale() * texScale
	scaleInt := func(x float32) int {
		return int(math.Round(float64(x * multiple)))
	}

	return scaleInt(pos.X), scaleInt(pos.Y)
}

// Resize sizes the canvas and its overlays to the given size
func (c *HeadlessCanvas) Resize(size fyne.Size) {
	// This might not be the ideal solution, but it effectively avoid the first frame to be blurry due to the
	// rounding of the size to the loower integer when scale == 1. It does not affect the other cases as far as we tested.
	// This can easily be seen with fyne/cmd/hello and a scale == 1 as the text will happear blurry without the following line.
	nearestSize := fyne.NewSize(float32(math.Ceil(float64(size.Width))), float32(math.Ceil(float64(size.Height))))

	c.Lock()
	c.size = nearestSize
	c.Unlock()

	for _, overlay := range c.Overlays().List() {
		if p, ok := overlay.(*widget.PopUp); ok {
			// TODO: remove this when #707 is being addressed.
			// “Notifies” the PopUp of the canvas size change.
			p.Refresh()
		} else {
			overlay.Resize(size)
		}
	}

	c.RLock()
	content := c.content
	contentSize := c.contentSize(nearestSize)
	contentPos := c.contentPos()
	menu := c.menu
	menuHeight := c.menuHeight()
	c.RUnlock()

	if content != nil {
		content.Resize(contentSize)
		content.Move(contentPos)
	}

	if menu != nil {
		menu.Refresh()
		menu.Resize(fyne.NewSize(nearestSize.Width, menuHeight))
	}
}

// Scale sets the scale cor the canvas
func (c *HeadlessCanvas) Scale() float32 {
	c.RLock()
	defer c.RUnlock()

	return c.scale
}

// SetContent sets the content for the canvas
func (c *HeadlessCanvas) SetContent(content fyne.CanvasObject) {
	content.Resize(content.MinSize()) // give it the space it wants then calculate the real min
	c.Lock()
	// the pass above makes some layouts wide enough to wrap, so we ask again what the true min is.
	newSize := c.size.Max(c.canvasSize(content.MinSize()))

	c.setContent(content)
	c.Unlock()

	c.Resize(newSize)
	c.SetDirty()
}

// SetOnTypedKey sets the handler for a typed key
func (c *HeadlessCanvas) SetOnTypedKey(handler func(*fyne.KeyEvent)) {
	c.Lock()
	defer c.Unlock()

	c.onTypedKey = handler
}

// SetOnTypedRune sets the handler for a typed rune
func (c *HeadlessCanvas) SetOnTypedRune(handler func(rune)) {
	c.Lock()
	defer c.Unlock()

	c.onTypedRune = handler
}

// SetPadded sets the padded flag for the canvas
func (c *HeadlessCanvas) SetPadded(padded bool) {
	c.Lock()
	content := c.content
	c.padded = padded
	pos := c.contentPos()
	c.Unlock()

	if content != nil {
		content.Move(pos)
	}
}

// SetScale sets the scale for the canvas
func (c *HeadlessCanvas) SetScale(scale float32) {
	c.Lock()
	defer c.Unlock()

	c.scale = scale
}

// Size returns the size of the canvas
func (c *HeadlessCanvas) Size() fyne.Size {
	c.RLock()
	defer c.RUnlock()

	return c.size
}

// canvasSize computes the needed canvas size for the given content size
func (c *HeadlessCanvas) canvasSize(contentSize fyne.Size) fyne.Size {
	canvasSize := contentSize.Add(fyne.NewSize(0, c.menuHeight()))
	if c.padded {
		return canvasSize.Add(fyne.NewSquareSize(theme.Padding() * 2))
	}
	return canvasSize
}

func (c *HeadlessCanvas) contentPos() fyne.Position {
	contentPos := fyne.NewPos(0, c.menuHeight())
	if c.padded {
		return contentPos.Add(fyne.NewSquareOffsetPos(theme.Padding()))
	}
	return contentPos
}

func (c *HeadlessCanvas) contentSize(canvasSize fyne.Size) fyne.Size {
	contentSize := fyne.NewSize(canvasSize.Width, canvasSize.Height-c.menuHeight())
	if c.Padded() {
		return contentSize.Subtract(fyne.NewSquareSize(theme.Padding() * 2))
	}
	return contentSize
}

func (c *HeadlessCanvas) menuHeight() float32 {
	if c.menu == nil {
		return 0 // no menu or native menu -> does not consume space on the canvas
	}

	return c.menu.MinSize().Height
}

func (c *HeadlessCanvas) objectTrees() []fyne.CanvasObject {
	trees := make([]fyne.CanvasObject, 0, len(c.Overlays().List())+1)
	if c.content != nil {
		trees = append(trees, c.content)
	}
	trees = append(trees, c.Overlays().List()...)
	return trees
}

func (c *HeadlessCanvas) paint(size fyne.Size) {
	clips := &internal.ClipStack{}
	if c.Content() == nil {
		return
	}
	c.Painter().Clear()

	paint := func(node *common.RenderCacheNode, pos fyne.Position) {
		obj := node.Obj()
		if _, ok := obj.(fyne.Scrollable); ok {
			inner := clips.Push(pos, obj.Size())
			c.Painter().StartClipping(inner.Rect())
		}
		if size.Width <= 0 || size.Height <= 0 { // iconifying on Windows can do bad things
			return
		}
		c.Painter().Paint(obj, pos, size)
	}
	afterPaint := func(node *common.RenderCacheNode, pos fyne.Position) {
		if _, ok := node.Obj().(fyne.Scrollable); ok {
			clips.Pop()
			if top := clips.Top(); top != nil {
				c.Painter().StartClipping(top.Rect())
			} else {
				c.Painter().StopClipping()
			}
		}

		if c.debug {
			c.DrawDebugOverlay(node.Obj(), pos, size)
		}
	}
	c.WalkTrees(paint, afterPaint)
}

func (c *HeadlessCanvas) setContent(content fyne.CanvasObject) {
	c.content = content
	c.SetContentTreeAndFocusMgr(content)
}

func layoutAndCollect(objects []fyne.CanvasObject, o fyne.CanvasObject, size fyne.Size) []fyne.CanvasObject {
	objects = append(objects, o)
	switch c := o.(type) {
	case fyne.Widget:
		r := c.CreateRenderer()
		r.Layout(size)
		for _, child := range r.Objects() {
			objects = layoutAndCollect(objects, child, child.Size())
		}
	case *fyne.Container:
		if c.Layout != nil {
			c.Layout.Layout(c.Objects, size)
		}
		for _, child := range c.Objects {
			objects = layoutAndCollect(objects, child, child.Size())
		}
	}
	return objects
}
