package client

import (
	"fmt"
	"net/http"
	"strings"
)

// SEE: https://blog.golang.org/go1.13-errors

type JSONError struct {
	baseError error
	error     string
}

func (e *JSONError) Error() string {
	return e.baseError.Error()
}

func (e *JSONError) Unwrap() error {
	return e.baseError
}

type HTTPError struct {
	responseCode int
	baseError    error
	error        string
}

func (e *HTTPError) Error() string {
	return e.baseError.Error()
}

func (e *HTTPError) ResponseCode() int {
	return e.responseCode
}

func (e *HTTPError) Unwrap() error {
	return e.baseError
}

func HandleAPIError(httpResponse *http.Response, err error, errorMsg string) error {
	if err != nil {
		// var openAPIErr external.GenericOpenAPIError
		// if errors.As(err, &openAPIErr) {
		// 	log.Errorf("API response: %+v", string(openAPIErr.Body()))
		// }
		if httpResponse != nil {
			if httpResponse.StatusCode <= 300 {
				// openapi codegen suppresses json.UnmarshalTypeError and issues its own error
				// the API itself breaks the swagger spec,
				// this is the only way to catch and handle this behavior
				if strings.Contains(err.Error(), "json: cannot unmarshal") {
					newErr := JSONError{
						baseError: err,
						error:     fmt.Sprintf("%s: %v", errorMsg, err),
					}
					return &newErr
				}
			} else {
				newErr := HTTPError{
					responseCode: httpResponse.StatusCode,
					baseError:    err,
					error:        fmt.Sprintf("%s: %v", errorMsg, err),
				}
				return &newErr
			}
		}
		return fmt.Errorf("%s: %w", errorMsg, err)
	}
	if httpResponse.StatusCode > 300 {
		newErr := HTTPError{
			responseCode: httpResponse.StatusCode,
			baseError:    nil,
			error:        fmt.Sprintf("%s: %s", errorMsg, httpResponse.Status),
		}
		return &newErr
	}
	return nil
}
