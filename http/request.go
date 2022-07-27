package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func decodeRequest(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		// EOF means no body was sent and can be treated as non-error (further validation should likely fail)
		if errors.Is(err, io.EOF) {
			return nil
		}

		return err
	}

	return nil
}
