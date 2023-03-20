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

func (service *DatabaseService) UploadFile(uuid string, originUrl string) error {
	_, err := service.db.Exec(`INSERT INTO public."file"(uuid, origin_url, status) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
		uuid,
		originUrl,
		Uploaded)
	return err
}

func (service *DatabaseService) GetStatus(uuid string) (Status, error) {
	row := service.db.QueryRow(`SELECT status FROM public."file" WHERE uuid=$1`, uuid)
	var status Status
	err := row.Scan(&status)
	if err != nil {
		return "", err
	}
	return status, err
}
