package nohead

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/storage"
)

type headlessStorage struct {
	*internal.Docs
}

func (s *headlessStorage) RootURI() fyne.URI {
	return storage.NewFileURI(os.TempDir())
}

func (s *headlessStorage) docRootURI() (fyne.URI, error) {
	return storage.Child(s.RootURI(), "Documents")
}
