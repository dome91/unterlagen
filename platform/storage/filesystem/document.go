package filesystem

import (
	"io"
	"path/filepath"
	fp "path/filepath"
	"unterlagen/features/archive"
	"unterlagen/platform/configuration"

	"github.com/spf13/afero"
)

var _ archive.DocumentStorage = &DocumentStorage{}
var _ archive.DocumentPreviewStorage = &DocumentPreviewStorage{}

type DocumentStorage struct {
	fs afero.Fs
}

// Retrieve implements unterlagen.DocumentStorage.
func (storage *DocumentStorage) Retrieve(filepath string, consumer archive.DocumentConsumer) error {
	file, err := storage.fs.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	return consumer(file)
}

func (storage *DocumentStorage) Store(filepath string, r io.Reader) error {
	// Create all necessary parent directories
	if err := storage.fs.MkdirAll(fp.Dir(filepath), 0755); err != nil {
		return err
	}

	// Create or truncate the file
	file, err := storage.fs.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the contents from the reader to the file
	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	return nil
}

// Delete implements unterlagen.DocumentStorage.
func (storage *DocumentStorage) Delete(filepath string) error {
	if err := storage.fs.Remove(filepath); err != nil {
		return err
	}

	// Check if directory is empty
	dir := fp.Dir(filepath)
	entries, err := afero.ReadDir(storage.fs, dir)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return storage.fs.Remove(dir)
	}

	return nil
}
func (storage *DocumentStorage) Size(filepath string) (int64, error) {
	fileInfo, err := storage.fs.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func NewDocumentStorage(configuration configuration.Configuration) *DocumentStorage {
	var fs afero.Fs
	if configuration.Production {
		fs = afero.NewOsFs()
	} else {
		fs = afero.NewMemMapFs()
	}

	base := filepath.Join(configuration.Data.Directory, "archive")
	fs = afero.NewBasePathFs(fs, base)
	return &DocumentStorage{
		fs: fs,
	}
}

type DocumentPreviewStorage struct {
	fs afero.Fs
}

// Delete implements document.PreviewStorage.
func (storage *DocumentPreviewStorage) Delete(preview string) error {
	if err := storage.fs.Remove(preview); err != nil {
		return err
	}

	// Check if directory is empty
	dir := fp.Dir(preview)
	entries, err := afero.ReadDir(storage.fs, dir)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return storage.fs.Remove(dir)
	}

	return nil
}

// Retrieve implements document.PreviewStorage.
func (storage *DocumentPreviewStorage) Retrieve(filepath string, consumer func(r io.Reader) error) error {
	file, err := storage.fs.Open(filepath)
	if err != nil {
		return err
	}
	//defer file.Close()
	return consumer(file)
}

// Store implements document.PreviewStorage.
func (storage *DocumentPreviewStorage) Store(filepath string, r io.Reader) error {
	// Create all necessary parent directories
	if err := storage.fs.MkdirAll(fp.Dir(filepath), 0755); err != nil {
		return err
	}

	// Create or truncate the file
	file, err := storage.fs.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the contents from the reader to the file
	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	return nil
}

func NewDocumentPreviewStorage(configuration configuration.Configuration) *DocumentPreviewStorage {
	var fs afero.Fs
	if configuration.Production {
		fs = afero.NewOsFs()
	} else {
		fs = afero.NewMemMapFs()
	}

	base := filepath.Join(configuration.Data.Directory, "archive")
	fs = afero.NewBasePathFs(fs, base)

	return &DocumentPreviewStorage{
		fs: fs,
	}
}
