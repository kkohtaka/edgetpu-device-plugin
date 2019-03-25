package fileutil

import (
	"io/ioutil"
	"os"
	"path"

	errors "golang.org/x/xerrors"
)

func FindFiles(dir, target string) ([]string, error) {
	_, err := os.Stat(dir)
	if err != nil {
		return nil, errors.Errorf("stat dir: %w", err)
	}

	var res []string
	if files, err := ioutil.ReadDir(dir); err != nil {
		return nil, errors.Errorf("list files in directory: %w", err)
	} else {
		for _, file := range files {
			p := path.Join(dir, file.Name())
			if file.IsDir() {
				if r, err := FindFiles(p, target); err != nil {
					return nil, errors.Errorf("find file in directory: %w", err)
				} else {
					res = append(res, r...)
				}
			} else if file.Name() == target {
				res = append(res, p)
			}
		}
	}
	return res, nil
}
