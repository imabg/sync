package response

import (
	"encoding/json"
	"net/http"

	"github.com/imabg/sync/pkg/errors"
	"github.com/imabg/sync/pkg/uuid"
)


func PaseRequest (r *http.Request, userType interface{}) error {
	return json.NewDecoder(r.Body).Decode(&userType)
}

func setBasicHeaders (w http.ResponseWriter, code int) {
	w.Header().Add("x-sync-id", uuid.GenerateUUID())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
}

func Send(w http.ResponseWriter, code int, data interface{}) {
	setBasicHeaders(w, code)
	databytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(databytes)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	
}

func SendWithError(w http.ResponseWriter, code int, e errors.CustomError) {
	setBasicHeaders(w, code)
	databytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(databytes)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	
}