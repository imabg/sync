package validate

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/imabg/sync/pkg/response"
)

var validate *validator.Validate

func SetupValidation() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}



func ValidateRequestPayload (r *http.Request, payload interface{}) error {

	err := response.GetPayload(r, payload)
	if err != nil {
		return err
	}

	err = validate.Struct(payload)

	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		// from here you can create your own error messages in whatever language you wish
	}
	return err
}
