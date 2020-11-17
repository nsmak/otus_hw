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

const tagKey = "validate"

var (
	ErrInputIsNotStruct = errors.New("input value is not a struct")
	ErrInvalidLength    = errors.New("length is invalid")
	ErrNotMatchRegexp   = errors.New("string is not matched by regexp")
	ErrNotIncludedInSet = errors.New("not included in validation set")
	ErrLessThanMin      = errors.New("less than the minimum")
	ErrMaxMoreMax       = errors.New("more than maximum")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for _, e := range v {
		builder.WriteString("field: ")
		builder.WriteString(e.Field)
		builder.WriteString(" - ")
		builder.WriteString(e.Err.Error())
	}
	return builder.String()
}

func Validate(v interface{}) error {
	var vErrors ValidationErrors
	vv := reflect.ValueOf(v)

	if vv.Kind() != reflect.Struct {
		return ValidationErrors{
			ValidationError{
				Field: "-",
				Err:   ErrInputIsNotStruct,
			},
		}
	}

	t := vv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagKey)
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
		return nil
	}
}

func stringValidators(tag, field, value string) []stringValidator {
	validatorsRaw := strings.Split(tag, "|")
	validators := make([]stringValidator, len(validatorsRaw))

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
				Err:   fmt.Errorf("%w: need %d, now %d", ErrInvalidLength, cond, l),
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
				Err:   fmt.Errorf("%w: %s", ErrNotMatchRegexp, sv.condition),
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
			Err:   fmt.Errorf("%w: value %s, set %v", ErrNotIncludedInSet, sv.value, set),
		}

	default:
		log.Printf("unknown validator's name %s", sv.name)
	}
	return nil
}

func intValidators(tag, field string, value int) []intValidator {
	validatorsRaw := strings.Split(tag, "|")
	validators := make([]intValidator, len(validatorsRaw))

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
				Err:   fmt.Errorf("%w: value %d, condition %d", ErrLessThanMin, iv.value, cond),
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
				Err:   fmt.Errorf("%w: value %d, condition %d", ErrMaxMoreMax, iv.value, cond),
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
			Err:   fmt.Errorf("%w: value %d, set %v", ErrNotIncludedInSet, iv.value, set),
		}
	default:
		log.Printf("unknown validator's name %s", iv.name)
	}
	return nil
}
