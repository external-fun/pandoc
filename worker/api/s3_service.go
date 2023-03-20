package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

type S3Service struct {
	session *session.Session
	config  *S3Config
}

type S3Config struct {
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
	BucketName      string
}

func NewS3Service(config *S3Config) (*S3Service, error) {
	customResolver := func(service string, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		return endpoints.ResolvedEndpoint{
			URL:           config.Endpoint,
			SigningRegion: "ru-msk",
		}, nil
	}

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String("ru-msk"),
			Credentials: credentials.NewStaticCredentials(
				config.AccessKeyId,
				config.SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
			EndpointResolver: endpoints.ResolverFunc(customResolver),
		})
	if err != nil {
		return nil, err
	}
	return &S3Service{
		session: sess,
		config:  config,
	}, nil
}

func (service *S3Service) Upload(name string, data io.Reader) error {
	uploader := s3manager.NewUploader(service.session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(service.config.BucketName),
		ACL:    aws.String("public-read"),
		Key:    aws.String(name),
		Body:   data,
	})
	if err != nil {
		return err
	}
	return nil
}
