package controller

import (
	"context"
	"log"

	"github.com/alvinmatias69/editor-history/internal/constants"
	"github.com/alvinmatias69/editor-history/internal/entity"
	"github.com/google/uuid"
)

type DocumentController struct {
	documentRepository           documentRepository
	documentOperationsRepository documentOperationsRepository
	sessionRepository            sessionRepository
}

func NewDocumentController(
	documentRepository documentRepository,
	documentOperationsRepository documentOperationsRepository,
	sessessionRepository sessionRepository,
) *DocumentController {
	return &DocumentController{
		documentRepository:           documentRepository,
		documentOperationsRepository: documentOperationsRepository,
		sessionRepository:            sessessionRepository,
	}
}

func (c *DocumentController) Create(ctx context.Context) (entity.Document, error) {
	document := entity.Document{
		ID:         uuid.NewString(),
		ReadOnlyID: uuid.NewString(),
		IsFinished: false,
	}

	err := c.documentRepository.Create(ctx, document)

	return document, err
}

func (c *DocumentController) Finish(ctx context.Context, sessionRequest entity.SessionRequest) error {
	err := c.checkSession(ctx, sessionRequest)
	if err != nil {
		return err
	}

	return c.documentRepository.Finish(ctx, sessionRequest.DocumentID)
}

func (c *DocumentController) Get(ctx context.Context, request entity.GetDocumentRequest) (entity.GetDocumentResponse, error) {
	if request.AccessType == constants.Write {
		err := c.checkSession(ctx, entity.SessionRequest{
			DocumentID: request.DocumentID,
			EditorID:   request.EditorID,
		})

		if err != nil {
			return entity.GetDocumentResponse{}, err
		}
	}

	var (
		document entity.Document
		err      error
	)
	if request.AccessType == constants.Read {
		document, err = c.documentRepository.GetForRead(ctx, request.DocumentID)
	} else {
		document, err = c.documentRepository.GetForWrite(ctx, request.DocumentID)
	}

	if err != nil {
		return entity.GetDocumentResponse{}, err
	}

	operations, err := c.documentOperationsRepository.Get(ctx, document.ID)
	if err != nil {
		return entity.GetDocumentResponse{}, err
	}

	response := entity.GetDocumentResponse{
		ID:         document.ID,
		ReaderID:   document.ReadOnlyID,
		IsFinished: document.IsFinished,
		Operations: make([]entity.DocumentOperationResponse, 0, len(operations)),
	}

	if request.AccessType == constants.Read {
		response.ID = ""
	}

	for _, operation := range operations {
		response.Operations = append(response.Operations, entity.DocumentOperationResponse{
			ID:                 operation.ID,
			EditorID:           operation.EditorID,
			PositionStart:      operation.PositionStart,
			PositionEnd:        operation.PositionEnd,
			Value:              operation.Value,
			Operation:          operation.Operation,
			OperationTimestamp: operation.OperationTimestamp,
		})
	}

	if request.AccessType == constants.Write {
		err = c.sessionRepository.Set(ctx, entity.SessionRequest{
			DocumentID: request.DocumentID,
			EditorID:   request.EditorID,
		})
		if err != nil {
			log.Printf("Error setting session: %v", err)
		}
	}

	return response, nil
}

func (c *DocumentController) AddOperations(ctx context.Context, request entity.AddOperationsRequest) error {
	log.Printf("request: %v\n", request)
	err := c.checkSession(ctx, entity.SessionRequest{
		DocumentID: request.DocumentID,
		EditorID:   request.EditorID,
	})
	if err != nil {
		return err
	}

	document, err := c.documentRepository.GetForWrite(ctx, request.DocumentID)
	if err != nil {
		return err
	}

	if document.IsFinished {
		return constants.DocumentIsFinalized
	}

	operations := make([]entity.DocumentOperation, 0, len(request.OperationRequests))
	for _, operation := range request.OperationRequests {
		operations = append(operations, entity.DocumentOperation{
			ID:                 operation.ID,
			DocumentID:         request.DocumentID,
			EditorID:           request.EditorID,
			PositionStart:      operation.PositionStart,
			PositionEnd:        operation.PositionEnd,
			Value:              operation.Value,
			Operation:          operation.Operation,
			OperationTimestamp: operation.Timestamp,
		})
	}

	num, err := c.documentOperationsRepository.Add(ctx, operations)
	if err == nil {
		log.Printf("Success adding %v operations", num)
	}

	sessionErr := c.sessionRepository.Set(ctx, entity.SessionRequest{
		DocumentID: request.DocumentID,
		EditorID:   request.EditorID,
	})
	if sessionErr != nil {
		log.Printf("Error setting session: %v", sessionErr)
	}

	return err
}

func (c *DocumentController) OverrideEditor(ctx context.Context, request entity.SessionRequest) error {
	return c.sessionRepository.Set(ctx, request)
}

func (c *DocumentController) checkSession(ctx context.Context, sessionRequest entity.SessionRequest) error {
	currentEditorID, err := c.sessionRepository.Get(ctx, sessionRequest.DocumentID)
	if err != nil {
		return err
	}

	if currentEditorID != "" && currentEditorID != sessionRequest.EditorID {
		return constants.DocumentEditedByAnotherUser
	}

	return nil
}
