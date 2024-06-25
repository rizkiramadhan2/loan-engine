package validate

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// Validate struct
type Validate struct {
	Validator *validator.Validate
	trans     ut.Translator
}

// Options struct
type Options struct {
}

// ErrorValidation struct format for returning message
type ErrorValidation struct {
	Message string `json:"message"`
}

type Assertion struct {
	AssertType     int         `json:"assert_type"`
	AssertValid    bool        `json:"assertTypeValid"`
	Visibility     bool        `json:"visibility"` // TODO: deprecate
	IsEnabled      bool        `json:"is_enabled"`
	Locator        string      `json:"locator"`
	Validator      string      `json:"validator"`
	ValidatorValid bool        `json:"validatorValid"`
	Value          interface{} `json:"value"`
	ValidatorText  string      `json:"validator_text"`
	Valid          bool        `json:"valid"`
	ResponseType   int         `json:"response_type"`
	ExcludeKey     []string    `json:"exclude_key"`
}

// New method construct a Validate struct with pre registered custom rule
func New(cfg *Options) *Validate {
	v := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(v, trans)

	_ = v.RegisterValidation("datetime", validateIsDatetime)

	// use json tag for message
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validate{
		Validator: v,
		trans:     trans,
	}
}

// ValidateStruct is a method under Validate struct that accept any struct (interface{}) to perform validation
func (v *Validate) validateStruct(structType interface{}) []ErrorValidation {
	err := v.Validator.Struct(structType)
	if err == nil {
		return nil
	}

	var errs []ErrorValidation
	if _, ok := err.(*validator.InvalidValidationError); ok {
		errs = append(errs, ErrorValidation{
			Message: err.Error(),
		})
		return errs
	}

	return v.translateErrMessage(err)
}

// ValidateStruct is a method under Validate struct that accept any struct (interface{}) to perform validation
func (v *Validate) ValidateStruct(structType interface{}) []string {
	msg := []string{}

	val := v.validateStruct(structType)
	for _, v := range val {
		msg = append(msg, v.Message)
	}

	return msg
}

func (v *Validate) translateErrMessage(err error) []ErrorValidation {
	var errs []ErrorValidation

	for _, err := range err.(validator.ValidationErrors) {
		var errorMessage ErrorValidation
		switch err.Tag() {

		case "datetime":
			errorMessage.Message = fmt.Sprintf(`%s format is not valid`,
				err.Field())
		default:
			errorMessage.Message = err.Translate(v.trans)
		}
		errs = append(errs, errorMessage)
	}

	return errs
}

// validateIsDatetime is the validation function for validating if the current field's value is a valid datetime string.
func validateIsDatetime(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	if field.Kind() == reflect.String {
		_, err := time.Parse(param, field.String())

		return err == nil
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}
