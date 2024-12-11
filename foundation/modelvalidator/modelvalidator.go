// Package modelvalidator contains the support for validating models (structs).
package modelvalidator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/es"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	es_translations "github.com/go-playground/validator/v10/translations/es"
)

type validation struct {
	fn      func(fl validator.FieldLevel) bool
	message string
}

var extravalidations = map[string]validation{
	"rut":   {fn: validateRUT, message: "rut"},
	"phone": {fn: validateFormatPhone, message: "phone number"},
}

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale & translation information.
var translator ut.Translator

func init() {
	validate = validator.New()

	// Create a translator for spanish. Thus, the error messages are
	// more human-readable than technical.
	translator, _ = ut.New(es.New(), es.New()).GetTranslator("es")

	// Register the spanish error messages for use.
	es_translations.RegisterDefaultTranslations(validate, translator)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	registerValidations(translator)

}

func Check(val any, returnFirstError bool) error {

	if err := validate.Struct(val); err != nil {
		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
		if returnFirstError {
			return fmt.Errorf("%s", verrors[0].Translate(translator))
		}
		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Err:   verror.Translate(translator),
			}
			fields = append(fields, field)
		}
		return fields
	}

	return nil
}

func registerValidations(translator ut.Translator) {

	for k, v := range extravalidations {
		k := k
		v := v

		validate.RegisterValidation(k, v.fn)
		validate.RegisterTranslation(k, translator, func(ut ut.Translator) error {
			return ut.Add(k, fmt.Sprintf("Invalid %s", v.message), true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(k, fe.Field())
			return t
		})
	}

}
