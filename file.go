package venom

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	jsonKey = "json"
)

// extensionMap is the collection of file extensions to the IOFileLoaders that
// can load files with the associated extensions
var extensionMap = map[string]IOFileLoader{
	jsonKey: JSONLoader,
}

// RegisterExtension registers an IOFileLoader for the provided file extension
func RegisterExtension(ext string, loader IOFileLoader) {
	extensionMap[ext] = loader
}

// ErrNoFileLoader is the error returned when a file is attempted to be loaded
// without a matching extension IOFileLoader
type ErrNoFileLoader struct {
	ext string
}

// Error implements the error interface and returns a custom error message for
// the current ErrNoFileLoader instance
func (e ErrNoFileLoader) Error() string {
	return fmt.Sprintf("venom: no loader for extension %q", e.ext)
}

// IOFileLoader is the function signature for a function which can load an
// io.Reader into a map[string]interface{}
type IOFileLoader func(io.Reader) (map[string]interface{}, error)

// JSONLoader is an IOFileLoader which loads JSON config data
func JSONLoader(r io.Reader) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// LoadFile loads the file from the provided path into Venoms configs. If the
// file can't be opened, if no loader for the files extension exists, or if
// loading the file fails, an error is returned
func (v *Venom) LoadFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}

	ext := strings.TrimLeft(filepath.Ext(name), ".")
	loader, ok := extensionMap[ext]
	if !ok {
		return ErrNoFileLoader{ext}
	}

	data, err := loader(file)
	if err != nil {
		return err
	}

	v.Merge(FileLevel, data)
	return nil
}

func findFiles(dir string, recurse bool) (files sort.StringSlice) {
	files = sort.StringSlice{}

	path, _ := filepath.Abs(dir)

	walk := func(file string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for extension := range extensionMap {
			if !i.IsDir() && strings.HasSuffix(file, extension) {
				files = append(files, strings.Replace(file, "\\", "/", -1))
			} else if i.IsDir() && !recurse && path != file {
				// don't recurse into subdirectories
				return filepath.SkipDir
			}
		}
		return nil
	}
	_ = filepath.Walk(path, walk)
	return
}

// LoadDirectory loads any config files found in the provided directory,
// optionally recursing into any sub-directories
func (v *Venom) LoadDirectory(dir string, recurse bool) error {
	configFiles := findFiles(dir, recurse)
	for _, file := range configFiles {
		if err := v.LoadFile(file); err != nil {
			return err
		}
	}
	return nil
}
