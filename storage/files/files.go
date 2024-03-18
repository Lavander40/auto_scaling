package files

import (
	"auto_scaling/lib/e"
	"auto_scaling/storage"
	"encoding/gob"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func New(path string) Storage {
	return Storage{
		basePath: path,
	}
}

func (s Storage) Save(call *storage.Call) error {
	filePath := filepath.Join(s.basePath, call.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return e.WrapErr("can't save", err)
	}

	fileName, err := fileName(call)
	if err != nil {
		return e.WrapErr("can't save", err)
	}

	filePath = filepath.Join(filePath, fileName)

	log.Print(filePath)

	file, err := os.Create(filePath)
	if err != nil {
		return e.WrapErr("can't save", err)
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(call); err != nil {
		return e.WrapErr("can't save", err)
	}

	return nil
}

func (s Storage) PickLastCalls(userName string) ([]*storage.Call, error) {
	path := filepath.Join(s.basePath, userName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, storage.ErrNoDir
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, e.WrapErr("can't find last calls", err)
	}

	if len(files) == 0 {
		return nil, e.WrapErr("can't find last calls", storage.ErrEmpty)
	}

	var file []fs.DirEntry
	// SELECT * FROM calls SORT BY time_created LIMIT 10
	if len(files) > 10 {
		file = files[len(files)-10:]
	} else {
		file = files
	}
	
	var res []*storage.Call

	for _, f := range file {
		r, err := s.decodecall(filepath.Join(path, f.Name()))
		if err != nil {
			return nil, e.WrapErr("can't find last calls", err)
		}
		res = append(res, r)
	}
	
	return res, nil
}

func (s Storage) Remove(p *storage.Call) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.WrapErr("can't remove", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return e.WrapErr("can't remove", err)
	}

	return nil
}

func (s Storage) decodecall(filePath string) (*storage.Call, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.WrapErr("can't decode call", err)
	}
	defer f.Close()

	var p storage.Call

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.WrapErr("can't decode call", err)
	}

	return &p, nil
}

func fileName(p *storage.Call) (string, error) {
	return p.Hash()
}
