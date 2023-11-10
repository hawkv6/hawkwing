package config

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-playground/validator"
)

var validate *validator.Validate

func Validate() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("mapstructure")
	})

	validate.RegisterStructValidation(ServiceConfigValidation, ServiceConfig{})
	validate.RegisterStructValidation(ApplicationValidation, Application{})
	validate.RegisterStructValidation(IntentValidation, Intent{})

	err := validate.Struct(Params)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Info(err)
		}

		for _, err := range err.(validator.ValidationErrors) {
			e := validationError{
				Namespace:   err.Namespace(),
				Field:       err.Field(),
				ActualValue: fmt.Sprintf("%v", err.Value()),
				Message:     err.Tag(),
			}

			indent, err := json.MarshalIndent(e, "", "  ")
			if err != nil {
				log.Fatalf("failed to validate config: %v", err)
			}
			fmt.Println(string(indent))
			log.Fatalf("failed to validate config with errors above")
		}
	}

}
