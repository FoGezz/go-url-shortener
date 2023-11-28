package storage

import (
	"encoding/json"
	"errors"
	"os"
)

func (lm *LinksMapping) SaveJsonToFile(path string) {
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
	jsoned, err := json.MarshalIndent(lm.byShortMap, "", "	")
	if err != nil {
		panic(err)
	}
	_, err = file.Write(jsoned)
	if err != nil {
		panic(err)
	}
}

func (lm *LinksMapping) LoadFromJsonFile(path string) {
	if path == "" {
		return
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return
	}
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &lm.byShortMap)
	if err != nil {
		panic(err)
	}
	for k, v := range lm.byShortMap {
		lm.byFullMap[v] = k
	}
}
