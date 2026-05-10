package constants

import "errors"

var (
	UserNotFound                = errors.New("User not found")
	UserExists                  = errors.New("User already exists")
	DocumentEditedByAnotherUser = errors.New("Document is currently edited by another user")
	DocumentNotFound            = errors.New("Document not found")
	DocumentIsFinalized         = errors.New("Document is finalized")
)
