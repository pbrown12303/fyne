package nohead

import (
	"fyne.io/fyne/v2"
)

type headlessDevice struct {
}

// Declare conformity with Device
var _ fyne.Device = (*headlessDevice)(nil)

func (d *headlessDevice) Orientation() fyne.DeviceOrientation {
	return fyne.OrientationVertical
}

func (d *headlessDevice) HasKeyboard() bool {
	return false
}

func (d *headlessDevice) SystemScale() float32 {
	return d.SystemScaleForWindow(nil)
}

func (d *headlessDevice) SystemScaleForWindow(fyne.Window) float32 {
	return 1
}

func (d *headlessDevice) IsBrowser() bool {
	return false
}

func (d *headlessDevice) IsMobile() bool {
	return false
}
