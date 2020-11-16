package hw09_struct_validator //nolint:golint,stylecheck
import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errStr := ""
	for _, e := range v {
		errStr += fmt.Sprintf("field: %s - %s", e.Field, e.Err.Error())
	}
	return errStr
}

func Validate(v interface{}) error {
	var vErrors ValidationErrors

	vv := reflect.ValueOf(v)

	if vv.Kind() != reflect.Struct {
		return errors.New("input value is not a struct")
	}

	t := vv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			log.Println("no one tag for key 'validate'")
			continue
		}
		fv := vv.Field(i)

		if !fv.CanInterface() {
			continue
		}

		errs := validateValue(tag, field.Name, fv)
		if len(errs) > 0 {
			vErrors = append(vErrors, errs...)
		}
	}

	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}

func validateValue(tag, field string, v reflect.Value) ValidationErrors {
	switch v.Kind() { //nolint:exhaustive  // есть default case
	case reflect.String:
		var errs ValidationErrors
		validators := stringValidators(tag, field, v.String())
		for _, validator := range validators {
			err := validator.validate()
			if err != nil {
				log.Println(err.Err.Error())
				errs = append(errs, *err)
			}
		}
		return errs

	case reflect.Int:
		var errs ValidationErrors
		validators := intValidators(tag, field, int(v.Int()))
		for _, validator := range validators {
			err := validator.Validate()
			if err != nil {
				errs = append(errs, *err)
			}
		}
		return errs

	case reflect.Slice:
		var errs ValidationErrors
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			errs = append(errs, validateValue(tag, field, elem)...)
		}
		return errs

	default:
		log.Println("unsupported type")
		return ValidationErrors{}
	}
}

func stringValidators(tag, field, value string) []stringValidator {
	validators := make([]stringValidator, 0)

	validatorsRaw := strings.Split(tag, "|")
	for _, valRaw := range validatorsRaw {
		val := strings.Split(valRaw, ":")
		if len(val) != 2 {
			log.Printf("invalid value for tag %s\n", tag)
			continue
		}
		sVal := stringValidator{
			name:      val[0],
			condition: val[1],
			field:     field,
			value:     value,
		}
		validators = append(validators, sVal)
	}
	return validators
}

type stringValidator struct {
	name      string
	condition string
	field     string
	value     string
}

func (sv stringValidator) validate() *ValidationError {
	switch sv.name {
	case "len":
		cond, err := strconv.Atoi(sv.condition)
		if err != nil {
			log.Println("len value is not int")
			return nil
		}
		l := len(sv.value)
		if l != cond {
			return &ValidationError{
				Field: sv.field,
				Err:   fmt.Errorf("length is invalid need: %d, now: %d", cond, l),
			}
		}

	case "regexp":
		matched, err := regexp.MatchString(sv.condition, sv.value)
		if err != nil {
			return &ValidationError{Field: sv.field, Err: err}
		}
		if !matched {
			return &ValidationError{
				Field: sv.field,
				Err:   fmt.Errorf("string is not matched by regexp"),
			}
		}

	case "in":
		set := strings.Split(sv.condition, ",")
		for _, e := range set {
			if sv.value == e {
				return nil
			}
		}
		return &ValidationError{
			Field: sv.field,
			Err:   fmt.Errorf("value (%s) not included in validation set (%v)", sv.value, set),
		}

	default:
		log.Printf("unknown validator's name %s", sv.name)
	}
	return nil
}

func intValidators(tag, field string, value int) []intValidator {
	validators := make([]intValidator, 0)

	validatorsRaw := strings.Split(tag, "|")
	for _, valRaw := range validatorsRaw {
		val := strings.Split(valRaw, ":")
		if len(val) != 2 {
			log.Printf("invalid value for tag %s\n", tag)
			continue
		}
		sVal := intValidator{
			name:      val[0],
			condition: val[1],
			field:     field,
			value:     value,
		}
		validators = append(validators, sVal)
	}
	return validators
}

type intValidator struct {
	name      string
	condition string
	field     string
	value     int
}

func (iv intValidator) Validate() *ValidationError {
	switch iv.name {
	case "min":
		cond, err := strconv.Atoi(iv.condition)
		if err != nil {
			log.Println("min value is not int")
			return nil
		}
		if iv.value < cond {
			return &ValidationError{
				Field: iv.field,
				Err:   fmt.Errorf("value (%d) is less thаn condition (%d)", iv.value, cond),
			}
		}

	case "max":
		cond, err := strconv.Atoi(iv.condition)
		if err != nil {
			log.Println("max value is not int")
			return nil
		}
		if iv.value > cond {
			return &ValidationError{
				Field: iv.field,
				Err:   fmt.Errorf("value (%d) is more thаn condition (%d)", iv.value, cond),
			}
		}

	case "in":
		set := strings.Split(iv.condition, ",")
		for _, e := range set {
			intVal, err := strconv.Atoi(e)
			if err != nil {
				log.Println("set's value is not int")
				return nil
			}
			if iv.value == intVal {
				return nil
			}
		}
		return &ValidationError{
			Field: iv.field,
			Err:   fmt.Errorf("value (%d) not included in validation set (%v)", iv.value, set),
		}
	default:
		log.Printf("unknown validator's name %s", iv.name)
	}
	return nil
}
