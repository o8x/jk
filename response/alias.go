package response

import "net/http"

var (
	NoContent = &Response{
		StatusCode: http.StatusNoContent,
	}

	BadRequest = &Response{
		StatusCode: http.StatusBadRequest,
	}
)
