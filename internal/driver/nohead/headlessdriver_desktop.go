package nohead

import (
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"fyne.io/fyne/v2"
)

var (
	systrayIcon fyne.Resource
	setup       sync.Once
)

func goroutineID() (id uint64) {
	var buf [30]byte
	runtime.Stack(buf[:], false)
	for i := 10; buf[i] != ' '; i++ {
		id = id*10 + uint64(buf[i]&15)
	}
	return id
}

func (d *headlessDriver) catchTerm() {
	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	d.Quit()
}

func addMissingQuitForMenu(menu *fyne.Menu, d *headlessDriver) {
	var lastItem *fyne.MenuItem
	if len(menu.Items) > 0 {
		lastItem = menu.Items[len(menu.Items)-1]
		if lastItem.Label == "Quit" {
			lastItem.IsQuit = true
		}
	}
	if lastItem == nil || !lastItem.IsQuit { // make sure the menu always has a quit option
		quitItem := fyne.NewMenuItem("Quit", nil)
		quitItem.IsQuit = true
		menu.Items = append(menu.Items, fyne.NewMenuItemSeparator(), quitItem)
	}
	for _, item := range menu.Items {
		if item.IsQuit && item.Action == nil {
			item.Action = d.Quit
		}
	}
}
