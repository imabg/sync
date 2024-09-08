package validate

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/imabg/sync/pkg/errors"
	"github.com/imabg/sync/pkg/response"
)

var validate *validator.Validate

func SetupValidation() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func GetPayload(r *http.Request, payload interface{}) error {
	err := response.PaseRequest(r, payload)
	if err != nil {
		return err
	}

	err = validate.Struct(payload)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsg := fmt.Sprintf("Error:Field validation for '%s' failed on the '%s' tag", err.Field(), err.Tag())
			return errors.NewCustomError(http.StatusBadRequest, errMsg, "", "")
		}
	}
	return err
}
