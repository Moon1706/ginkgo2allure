package filemanager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/pkg/errors"
)

const (
	fileSystemPermissionCode = 0644
)

type FileManager interface {
	SaveJSONResult(result allure.Result) error
}

type fileManager struct {
	resultsFolderPath string
}

func NewFileManager(resultsFolderPath string) FileManager {
	return &fileManager{
		resultsFolderPath: resultsFolderPath,
	}
}

func (m *fileManager) createFile(name string, content []byte) error {
	return os.WriteFile(filepath.Join(m.resultsFolderPath, name), content, fileSystemPermissionCode)
}

func (m *fileManager) SaveJSONResult(result allure.Result) error {
	bResult, err := json.Marshal(result)
	if err != nil {
		return errors.Wrap(err, "Failed marshal Result")
	}

	err = m.createFile(fmt.Sprintf("%s-result.json", result.UUID), bResult)
	if err != nil {
		return errors.Wrap(err, "Cannot save Result")
	}
	return nil
}
