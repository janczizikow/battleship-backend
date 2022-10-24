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
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type room struct {
	Id        int     `json:"id"`
	Name      *string `json:"name"`
	Player1Id string  `json:"player1Id"`
}

var rooms = map[int64]room{}

func (r *room) GenerateId() {
	id := rand.Intn(9999)
	if _, exists := rooms[int64(id)]; exists {
		r.GenerateId()
	} else {
		r.Id = id
	}
}

func (h Handler) Create(w http.ResponseWriter, req *http.Request) {
	data, err := io.ReadAll(req.Body)

	if err != nil {
		h.logger.Error("failed to read request body:", zap.Error(err))
		resp := errorResponse{Status: http.StatusBadRequest, Message: "unable to process the request"}
		WriteErrorResponse(w, req, resp)
		return
	}

	var newRoom room
	err = json.Unmarshal(data, &newRoom)

	if err != nil {
		h.logger.Error("failed to unmarshal:", zap.Error(err))
		resp := errorResponse{Status: http.StatusBadRequest, Message: "unable to process the request"}
		WriteErrorResponse(w, req, resp)
		return
	}

	newRoom.GenerateId()
	h.logger.Info("appending new room", zap.String("roomName", *newRoom.Name))
	rooms[int64(newRoom.Id)] = newRoom

	err = WriteJSON(w, http.StatusCreated, newRoom)
	if err != nil {
		h.logger.Error("failed to marshal error room response:", zap.Error(err))
		WriteErrorResponse(w, req, errorResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}
}

func (h Handler) Join(w http.ResponseWriter, req *http.Request) {
	param := mux.Vars(req)["roomCode"]
	roomId, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		WriteErrorResponse(w, req, errorResponse{Status: http.StatusNotFound, Message: "Room Not Found"})
		return
	}

	foundRoom, exists := rooms[roomId]

	if !exists {
		h.logger.Info("couldn't find a room", zap.Int64("roomId", roomId))
		WriteErrorResponse(w, req, errorResponse{Status: http.StatusNotFound, Message: "Room Not Found"})
		return
	}

	err = WriteJSON(w, http.StatusCreated, foundRoom)

	if err != nil {
		h.logger.Error("failed to marshal error room response:", zap.Error(err))
		WriteErrorResponse(w, req, errorResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}
}

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, resp errorResponse) {
	data := map[string]interface{}{"status": resp.Status, "error": resp.Message}
	err := WriteJSON(w, resp.Status, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	_, err = w.Write(bytes)

	return err
}
