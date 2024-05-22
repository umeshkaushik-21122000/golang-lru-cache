package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/umeshkaushik-21122000/golang-lru-cache/app"
	"github.com/umeshkaushik-21122000/golang-lru-cache/types"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler interface {
	Get(http.ResponseWriter, *http.Request)
	Set(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

type HandlerImpl struct {
	lruCache app.LRUCacheInterface
}

func NewHandler(lruCache app.LRUCacheInterface) *HandlerImpl {
	return &HandlerImpl{
		lruCache: lruCache,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections by default
	},
}

func (h *HandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // Ensure ticker is stopped to avoid resource leak
	for {
		select {
		case <-ticker.C:
			res := h.lruCache.GetAll()
			fmt.Println(res)
			bytes, err := json.Marshal(res)
			if err != nil {
				log.Println("Error marshalling response:", err)
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
				log.Println("Error writing JSON:", err)
				continue
			}
		}
	}
}

func (h *HandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	integerNumber, err := strconv.Atoi(key)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	serviceRes := h.lruCache.Get(integerNumber)
	res := types.KeyValue{
		Key:   integerNumber,
		Value: serviceRes,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}

func (h *HandlerImpl) Set(w http.ResponseWriter, r *http.Request) {
	var kv types.KeyValue
	if err := json.NewDecoder(r.Body).Decode(&kv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("--------------->>>>>>>>>", kv)
	h.lruCache.Put(kv.Key, kv.Value, kv.Expiration)
	w.WriteHeader(http.StatusOK)
}

func (h *HandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	integerNumber, err := strconv.Atoi(key)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.lruCache.Remove(integerNumber)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}
	w.WriteHeader(http.StatusOK)
}
