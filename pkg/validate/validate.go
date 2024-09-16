package validate

import (
	brazilian_portuguese "github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	br_translations "github.com/go-playground/validator/v10/translations/pt_BR"
)

var Validate *validator.Validate
var Translate ut.Translator

func init() {
	ptbr := brazilian_portuguese.New()
	uni := ut.New(ptbr, ptbr)
	Translate, _ = uni.GetTranslator("pt_BR")

	Validate = validator.New()
	br_translations.RegisterDefaultTranslations(Validate, Translate)
}

func Struct(str interface{}) []string {
	var errors []string

	err := Validate.Struct(str)

	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			errors = append(errors, "Erro interno ao validar campos")
			return errors
		}

		for _, e := range errs {
			errors = append(errors, e.Translate(Translate))
		}
	}

	return errors
}
