package rooms

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewHandler(logger *zap.Logger) Handler {
	return Handler{logger: logger}
}

type Handler struct {
	logger *zap.Logger
}

type errorResponse struct {
	Message string `json:"message"`
}

type room struct {
	Id        int     `json:"id"`
	Name      *string `json:"name"`
	Player1Id string  `json:"player1Id"`
}

var rooms = map[int64]room{}

// TODO: check why this doesn't work:
// rooms := new(map[int64]room)

func (h Handler) Create(w http.ResponseWriter, req *http.Request) {
	data, err := io.ReadAll(req.Body)

	if err != nil {
		h.logger.Error("failed to read request body:", zap.Error(err))
		resp := errorResponse{Message: "unable to process the request"}
		bytes, err := json.Marshal(resp)
		if err != nil {
			h.logger.Error("failed to read request body:", zap.Error(err))
			w.Write([]byte(`{"error": "Internal Server Error"}`))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes)
		return
	}

	var newRoom room
	err = json.Unmarshal(data, &newRoom)

	if err != nil {
		h.logger.Error("failed to unmarshal:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		resp := errorResponse{Message: "unable to process the request"}
		bytes, err := json.Marshal(resp)
		if err != nil {
			h.logger.Error("failed to marshal error response:", zap.Error(err))
			w.Write([]byte(`{"error": "Internal Server Error"}`))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes)
		return
	}

	// TODO: check if the id already exists in the map
	newRoom.Id = rand.Intn(9999)
	h.logger.Info("appending new room", zap.String("roomName", *newRoom.Name))
	rooms[int64(newRoom.Id)] = newRoom

	bytes, err := json.Marshal(newRoom)
	if err != nil {
		h.logger.Error("failed to marshal error room response:", zap.Error(err))
		w.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}
	w.Write(bytes)
}

func (h Handler) Join(w http.ResponseWriter, req *http.Request) {
	param := mux.Vars(req)["roomCode"]
	roomId, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		return
	}

	foundRoom := rooms[roomId]
	// TODO: check if we actually found a room

	bytes, err := json.Marshal(foundRoom)

	if err != nil {
		h.logger.Error("failed to marshal error room response:", zap.Error(err))
		w.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}
	w.Write(bytes)
}
