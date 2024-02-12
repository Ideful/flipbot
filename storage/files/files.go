package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"

	"github.com/Ideful/flipbot/storage"
)

const defaultPerm = 0774

var ErrNoSavedPages = errors.New("no saved pages")

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {
	fPath := path.Join(s.basePath, page.UserName)

	if err := os.Mkdir(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("Hash issue:%v", err)
	}

	fPath = filepath.Join(fPath, fName)
	file, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("filepath create issue:%v", err)
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return fmt.Errorf("gob error:%v", err)
	}
	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("readDir issue:%v", err)
	}
	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("file remove issue:%v", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("issue while trying to remove %s:%v", path, err)
	}
	return nil
}

func (s Storage) Exists(p *storage.Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, fmt.Errorf("file existence:%v", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fName)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil

	case err != nil:
		return false, fmt.Errorf("issue while checking existence of %s:%v", path, err)
	}
	return true, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("filepath open error:%v", err)
	}
	defer f.Close()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("gob issue:%v", err)
	}
	return &p, nil
}
