package OpenExcelFile

import (
	"github.com/tealeg/xlsx"
)

func OpenExcelFile(fileName string) (*xlsx.File, error) {
	// Open Excel file
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}
