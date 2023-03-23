package api

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type StatService struct {
	db *DatabaseService
}

func NewStatService(db *DatabaseService) *StatService {
	return &StatService{
		db: db,
	}
}

var upgrader = websocket.Upgrader{}

func (service *StatService) statHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer conn.Close()

		for {
			time.Sleep(5 * time.Second)
			err := conn.WriteMessage(websocket.BinaryMessage, []byte("hello world"))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (service *StatService) Serve(addr string) {
	http.HandleFunc("api/v1/stat", service.statHandler())
	http.ListenAndServe(addr, nil)
}
