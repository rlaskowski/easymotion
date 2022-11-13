package storage

type FileStorage struct {
	path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path}
}

func (f *FileStorage) Store(name string, data []byte) error {
	return nil
}
