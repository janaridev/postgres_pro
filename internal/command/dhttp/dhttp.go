package dhttp

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	stdstring "strings"

	"github.com/go-playground/validator/v10"
	"github.com/janaridev/postgres_pro/internal/command/service"
	"github.com/janaridev/postgres_pro/internal/command/types"
	"github.com/janaridev/postgres_pro/pkg/api/response"
	"github.com/janaridev/postgres_pro/pkg/logger/sl"
)

type Handler struct {
	log     *slog.Logger
	v       *validator.Validate
	service types.CommandService
}

func New(log *slog.Logger, service types.CommandService) *Handler {
	v := validator.New()

	return &Handler{
		log:     log,
		v:       v,
		service: service,
	}
}

const (
	applicationJSON = "application/json"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	cmds, err := h.service.List(r.Context())
	if err != nil {
		response.SendJSONResponse(w, applicationJSON, http.StatusInternalServerError, response.NewErrorResponse("something went wrong"))
		return
	}

	var res []types.GetCommandResponse
	for _, cmd := range cmds {
		res = append(res, types.GetCommandResponse{
			ID:        cmd.ID,
			Name:      cmd.Name,
			Raw:       cmd.Raw,
			Logs:      cmd.Logs.String,
			ErrorMsg:  cmd.ErrorMsg.String,
			Status:    cmd.Status.String,
			CreatedAt: cmd.CreatedAt,
			UpdatedAt: cmd.UpdatedAt.Time,
		})
	}

	response.SendJSONResponse(w, applicationJSON, http.StatusOK, res)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "dhttp.Get"

	logger := h.log.With(slog.String("op", op))

	parts := stdstring.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		logger.Warn("failed to get id", sl.Err(err))

		response.SendJSONResponse(w, applicationJSON, http.StatusBadRequest, response.NewErrorResponse("provide correct id"))
		return
	}

	cmd, err := h.service.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCommandNotFound) {
			response.SendJSONResponse(w, applicationJSON, http.StatusNotFound, response.NewErrorResponse("command not found"))
			return
		}

		response.SendJSONResponse(w, applicationJSON, http.StatusInternalServerError, response.NewErrorResponse("failed to get command"))
		return
	}

	res := types.GetCommandResponse{
		ID:        cmd.ID,
		Name:      cmd.Name,
		Raw:       cmd.Raw,
		Logs:      cmd.Logs.String,
		ErrorMsg:  cmd.ErrorMsg.String,
		Status:    cmd.Status.String,
		CreatedAt: cmd.CreatedAt,
		UpdatedAt: cmd.UpdatedAt.Time,
	}

	response.SendJSONResponse(w, applicationJSON, http.StatusOK, res)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "dhttp.Create"

	logger := h.log.With(slog.String("op", op))

	var req types.CreateCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("failed to decode payload", sl.Err(err))

		response.SendJSONResponse(w, applicationJSON, http.StatusBadRequest, response.NewErrorResponse("failed to decode payload"))
		return
	}

	if err := h.v.Struct(req); err != nil {
		response.SendJSONResponse(w, applicationJSON, http.StatusBadRequest, err.(validator.ValidationErrors))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.(validator.ValidationErrors))
		return
	}

	cmd, err := h.service.Create(r.Context(), req.Name, req.Raw)
	if err != nil {
		if errors.Is(err, service.ErrCommandAlreadyExists) {
			response.SendJSONResponse(w, applicationJSON, http.StatusBadRequest, response.NewErrorResponse("command already exists"))
			return
		}

		response.SendJSONResponse(w, applicationJSON, http.StatusInternalServerError, response.NewErrorResponse("failed to create command"))
		return
	}

	res := types.CreateCommandResponse{
		ID: cmd.ID,
	}

	response.SendJSONResponse(w, applicationJSON, http.StatusCreated, res)
}

func (h *Handler) Stop(w http.ResponseWriter, r *http.Request) {
	const op = "dhttp.Remove"

	logger := h.log.With(slog.String("op", op))

	parts := stdstring.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		logger.Warn("failed to get id", sl.Err(err))

		response.SendJSONResponse(w, applicationJSON, http.StatusBadRequest, response.NewErrorResponse("provide correct id"))
		return
	}

	if err := h.service.Stop(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrCommandNotFound) {
			response.SendJSONResponse(w, applicationJSON, http.StatusNotFound, response.NewErrorResponse("command not found"))
			return
		}

		response.SendJSONResponse(w, applicationJSON, http.StatusInternalServerError, response.NewErrorResponse("something went wrong"))
		return
	}

	response.SendJSONResponse(w, applicationJSON, http.StatusOK, nil)
}
