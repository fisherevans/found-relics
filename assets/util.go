package assets

import (
	"embed"
	"fmt"
	"io/fs"
)

func getAllFileBytes(efs *embed.FS) ([][]byte, error) {
	files, err := getAllFilePaths(efs)
	if err != nil {
		return nil, err
	}
	var datas [][]byte
	for _, file := range files {
		d, err := efs.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read embed file %s: %w", file, err)
		}
		datas = append(datas, d)
	}
	return datas, nil
}

func getAllFilePaths(efs *embed.FS) ([]string, error) {
	var files []string
	collect := func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	}
	err := fs.WalkDir(efs, ".", collect)
	if err != nil {
		return nil, err
	}
	return files, nil
}
