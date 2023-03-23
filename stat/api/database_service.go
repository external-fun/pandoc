package api

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

type DatabaseService struct {
	db *sql.DB
}

type Status string

const (
	Uploaded   = "uploaded"
	Converting = "converting"
	Converted  = "converted"
)

func getConnectionUrl() string {
	return os.Getenv("DB_CONNECT_URL")
}

func NewDatabaseService() (*DatabaseService, error) {
	db, err := sql.Open("postgres", getConnectionUrl())
	if err != nil {
		return nil, err
	}

	return &DatabaseService{
		db: db,
	}, nil
}

func (service *DatabaseService) UploadFile(uuid string, convertedUrl string) error {
	_, err := service.db.Exec(`UPDATE public."file" SET converted_url=$2, status=$3 WHERE uuid=$1`,
		uuid,
		convertedUrl,
		Converted)
	return err
}
