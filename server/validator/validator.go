package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	errors "github.com/vasupal1996/goerror"
)

// Validator container validator library and transaltor
type Validator struct {
	V *validator.Validate
	T *ut.Translator
}

// NewValidation create new Validator struct instance
func NewValidation() *Validator {
	v := &Validator{
		V: validator.New(),
	}
	trans := initializeTranslation(v.V)
	v.T = trans
	registerFunc(v.V)
	return v
}

// Initialize initializes and returns the UniversalTranslator instance for the application
func initializeTranslation(validate *validator.Validate) *ut.Translator {

	// initialize translator
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")
	// initialize translations
	en_translations.RegisterDefaultTranslations(validate, trans)
	return &trans
}

func registerFunc(validate *validator.Validate) {
	// register function to get tag name from json tags.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	validate.RegisterValidation("required_with_field", requiredWithField)
}

// Validate validates the struct
// Note: do not pass slice of struct
func (v *Validator) Validate(form interface{}) []error {
	var errResp []error
	if err := v.V.Struct(form); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			err := errors.New(e.Translate(*v.T), &errors.BadRequest)
			key := strings.SplitAfterN(e.Namespace(), ".", 2)
			err = errors.SetContext(err, key[1], e.Translate(*v.T))
			errResp = append(errResp, err)
		}
	}
	return errResp
}

// requiredWithField validates B field is also required if A field is not empty or nil
var requiredWithField validator.Func = func(fl validator.FieldLevel) bool {
	// Checking if current field is empty or not
	var otherField string = fl.Param()
	var otherFieldVal reflect.Value

	if fl.Parent().Kind() == reflect.Ptr {
		otherFieldVal = fl.Parent().Elem().FieldByName(otherField)
	} else {
		otherFieldVal = fl.Parent().FieldByName(otherField)
	}

	switch isNilOrZeroValue(fl.Field()) {
	case true:
		if isNilOrZeroValue(otherFieldVal) {
			return false
		}
	case false:
		if !(isNilOrZeroValue(otherFieldVal)) {
			return false
		}
	}
	return true
}

func isNilOrZeroValue(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsZero()
	}
}
