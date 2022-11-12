package fileworker

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/vladjong/user_balance/internal/entities"
)

type workerCsv struct{}

func New() *workerCsv {
	return &workerCsv{}
}

func (f *workerCsv) Record(records []entities.Report, header []string, date string) (string, error) {
	filename := fmt.Sprintf("data/report_%s.csv", date)
	outputFile, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()
	if err := writer.Write(header); err != nil {
		return "", err
	}
	for _, record := range records {
		var csvRow []string
		csvRow = append(csvRow, fmt.Sprint(record.Id), record.Name, record.AllSum.String())
		if err := writer.Write(csvRow); err != nil {
			return "", err
		}
	}
	return filename, nil
}
