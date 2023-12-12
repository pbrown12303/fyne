package nohead

import (
	"image"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/theme"
)

// HeadlessPainter is a simple software painter that can paint a canvas in memory.
type HeadlessPainter struct {
	canvas *HeadlessCanvas
}

// Declare compliance with the painter.Painter interface
var _ common.Painter = (*HeadlessPainter)(nil)

// NewHeadlessPainter creates a new HeadlessPainter.
func NewHeadlessPainter() *HeadlessPainter {
	return &HeadlessPainter{}
}

// Init tell a new painter to initialise, usually called after a context is available
func (hp *HeadlessPainter) Init() {

}

// Capture requests that the specified canvas be drawn to an in-memory image
func (hp *HeadlessPainter) Capture(c fyne.Canvas) image.Image {
	bounds := image.Rect(0, 0, scale.ToScreenCoordinate(c, c.Size().Width), scale.ToScreenCoordinate(c, c.Size().Height))
	switch typedCanvas := c.(type) {
	case *HeadlessCanvas:
		if typedCanvas.Painter() != nil {
			height := bounds.Dy()
			width := bounds.Dx()
			size := fyne.NewSize(float32(width), float32(height))
			hp.Paint(hp.canvas.content, fyne.NewPos(0, 0), size)
		}
		return typedCanvas.renderedImage
	}
	return nil
}

// Clear tells our painter to prepare a fresh paint
func (hp *HeadlessPainter) Clear() {

}

// Free is used to indicate that a certain canvas object is no longer needed
func (hp *HeadlessPainter) Free(fyne.CanvasObject) {

}

// Paint renders the canvas object into the HeadlessCanvas's renderedImage. The old renderedImage is discarded.
func (hp *HeadlessPainter) Paint(canvasObject fyne.CanvasObject, position fyne.Position, size fyne.Size) {
	can := hp.canvas
	bounds := image.Rect(0, 0, scale.ToScreenCoordinate(can, can.Size().Width), scale.ToScreenCoordinate(can, can.Size().Height))
	base := image.NewNRGBA(bounds)

	if can.transparent {
		draw.Draw(base, bounds, image.NewUniform(theme.BackgroundColor()), image.Point{}, draw.Src)
	}

	paint := func(obj fyne.CanvasObject, pos, clipPos fyne.Position, clipSize fyne.Size) bool {
		w := fyne.Min(clipPos.X+clipSize.Width, can.Size().Width)
		h := fyne.Min(clipPos.Y+clipSize.Height, can.Size().Height)
		clip := image.Rect(
			scale.ToScreenCoordinate(can, clipPos.X),
			scale.ToScreenCoordinate(can, clipPos.Y),
			scale.ToScreenCoordinate(can, w),
			scale.ToScreenCoordinate(can, h),
		)
		switch o := obj.(type) {
		case *canvas.Image:
			drawImage(can, o, pos, base, clip)
		case *canvas.Text:
			drawText(can, o, pos, base, clip)
		case gradient:
			drawGradient(can, o, pos, base, clip)
		case *canvas.Circle:
			drawCircle(can, o, pos, base, clip)
		case *canvas.Line:
			drawLine(can, o, pos, base, clip)
		case *canvas.Raster:
			drawRaster(can, o, pos, base, clip)
		case *canvas.Rectangle:
			drawRectangle(can, o, pos, base, clip)
		}

		return false
	}

	driver.WalkVisibleObjectTree(can.Content(), paint, nil)
	for _, o := range can.Overlays().List() {
		driver.WalkVisibleObjectTree(o, paint, nil)
	}

	can.renderedImage = base
}

// SetFrameBufferScale tells us when we have more than 1 framebuffer pixel for each output pixel
func (hp *HeadlessPainter) SetFrameBufferScale(float32) {

}

// SetOutputSize is used to change the resolution of our output viewport
func (hp *HeadlessPainter) SetOutputSize(int, int) {

}

// StartClipping tells us that the following paint actions should be clipped to the specified area.
func (hp *HeadlessPainter) StartClipping(fyne.Position, fyne.Size) {

}

// StopClipping stops clipping paint actions.
func (hp *HeadlessPainter) StopClipping() {

}
