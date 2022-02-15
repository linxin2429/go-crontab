package common

import (
	"github.com/pkg/errors"
	"net/http"
)

func HttpInternalErrorHandle(w http.ResponseWriter, err error) {
	http.Error(w, "internal error", http.StatusInternalServerError)
	Logger.Errorln(errors.Errorf("%w", err))
}
