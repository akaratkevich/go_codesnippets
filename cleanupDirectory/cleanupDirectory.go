package cleanupDirectory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CleanupDirectory(postfix string) error {
	dir := "D:\\vdcMigrations"
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), postfix) {
			fmt.Println("Removing file:", path)
			return os.Remove(path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
