package common

import (
	"net/http"
)

func HttpInternalErrorHandle(w http.ResponseWriter, err error) {
	http.Error(w, "internal error", http.StatusInternalServerError)
	Logger.Errorln(err)
}
