package filemanager_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	fmngr "github.com/Moon1706/ginkgo2allure/pkg/convert/file_manager"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/assert"
)

func TestFileManagerSuccess(t *testing.T) {
	id := uuid.New()
	resultsFolderName := fmt.Sprintf("test-%s", id.String())
	resultsPath := filepath.Join(os.TempDir(), resultsFolderName)
	err := os.MkdirAll(resultsPath, os.ModePerm)
	assert.Empty(t, err, "tmp test dir should create successful")

	fm := fmngr.NewFileManager(resultsPath)
	err = fm.SaveJSONResult(allure.Result{
		UUID: id,
	})
	assert.Empty(t, err, "report saved successful")
}
