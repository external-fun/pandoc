package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type ConverterService struct {
	s3 *S3Service
	db *DatabaseService
}

type RequestData struct {
	Uuid      string `json:"uuid"`
	OriginUrl string `json:"originUrl"`
	From      string `json:"from"`
	To        string `json:"to"`
}

var (
	RabbitMqUrl = os.Getenv("RABBIT_MQ_URL")
	QueueName   = os.Getenv("RABBIT_MQ_NAME")
)

func TryDial(url string, tries int) *amqp.Connection {
	for i := 0; i < tries; i++ {
		conn, err := amqp.Dial(RabbitMqUrl)
		if err == nil {
			return conn
		}
		time.Sleep(5 * time.Second)
	}
	panic(fmt.Sprintf("couldn't dial %s", url))
}

func NewConverterService(s3 *S3Service, db *DatabaseService) *ConverterService {
	return &ConverterService{
		s3: s3,
		db: db,
	}
}

// TODO: Maybe we should panic
func (service *ConverterService) Serve() {
	conn := TryDial(RabbitMqUrl, 10)
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(QueueName, false, false, false, false, nil)
	if err != nil {
		log.Println(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Println(err)
	}

	for msg := range msgs {
		log.Printf("Got %s\n", msg.Body)

		var data RequestData
		err = json.Unmarshal(msg.Body, &data)
		if err != nil {
			log.Println(err)
			continue
		}

		id, err := service.convertFile(&data)
		if err != nil {
			log.Println(err)
			err = msg.Nack(false, true)
			if err != nil {
				log.Println(err)
			}
			continue
		}

		err = msg.Ack(false)
		if err != nil {
			log.Println(err)
			continue
		}

		err = service.db.UploadFile(data.Uuid, getConvertedUrl(id))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (service *ConverterService) convertFile(r *RequestData) (string, error) {
	file, err := downloadFile(r.OriginUrl)
	if err != nil {
		return "", err
	}
	defer file.Close()

	id := uuid.New().String()
	cmd := exec.Command("pandoc", "-f", r.From, "-t", r.To, "-o", id, file.Name())
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}

	result, err := os.Open(id)
	if err != nil {
		return "", err
	}
	defer result.Close()

	err = service.s3.Upload(id, result)
	if err != nil {
		return "", err
	}

	return id, nil
}

var (
	storageBaseUrl = os.Getenv("STORAGE_BASE_URL")
)

func getConvertedUrl(uuid string) string {
	return "http://" + storageBaseUrl + "/" + uuid
}

func downloadFile(url string) (*os.File, error) {
	file, err := os.CreateTemp("", "converter")
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		file.Close()
		return nil, err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}
