package nohead

import (
	"runtime"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver/common"
)

type funcData struct {
	f    func()
	done chan struct{} // Zero allocation signalling channel
}

type drawData struct {
	f    func()
	win  *headlessWindow
	done chan struct{} // Zero allocation signalling channel
}

// channel for queuing functions on the main thread
var funcQueue = make(chan funcData)
var drawFuncQueue = make(chan drawData)
var running uint32 // atomic bool, 0 or 1
var initOnce = &sync.Once{}

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
	mainGoroutineID = goroutineID()
}

// force a function f to run on the main thread
func runOnMain(f func()) {
	// If we are on main just execute - otherwise add it to the main queue and wait.
	// The "running" variable is normally false when we are on the main thread.
	if onMain := atomic.LoadUint32(&running) == 0; onMain {
		f()
		return
	}

	done := common.DonePool.Get().(chan struct{})
	defer common.DonePool.Put(done)

	funcQueue <- funcData{f: f, done: done}

	<-done
}

// // force a function f to run on the draw thread
// func runOnDraw(w *headlessWindow, f func()) {
// 	if drawOnMainThread {
// 		runOnMain(func() { w.RunWithContext(f) })
// 		return
// 	}
// 	done := common.DonePool.Get().(chan struct{})
// 	defer common.DonePool.Put(done)

// 	drawFuncQueue <- drawData{f: f, win: w, done: done}
// 	<-done
// }

func (d *headlessDriver) drawSingleFrame() {
	refreshingCanvases := make([]fyne.Canvas, 0)
	for _, win := range d.windowList() {
		w := win.(*headlessWindow)
		w.viewLock.RLock()
		canvas := w.canvas
		closing := w.closing
		visible := w.visible
		w.viewLock.RUnlock()

		// CheckDirtyAndClear must be checked after visibility,
		// because when a window becomes visible, it could be
		// showing old content without a dirty flag set to true.
		// Do the clear if and only if the window is visible.
		if closing || !visible || !canvas.CheckDirtyAndClear() {
			continue
		}

		d.repaintWindow(w)
		refreshingCanvases = append(refreshingCanvases, canvas)
	}
	cache.CleanCanvases(refreshingCanvases)
}

// func (d *headlessDriver) runGL() {
// 	if !atomic.CompareAndSwapUint32(&running, 0, 1) {
// 		return // Run was called twice.
// 	}
// 	close(d.waitForStart) // Signal that execution can continue.

// 	d.initGLFW()
// 	if d.trayStart != nil {
// 		d.trayStart()
// 	}
// 	fyne.CurrentApp().Lifecycle().(*app.Lifecycle).TriggerStarted()
// 	eventTick := time.NewTicker(time.Second / 60)
// 	for {
// 		select {
// 		case <-d.done:
// 			eventTick.Stop()
// 			d.drawDone <- struct{}{} // wait for draw thread to stop
// 			d.Terminate()
// 			fyne.CurrentApp().Lifecycle().(*app.Lifecycle).TriggerStopped()
// 			return
// 		case f := <-funcQueue:
// 			f.f()
// 			f.done <- struct{}{}
// 		case <-eventTick.C:
// 			d.tryPollEvents()
// 			windowsToRemove := 0
// 			for _, win := range d.windowList() {
// 				w := win.(*window)
// 				if w.viewport == nil {
// 					continue
// 				}

// 				if w.viewport.ShouldClose() {
// 					windowsToRemove++
// 					continue
// 				}

// 				w.viewLock.RLock()
// 				expand := w.shouldExpand
// 				fullScreen := w.fullScreen
// 				w.viewLock.RUnlock()

// 				if expand && !fullScreen {
// 					w.fitContent()
// 					w.viewLock.Lock()
// 					shouldExpand := w.shouldExpand
// 					w.shouldExpand = false
// 					view := w.viewport
// 					w.viewLock.Unlock()
// 					if shouldExpand {
// 						view.SetSize(w.shouldWidth, w.shouldHeight)
// 					}
// 				}

// 				if drawOnMainThread {
// 					d.drawSingleFrame()
// 				}
// 			}
// 			if windowsToRemove > 0 {
// 				oldWindows := d.windowList()
// 				newWindows := make([]fyne.Window, 0, len(oldWindows)-windowsToRemove)

// 				for _, win := range oldWindows {
// 					w := win.(*window)
// 					if w.viewport == nil {
// 						continue
// 					}

// 					if w.viewport.ShouldClose() {
// 						w.viewLock.Lock()
// 						w.visible = false
// 						v := w.viewport
// 						w.viewLock.Unlock()

// 						// remove window from window list
// 						v.Destroy()
// 						w.destroy(d)
// 						continue
// 					}

// 					newWindows = append(newWindows, win)
// 				}

// 				d.windowLock.Lock()
// 				d.windows = newWindows
// 				d.windowLock.Unlock()

// 				if len(newWindows) == 0 {
// 					d.Quit()
// 				}
// 			}
// 		}
// 	}
// }

func (d *headlessDriver) repaintWindow(w *headlessWindow) {
	canvas := w.canvas
	if canvas.EnsureMinSize() {
		w.viewLock.Lock()
		w.shouldExpand = true
		w.viewLock.Unlock()
	}
	canvas.FreeDirtyTextures()

	canvas.paint(canvas.Size())

}

// func (d *headlessDriver) startDrawThread() {
// 	settingsChange := make(chan fyne.Settings)
// 	fyne.CurrentApp().Settings().AddChangeListener(settingsChange)
// 	var drawCh <-chan time.Time
// 	if drawOnMainThread {
// 		drawCh = make(chan time.Time) // don't tick when on M1
// 	} else {
// 		drawCh = time.NewTicker(time.Second / 60).C
// 	}

// 	go func() {
// 		runtime.LockOSThread()

// 		for {
// 			select {
// 			case <-d.drawDone:
// 				return
// 			case f := <-drawFuncQueue:
// 				f.win.RunWithContext(f.f)
// 				f.done <- struct{}{}
// 			case set := <-settingsChange:
// 				painter.ClearFontCache()
// 				cache.ResetThemeCaches()
// 				app.ApplySettingsWithCallback(set, fyne.CurrentApp(), func(w fyne.Window) {
// 					c, ok := w.Canvas().(*HeadlessCanvas)
// 					if !ok {
// 						return
// 					}
// 					c.applyThemeOutOfTreeObjects()
// 					go c.reloadScale()
// 				})
// 			case <-drawCh:
// 				d.drawSingleFrame()
// 			}
// 		}
// 	}()
// }

func (d *headlessDriver) windowList() []fyne.Window {
	d.windowLock.RLock()
	defer d.windowLock.RUnlock()
	return d.windows
}

// refreshWindow requests that the specified window be redrawn
func refreshWindow(w *headlessWindow) {
	w.canvas.SetDirty()
}
