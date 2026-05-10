package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/alvinmatias69/editor-history/internal/config"
	"github.com/alvinmatias69/editor-history/internal/constants"
	"github.com/alvinmatias69/editor-history/internal/entity"
)

const idPath = "id"
const editorIDQueryKey = "editor-id"

type DocumentHandler struct {
	config     config.Base
	controller documentController
}

func NewDocumentHandler(config config.Base, controller documentController) *DocumentHandler {
	return &DocumentHandler{
		config:     config,
		controller: controller,
	}
}

func (h *DocumentHandler) Create(res http.ResponseWriter, req *http.Request) {
	ctx, cancelCtx := context.WithTimeout(req.Context(), time.Duration(h.config.DefaultTimeoutSeconds)*time.Second)
	defer cancelCtx()

	res.Header().Set("Content-Type", "application/json")

	document, err := h.controller.Create(ctx)
	if err != nil {
		log.Printf("error: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(entity.InternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(document)
}

func (h *DocumentHandler) Finish(res http.ResponseWriter, req *http.Request) {
	ctx, cancelCtx := context.WithTimeout(req.Context(), time.Duration(h.config.DefaultTimeoutSeconds)*time.Second)
	defer cancelCtx()

	res.Header().Set("Content-Type", "application/json")

	var requestBody entity.SessionRequest
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("error: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(entity.BadRequestError)
		return
	}
	requestBody.DocumentID = req.PathValue(idPath)

	err = h.controller.Finish(ctx, requestBody)
	if errors.Is(err, constants.DocumentNotFound) {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(entity.DocumentNotFoundError)
		return
	}

	if errors.Is(err, constants.DocumentEditedByAnotherUser) {
		res.WriteHeader(http.StatusConflict)
		json.NewEncoder(res).Encode(entity.DocumentEditedByOtherError)
		return
	}

	if err != nil {
		log.Printf("error: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(entity.InternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (h *DocumentHandler) GetForWrite(res http.ResponseWriter, req *http.Request) {
	ctx, cancelCtx := context.WithTimeout(req.Context(), time.Duration(h.config.DefaultTimeoutSeconds)*time.Second)
	defer cancelCtx()

	res.Header().Set("Content-Type", "application/json")

	editorID := req.URL.Query().Get(editorIDQueryKey)
	if editorID == "" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(entity.BadRequestError)
		return
	}

	payload := entity.GetDocumentRequest{
		DocumentID: req.PathValue(idPath),
		EditorID:   editorID,
		AccessType: constants.Write,
	}

	result, err := h.controller.Get(ctx, payload)
	if errors.Is(err, constants.DocumentNotFound) {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(entity.DocumentNotFoundError)
		return
	}

	if errors.Is(err, constants.DocumentEditedByAnotherUser) {
		res.WriteHeader(http.StatusConflict)
		json.NewEncoder(res).Encode(entity.DocumentEditedByOtherError)
		return
	}

	if err != nil {
		log.Printf("error: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(entity.InternalServerError)
		return
	}

	json.NewEncoder(res).Encode(result)
}

func (h *DocumentHandler) GetForRead(res http.ResponseWriter, req *http.Request) {
	ctx, cancelCtx := context.WithTimeout(req.Context(), time.Duration(h.config.DefaultTimeoutSeconds)*time.Second)
	defer cancelCtx()

	res.Header().Set("Content-Type", "application/json")

	payload := entity.GetDocumentRequest{
		DocumentID: req.PathValue(idPath),
		EditorID:   "",
		AccessType: constants.Read,
	}

	result, err := h.controller.Get(ctx, payload)
	if errors.Is(err, constants.DocumentNotFound) {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(entity.DocumentNotFoundError)
		return
	}

	if err != nil {
		log.Printf("error: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(entity.InternalServerError)
		return
	}

	json.NewEncoder(res).Encode(result)
}

func (h *DocumentHandler) AddOperations(res http.ResponseWriter, req *http.Request) {
	ctx, cancelCtx := context.WithTimeout(req.Context(), time.Duration(h.config.DefaultTimeoutSeconds)*time.Second)
	defer cancelCtx()

	res.Header().Set("Content-Type", "application/json")

	// TODO: Add validation
	var requestBody entity.AddOperationsRequest
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("error: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(entity.BadRequestError)
		return
	}
	requestBody.DocumentID = req.PathValue(idPath)

	err = h.controller.AddOperations(ctx, requestBody)
	if errors.Is(err, constants.DocumentNotFound) {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(entity.DocumentNotFoundError)
		return
	}

	if errors.Is(err, constants.DocumentIsFinalized) {
		res.WriteHeader(http.StatusConflict)
		json.NewEncoder(res).Encode(entity.DocumentFinalizedError)
		return
	}

	if errors.Is(err, constants.DocumentEditedByAnotherUser) {
		res.WriteHeader(http.StatusConflict)
		json.NewEncoder(res).Encode(entity.DocumentEditedByOtherError)
		return
	}

	if err != nil {
		log.Printf("error: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(entity.InternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (h *DocumentHandler) OverrideEditor(res http.ResponseWriter, req *http.Request) {
	ctx, cancelCtx := context.WithTimeout(req.Context(), time.Duration(h.config.DefaultTimeoutSeconds)*time.Second)
	defer cancelCtx()

	res.Header().Set("Content-Type", "application/json")

	// TODO: Add validation
	var requestBody entity.SessionRequest
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("error: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(entity.BadRequestError)
		return
	}

	requestBody.DocumentID = req.PathValue(idPath)

	err = h.controller.OverrideEditor(ctx, requestBody)
	if errors.Is(err, constants.DocumentNotFound) {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(entity.DocumentNotFoundError)
		return
	}

	if err != nil {
		log.Printf("error: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(entity.InternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}
