package response

import (
	"encoding/json"
	"net/http"
)


func GetPayload (r *http.Request, userType interface{}) error {
	return json.NewDecoder(r.Body).Decode(&userType)
}