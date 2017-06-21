package errors

import (
	"encoding/json"
	"fmt"
	"gokit/log"
)

type DataError struct {
	Code       int64  `json:"code"`
	Message    string `json:"message"`
	_cache     string
	_jsonCache string
}

func (e *DataError) Error() string {
	if e._cache != "" {
		return e._cache
	}

	e._cache = fmt.Sprintf(` *****code: %d, #message: %s*****`, e.Code, e.Message)
	return e._cache
}

func (e *DataError) Json() string {
	if e._jsonCache != "" {
		return e._jsonCache
	}

	bin, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		log.Error("Marshal response: ", *e, " to json failed, error: ", err)
		e._jsonCache = fmt.Sprintf(`{"code": %d, "message": "%s"}`, e.Code, e.Message)
	} else {
		e._jsonCache = string(bin)
	}

	return e._jsonCache
}

func (e *DataError) With(m string) *DataError {
	return Error(e.Code, m)
}

func Error(c int64, m string) *DataError {
	return &DataError{Code: c, Message: m}
}
