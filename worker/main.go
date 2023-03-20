package main

import (
	"github.com/external-fun/pandoc/worker/api"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	s3, err := api.NewS3Service(&api.S3Config{
		Endpoint:        os.Getenv("S3_ENDPOINT"),
		AccessKeyId:     os.Getenv("S3_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
		BucketName:      os.Getenv("S3_BUCKET_NAME"),
	})
	if err != nil {
		panic("Couldn't connect to S3 with config ")
	}
	db, err := api.NewDatabaseService()
	if err != nil {
		panic("Couldn't connect to Db")
	}
	service := api.NewConverterService(s3, db)
	service.Serve()
}
