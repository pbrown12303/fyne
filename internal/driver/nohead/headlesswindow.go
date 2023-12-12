package nohead

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/common"
)

// This content is a meld of glfw.window.go and glfw.window_desktop.go

// Declare conformity to Window interface
var _ fyne.Window = (*headlessWindow)(nil)

type headlessWindow struct {
	common.Window

	viewport     *headlessWindow
	viewLock     sync.RWMutex
	decorate     bool
	closing      bool
	visible      bool
	shouldExpand bool

	title              string
	fullScreen         bool
	fixedSize          bool
	focused            bool
	onClosed           func()
	onCloseIntercepted func()

	canvas    *HeadlessCanvas
	clipboard fyne.Clipboard
	driver    *headlessDriver
	menu      *fyne.MainMenu
}

// NewHeadlessWindow creates and registers a new window for test purposes
func NewHeadlessWindow(content fyne.CanvasObject) fyne.Window {
	window := fyne.CurrentApp().NewWindow("")
	window.SetContent(content)
	return window
}

func (w *headlessWindow) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *headlessWindow) CenterOnScreen() {
	// no-op
}

func (w *headlessWindow) Clipboard() fyne.Clipboard {
	if w.clipboard == nil {
		w.clipboard = &headlessClipboard{window: w.viewport}
	}
	return w.clipboard
}

func (w *headlessWindow) Close() {
	if w.onClosed != nil {
		w.onClosed()
	}
	w.focused = false
	w.driver.removeWindow(w)
}

func (w *headlessWindow) Content() fyne.CanvasObject {
	return w.Canvas().Content()
}

func (w *headlessWindow) FixedSize() bool {
	return w.fixedSize
}

func (w *headlessWindow) FullScreen() bool {
	return w.fullScreen
}

func (w *headlessWindow) Hide() {
	w.focused = false
}

func (w *headlessWindow) Icon() fyne.Resource {
	return fyne.CurrentApp().Icon()
}

func (w *headlessWindow) isClosing() bool {
	w.viewLock.RLock()
	closing := w.closing || w.viewport == nil
	w.viewLock.RUnlock()
	return closing
}

func (w *headlessWindow) MainMenu() *fyne.MainMenu {
	return w.menu
}

func (w *headlessWindow) Padded() bool {
	return w.canvas.Padded()
}

func (w *headlessWindow) RequestFocus() {
	for _, win := range w.driver.AllWindows() {
		win.(*headlessWindow).focused = false
	}

	w.focused = true
}

func (w *headlessWindow) Resize(size fyne.Size) {
	w.canvas.Resize(size)
}

func (w *headlessWindow) SetContent(obj fyne.CanvasObject) {
	w.Canvas().SetContent(obj)
}

func (w *headlessWindow) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
}

func (w *headlessWindow) SetIcon(_ fyne.Resource) {
	// no-op
}

func (w *headlessWindow) SetFullScreen(fullScreen bool) {
	w.fullScreen = fullScreen
}

func (w *headlessWindow) SetMainMenu(menu *fyne.MainMenu) {
	w.menu = menu
}

func (w *headlessWindow) SetMaster() {
	// no-op
}

func (w *headlessWindow) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *headlessWindow) SetCloseIntercept(callback func()) {
	w.onCloseIntercepted = callback
}

func (w *headlessWindow) SetOnDropped(dropped func(fyne.Position, []fyne.URI)) {

}

func (w *headlessWindow) SetPadded(padded bool) {
	w.canvas.SetPadded(padded)
}

func (w *headlessWindow) SetTitle(title string) {
	w.title = title
}

func (w *headlessWindow) Show() {
	w.RequestFocus()
}

func (w *headlessWindow) ShowAndRun() {
	w.Show()
}

func (w *headlessWindow) Title() string {
	return w.title
}

func (w *headlessWindow) view() *headlessWindow {
	w.viewLock.RLock()
	defer w.viewLock.RUnlock()

	if w.closing {
		return nil
	}
	return w.viewport
}
