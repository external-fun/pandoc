package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
)

type ConverterService struct {
	s3 *S3Service
	db *DatabaseService
	mq *MqService
}

func NewConverterService(s3 *S3Service, db *DatabaseService, mq *MqService) *ConverterService {
	return &ConverterService{
		s3: s3,
		db: db,
		mq: mq,
	}
}

type RequestError struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func requestError(error *RequestError) []byte {
	data, err := json.Marshal(error)
	if err != nil {
		log.Println("RequestError ", err)
		return []byte("")
	}
	return data
}

var (
	storageBaseUrl = os.Getenv("STORAGE_BASE_URL")
)

func getOriginUrl(uuid string) string {
	return fmt.Sprintf("http://%s/%s", storageBaseUrl, uuid)
}

type RequestData struct {
	Uuid      string `json:"uuid"`
	OriginUrl string `json:"originUrl"`
	From      string `json:"from"`
	To        string `json:"to"`
}

func (service *ConverterService) parseRequestData(r *http.Request) (*RequestData, error) {
	parts, err := r.MultipartReader()
	if err != nil {
		return &RequestData{}, err
	}

	data := RequestData{}
	for {
		part, err := parts.NextPart()
		if err == io.EOF {
			return &data, nil
		} else if err != nil {
			return nil, err
		}

		if part.FormName() == "file" {
			// TODO: should be able to revert in case of error
			id := uuid.New().String()
			err := service.s3.Upload(id, part)
			if err != nil {
				return nil, err
			}
			err = service.db.UploadFile(id, getOriginUrl(id))
			if err != nil {
				return nil, err
			}

			data.OriginUrl = getOriginUrl(id)
			data.Uuid = id
		}

		if part.FormName() == "from" {
			val, err := io.ReadAll(part)
			if err != nil {
				return nil, err
			}

			data.From = string(val)
		}

		if part.FormName() == "to" {
			val, err := io.ReadAll(part)
			if err != nil {
				return nil, err
			}

			data.To = string(val)
		}
	}
}

type UploadResponse struct {
	Uuid string `json:"uuid"`
}

func (service *ConverterService) uploadHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling request")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(204)
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(400)
			w.Write(requestError(&RequestError{
				Name:        "Wrong Method",
				Description: "Only POST is allowed",
			}))
			return
		}

		req, err := service.parseRequestData(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(requestError(&RequestError{
				Name:        "Internal error",
				Description: "Couldn't upload a file",
			}))
			return
		}
		err = service.mq.Upload(req)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(requestError(&RequestError{
				Name:        "Internal error",
				Description: "Couldn't upload a file",
			}))
			return
		}

		resp, _ := json.Marshal(UploadResponse{
			Uuid: req.Uuid,
		})
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		log.Println("Successfully uploaded file")
	}
}

type StatusResponse struct {
	Status string `json:"status"`
}

func (service *ConverterService) statusHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(204)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(requestError(&RequestError{
				Name:        "Wrong Method",
				Description: "Only GET is allowed",
			}))
			return
		}

		uuid := r.URL.Query().Get("uuid")
		if uuid == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(requestError(&RequestError{
				Name:        "Invalid UUID",
				Description: "No UUID is given",
			}))
			return
		}

		status, err := service.db.GetStatus(uuid)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(requestError(&RequestError{
				Name:        "Status not found",
				Description: "Status for UUID wasn't found",
			}))
			return
		}

		resp, _ := json.Marshal(StatusResponse{
			Status: string(status),
		})
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func (service *ConverterService) Serve(addr string) {
	http.HandleFunc("/api/v1/upload", service.uploadHandler())
	http.HandleFunc("/api/v1/status", service.statusHandler())
	http.ListenAndServe(addr, nil)
}
