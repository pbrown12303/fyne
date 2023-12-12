package nohead

import (
	"runtime"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/painter"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage/repository"
)

const defaultTitle = "Fyne Application"

// mainGoroutineID stores the main goroutine ID.
// This ID must be initialized in main.init because
// a main goroutine may not equal to 1 due to the
// influence of a garbage collector.
var mainGoroutineID uint64

var curWindow *headlessWindow

// A workaround on Apple M1/M2, just use 1 thread until fixed upstream.
const drawOnMainThread bool = runtime.GOOS == "darwin" && runtime.GOARCH == "arm64"

type headlessDriver struct {
	windowLock sync.RWMutex

	drawDone chan struct{}

	device       *headlessDevice
	painter      common.Painter
	windows      []fyne.Window
	windowsMutex sync.RWMutex
}

// Declare conformity with Driver
var _ fyne.Driver = (*headlessDriver)(nil)

// NewHeadlessDriver sets up and registers a new dummy driver for test purpose
func NewHeadlessDriver() *headlessDriver {
	drv := &headlessDriver{windowsMutex: sync.RWMutex{}}
	repository.Register("file", intRepo.NewFileRepository())

	// make a single dummy window for rendering tests
	drv.CreateWindow("")

	return drv
}

// NewHeadlessDriverWithPainter creates a new dummy driver that will pass the given
// painter to all canvases created
func NewHeadlessDriverWithPainter(painter common.Painter) *headlessDriver {
	return &headlessDriver{
		painter:      painter,
		windowsMutex: sync.RWMutex{},
	}
}

func (d *headlessDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	c := d.CanvasForObject(co)
	if c == nil {
		return fyne.NewPos(0, 0)
	}

	tc := c.(*HeadlessCanvas)
	return driver.AbsolutePositionForObject(co, tc.objectTrees())
}

func (d *headlessDriver) AllWindows() []fyne.Window {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	return d.windows
}

func (d *headlessDriver) CanvasForObject(fyne.CanvasObject) fyne.Canvas {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	// cheating: probably the last created window is meant
	return d.windows[len(d.windows)-1].Canvas()
}

func (d *headlessDriver) CreateWindow(title string) fyne.Window {
	var ret *headlessWindow
	if title == "" {
		title = defaultTitle
	}

	ret = &headlessWindow{title: title, decorate: true, driver: d}
	// This queue is destroyed when the window is closed.
	ret.InitEventQueue()
	go ret.RunEventQueue()

	ret.canvas = NewHeadlessCanvas()
	d.windowsMutex.Lock()
	d.windows = append(d.windows, ret)
	d.windowsMutex.Unlock()
	return ret
}

func (d *headlessDriver) Device() fyne.Device {
	if d.device == nil {
		d.device = &headlessDevice{}
	}
	return d.device
}

// RenderedTextSize looks up how bit a string would be if drawn on screen
func (d *headlessDriver) RenderedTextSize(text string, size float32, style fyne.TextStyle) (fyne.Size, float32) {
	return painter.RenderedTextSize(text, size, style)
}

func (d *headlessDriver) Run() {
	// no-op
}

func (d *headlessDriver) StartAnimation(a *fyne.Animation) {
	// currently no animations in test app, we just initialise it and leave
	a.Tick(1.0)
}

func (d *headlessDriver) StopAnimation(a *fyne.Animation) {
	// currently no animations in test app, do nothing
}

func (d *headlessDriver) Quit() {
	// no-op
}

func (d *headlessDriver) removeWindow(w *headlessWindow) {
	d.windowsMutex.Lock()
	i := 0
	for _, window := range d.windows {
		if window == w {
			break
		}
		i++
	}

	d.windows = append(d.windows[:i], d.windows[i+1:]...)
	d.windowsMutex.Unlock()
}
