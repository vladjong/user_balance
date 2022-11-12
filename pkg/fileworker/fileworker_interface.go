package fileworker

import "github.com/vladjong/user_balance/internal/entities"

type FileWorker interface {
	Record(records []entities.Report, header []string, date string) (string, error)
}
