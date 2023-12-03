package storage

import (
	"encoding/json"
	"errors"
	"os"
)

func (container *LinksMapping) SaveJSONToFile(path string) {
	if path == "" {
		return
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()
	jsoned, err := json.MarshalIndent(container.byShortMap, "", "	")
	if err != nil {
		panic(err)
	}
	_, err = file.Write(jsoned)
	if err != nil {
		panic(err)
	}
}

func (container *LinksMapping) LoadFromJSONFile(path string) error {
	if path == "" {
		return nil
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &container.byShortMap)
	if err != nil {
		return err
	}
	for k, v := range container.byShortMap {
		container.byFullMap[v] = k
	}

	return nil
}
