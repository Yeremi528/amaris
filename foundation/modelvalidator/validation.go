package modelvalidator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

func validateRUT(fl validator.FieldLevel) bool {

	rut := fmt.Sprintf("%v", fl.Field())

	rut = strings.ReplaceAll(rut, ".", "")
	rut = strings.ReplaceAll(rut, "-", "")
	cuerpo := rut[:len(rut)-1]
	dv := strings.ToUpper(string(rut[len(rut)-1]))

	if len(cuerpo) < 7 {
		return false
	}

	suma := 0
	multiplo := 2
	for i := 1; i <= len(cuerpo); i++ {
		index := multiplo * int(rut[len(cuerpo)-i]-'0')
		suma += index
		if multiplo < 7 {
			multiplo++
		} else {
			multiplo = 2
		}
	}

	dvEsperado := 11 - (suma % 11)
	if dv == "K" {
		dv = "10"
	} else if dv == "0" {
		dv = "11"
	}

	if fmt.Sprintf("%d", dvEsperado) != dv {
		return false
	}

	return true
}

func validateFormatPhone(fl validator.FieldLevel) bool {

	phoneNumber := fmt.Sprintf("%v", fl.Field())

	num, err := phonenumbers.Parse(phoneNumber, "CL")
	if err != nil {
		return false
	}

	if !phonenumbers.IsValidNumberForRegion(num, "CL") {
		return false
	}

	return true

}
