package entity

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

var (
	DocumentNotFoundError = Error{
		Code:    1,
		Message: "Document not found",
	}

	DocumentFinalizedError = Error{
		Code:    2,
		Message: "Document is already finalized",
	}

	DocumentEditedByOtherError = Error{
		Code:    3,
		Message: "Document is edited by other",
	}

	InternalServerError = Error{
		Code:    -1,
		Message: "Internal error, please try again",
	}

	BadRequestError = Error{
		Code:    -2,
		Message: "Bad request",
	}
)
