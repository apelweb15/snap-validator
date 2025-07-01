package validator

import (
	"fmt"
	"github.com/apelweb15/snap-validator/internal/models"
	"github.com/apelweb15/snap-validator/snap_validator_errors"
	"github.com/apelweb15/snap-validator/snap_validator_utils"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type validatorImpl[T any] struct {
	customValidation []map[string]string
	serviceCode      string
}

func New[T any](serviceCode string, customValidation ...map[string]string) Validator[T] {
	return &validatorImpl[T]{
		serviceCode:      serviceCode,
		customValidation: customValidation,
	}
}

func NewService[T any](customValidation ...map[string]string) Validator[T] {
	return validatorImpl[T]{
		customValidation: customValidation,
	}
}

func (v validatorImpl[T]) ValidateStructSnap(data interface{}, parentProperty ...models.ValidatorProperty) error {
	return v.validate(data, parentProperty...)
}

func (v validatorImpl[T]) ValidateStructSnapServiceCode(data interface{}, serviceCode string, parentProperty ...models.ValidatorProperty) error {
	v.serviceCode = serviceCode
	return v.validate(data, parentProperty...)
}

func (v validatorImpl[T]) validate(data interface{}, parentProperty ...models.ValidatorProperty) error {
	reflectType := reflect.TypeOf(data)
	reflectValue := reflect.ValueOf(data)
	for i := 0; i < reflectType.NumField(); i++ {
		fieldType := reflectType.Field(i)
		fieldValue := reflectValue.Field(i)
		errValidation := v.process(models.CustomValidator{
			FieldType:  fieldType,
			FieldValue: fieldValue,
		}, parentProperty...)
		if errValidation != nil {
			return errValidation
		}

		if fieldType.Type.Kind() == reflect.Struct {
			var f = reflect.ValueOf(fieldValue.Interface())
			if !reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(f.Type()).Interface()) {
				parentProperty := append(parentProperty, models.ValidatorProperty{
					FieldName:  fieldType.Name,
					JsonName:   fieldType.Tag.Get("json"),
					Array:      false,
					FieldType:  fieldType,
					FieldValue: fieldValue,
				})
				errStruct := v.validate(fieldValue.Interface(), parentProperty...)
				if errStruct != nil {
					return errStruct
				}
			}
		}

		if fieldType.Type.Kind() == reflect.Slice {
			var f = reflect.ValueOf(fieldValue.Interface())
			ret := make([]reflect.Value, f.Len())
			for j := 0; j < f.Len(); j++ {
				ret[j] = f.Index(j)

				parentProperty := append(parentProperty, models.ValidatorProperty{
					FieldName:  fieldType.Name,
					JsonName:   fieldType.Tag.Get("json"),
					Array:      true,
					FieldType:  fieldType,
					FieldValue: fieldValue,
					IdxArray:   j,
				})

				errSlice := v.validate(f.Index(j).Interface(), parentProperty...)
				if errSlice != nil {
					return errSlice
				}
			}
		}
	}
	return nil
}

func (v validatorImpl[T]) process(customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	fieldType := customValidator.FieldType
	//value := customValidator.FieldValue

	snapTag := fieldType.Tag.Get("snapValidator")
	listValidator := strings.Split(snapTag, "|")
	for _, validation := range listValidator {
		//region get validation type & value
		validations := strings.Split(validation, ":")
		validationParam := ""
		validationType := validations[0]
		if len(validations) == 2 {
			validationParam = validations[1]
		}

		//endregion get validation type & value

		switch validationType {
		case "required":
			errorValidate := v.validateRequired(customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break
		case "min_length":
			errorValidate := v.validateMinLength(validationParam, customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break
		case "max_length":
			errorValidate := v.validateMaxLength(validationParam, customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break
		case "iso_date":
			errorValidate := v.validateIsoDate(customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break
		case "after_time_now":
			errorValidate := v.validateAfterTimeNow(customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break
		case "alpha_numeric":
			errorValidate := v.validateAlphaNum(customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break
		case "alpha_numeric_symbol":
			errorValidate := v.validateAlphaNumSymbol(customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break
		case "numeric":
			errorValidate := v.validateNumeric(customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break

		case "string":
			errorValidate := v.validateString(customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break
		case "amount":
			errorValidate := v.validateAmount(customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break
		case "email":
			errorValidate := v.validateEmail(customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break
		case "in_data":
			errorValidate := v.validateInData(validationParam, customValidator, parentProperty...)
			if errorValidate != nil {
				return errorValidate
			}
			break

		}

	}
	return nil
}

func (v validatorImpl[T]) ValidateRequired(data interface{}) error {
	value := reflect.ValueOf(data)
	isValid := true
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			isValid = false
		}
	} else if value.Kind() == reflect.Struct {
		if reflect.DeepEqual(data, reflect.Zero(value.Type()).Interface()) {
			isValid = false
		}
	} else if snap_validator_utils.KindIsNumeric(value.Kind()) {
		n, _ := strconv.ParseFloat(fmt.Sprintf("%d", value.Interface()), 64)
		if n <= 0 {
			isValid = false
		}
	} else if value.Kind() == reflect.Slice {
		if reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface()) {
			isValid = false
		}
		if value.Len() == 0 {
			isValid = false
		}
	} else if value.Kind() == reflect.Ptr {
		rawType := reflect.TypeOf(data).Elem()
		if reflect.DeepEqual(data, reflect.New(rawType).Interface()) {
			isValid = false
		}
	}
	if !isValid {
		return &snap_validator_errors.ErrorValidation{
			Code:    "40002",
			Message: "Invalid mandatory field",
		}
	}
	return nil
}

func (v validatorImpl[T]) validateRequired(customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	value := customValidator.FieldValue
	errorCode := "40002"
	err := v.ValidateRequired(value.Interface())
	if err != nil {
		return v.parsingError(customValidator, errorCode, "required", parentProperty...)
	}
	return nil
}

func (v validatorImpl[T]) ValidateMaxLength(length int, data interface{}) error {
	value := reflect.ValueOf(data)
	lengthStr := 0
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}

		lengthStr = len([]rune(data.(string)))
	}
	if lengthStr > length {
		return &snap_validator_errors.ErrorValidation{
			Code:    "40001",
			Message: "Invalid field format",
		}
	}
	return nil
}

func (v validatorImpl[T]) validateMaxLength(length string, customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	fieldType := customValidator.FieldType
	fieldValue := customValidator.FieldValue
	errorCode := "40001"
	if length == "" {
		length = "0"
	}
	lengthInt, err := strconv.Atoi(length)
	if err != nil {
		return snap_validator_errors.NewError("500000", fieldType.Name)
	}
	errValidation := v.ValidateMaxLength(lengthInt, fieldValue.Interface())
	if errValidation != nil {
		return v.parsingError(customValidator, errorCode, "max_length", parentProperty...)
	}
	return nil
}

func (v validatorImpl[T]) ValidateMinLength(length int, data interface{}) error {
	value := reflect.ValueOf(data)
	lengthStr := 0
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}

		lengthStr = len([]rune(data.(string)))
	}
	if lengthStr < length {
		return &snap_validator_errors.ErrorValidation{
			Code:    "40001",
			Message: "Invalid field format",
		}
	}
	return nil
}

func (v validatorImpl[T]) validateMinLength(length string, customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	fieldType := customValidator.FieldType
	fieldValue := customValidator.FieldValue
	errorCode := "40001"
	if length == "" {
		length = "0"
	}
	lengthInt, err := strconv.Atoi(length)
	if err != nil {
		return snap_validator_errors.NewError("500000", fieldType.Name)
	}
	errValidation := v.ValidateMinLength(lengthInt, fieldValue.Interface())
	if errValidation != nil {
		return v.parsingError(customValidator, errorCode, "min_length", parentProperty...)
	}
	return nil
}

func (v validatorImpl[T]) ValidateIsoDate(data interface{}) error {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}
		_, errDate := time.Parse("2006-01-02T15:04:05-07:00", data.(string))
		if errDate != nil {
			return &snap_validator_errors.ErrorValidation{
				Code:    "40001",
				Message: "Invalid field format",
			}
		}
	}
	return nil
}

func (v validatorImpl[T]) validateIsoDate(customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	errorCode := "40001"
	fieldValue := customValidator.FieldValue

	err := v.ValidateIsoDate(fieldValue.Interface())
	if err != nil {
		return v.parsingError(customValidator, errorCode, "iso_date", parentProperty...)
	}
	return nil
}

func (v validatorImpl[T]) ValidateAfterTimeNow(data interface{}) error {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}
		tm, errDate := time.Parse("2006-01-02T15:04:05-07:00", data.(string))
		if errDate != nil {
			return &snap_validator_errors.ErrorValidation{
				Code:    "40001",
				Message: "Invalid field format",
			}
		}
		if time.Now().After(tm) {
			return &snap_validator_errors.ErrorValidation{
				Code:    "40001",
				Message: "Invalid field format date must greater than now",
			}
		}
	}
	return nil
}

func (v validatorImpl[T]) validateAfterTimeNow(customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	errorCode := "40001"
	fieldValue := customValidator.FieldValue
	err := v.ValidateAfterTimeNow(fieldValue.Interface())
	if err != nil {
		return v.parsingErrorWithSuffix("must greater than now", customValidator, errorCode, "after_time_now", parentProperty...)
	}
	return nil
}

func (v validatorImpl[T]) ValidateAlphaNum(data interface{}) error {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}
		cond := regexp.MustCompile(`^[a-zA-Z\d]+$`)
		if !cond.MatchString(data.(string)) {
			return &snap_validator_errors.ErrorValidation{
				Code:    "40001",
				Message: "Invalid field format",
			}
		}

	}
	return nil
}

func (v validatorImpl[T]) validateAlphaNum(customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	errorCode := "40001"
	fieldValue := customValidator.FieldValue

	err := v.ValidateAlphaNum(fieldValue.Interface())
	if err != nil {
		return v.parsingError(customValidator, errorCode, "alpha_numeric", parentProperty...)
	}

	return nil
}

func (v validatorImpl[T]) ValidateAlphaNumSymbol(data interface{}) error {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}

		cond := regexp.MustCompile(`^[a-zA-Z\d_-]+$`)
		if !cond.MatchString(data.(string)) {
			return &snap_validator_errors.ErrorValidation{
				Code:    "40001",
				Message: "Invalid field format",
			}
		}

	}
	return nil
}

func (v validatorImpl[T]) validateAlphaNumSymbol(customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	errorCode := "40001"
	fieldValue := customValidator.FieldValue
	err := v.ValidateAlphaNumSymbol(fieldValue.Interface())
	if err != nil {
		return v.parsingError(customValidator, errorCode, "alpha_numeric_symbol", parentProperty...)
	}

	return nil
}

func (v validatorImpl[T]) ValidateNumeric(data interface{}) error {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}
		cond := regexp.MustCompile(`^\d+$`)
		if !cond.MatchString(strings.TrimSpace(data.(string))) {
			return &snap_validator_errors.ErrorValidation{
				Code:    "40001",
				Message: "Invalid field format",
			}
		}
	}
	return nil
}

func (v validatorImpl[T]) validateNumeric(customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	errorCode := "40001"
	fieldValue := customValidator.FieldValue
	err := v.ValidateNumeric(fieldValue.Interface())
	if err != nil {
		return v.parsingError(customValidator, errorCode, "numeric", parentProperty...)
	}
	return nil
}

func (v validatorImpl[T]) ValidateString(data interface{}) error {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}
		cond := regexp.MustCompile(`^[a-zA-Z\d:; )(\+\/\&\#\!\@\%\*.,\'\"_-]+$`)
		if !cond.MatchString(data.(string)) {
			return &snap_validator_errors.ErrorValidation{
				Code:    "40001",
				Message: "Invalid field format",
			}
		}
	}
	return nil
}

func (v validatorImpl[T]) validateString(customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	errorCode := "40001"
	fieldValue := customValidator.FieldValue

	err := v.ValidateString(fieldValue.Interface())
	if err != nil {
		return v.parsingError(customValidator, errorCode, "string", parentProperty...)
	}
	return nil
}

func (v validatorImpl[T]) ValidateAmount(data interface{}) error {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}
		cond := regexp.MustCompile(`^(0|[1-9]\d{0,10})[.]0{2}$`)
		if !cond.MatchString(data.(string)) {
			return &snap_validator_errors.ErrorValidation{
				Code:    "40001",
				Message: "Invalid field format",
			}
		}
	}
	return nil
}

func (v validatorImpl[T]) validateAmount(customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	errorCode := "40001"
	fieldValue := customValidator.FieldValue
	err := v.ValidateAmount(fieldValue.Interface())
	if err != nil {
		return v.parsingError(customValidator, errorCode, "amount", parentProperty...)
	}
	return nil
}

func (v validatorImpl[T]) ValidateInData(allowed []any, data interface{}) error {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}
		exist, _ := snap_validator_utils.InArray(data.(string), allowed)
		if !exist {
			return &snap_validator_errors.ErrorValidation{
				Code:    "40001",
				Message: "Invalid field format",
			}
		}
	}
	return nil
}

func (v validatorImpl[T]) validateInData(allowed string, customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	errorCode := "40001"
	fieldType := customValidator.FieldType
	fieldValue := customValidator.FieldValue
	if fieldType.Type.Kind() == reflect.String {
		if fieldValue.Interface().(string) == "" {
			return nil
		}
		allowed = strings.ReplaceAll(allowed, "[", "")
		allowed = strings.ReplaceAll(allowed, "]", "")
		tSlice := strings.Split(allowed, ",")
		exist, _ := snap_validator_utils.InArray(fieldValue.Interface().(string), tSlice)
		if !exist {
			return v.parsingError(customValidator, errorCode, "in_data", parentProperty...)
		}
	}
	return nil
}

func (v validatorImpl[T]) ValidateEmail(data interface{}) error {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}
		_, err := mail.ParseAddress(data.(string))
		if err != nil {
			return &snap_validator_errors.ErrorValidation{
				Code:    "40001",
				Message: "Invalid field format",
			}
		}
	}
	return nil
}

func (v validatorImpl[T]) validateEmail(customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	errorCode := "40001"
	fieldValue := customValidator.FieldValue

	err := v.ValidateEmail(fieldValue.Interface())
	if err != nil {
		return v.parsingError(customValidator, errorCode, "email", parentProperty...)
	}
	return nil
}

func (v validatorImpl[T]) ValidateUrl(data interface{}) error {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.String {
		if data.(string) == "" {
			return nil
		}
		_, err := url.ParseRequestURI(data.(string))
		if err != nil {
			return &snap_validator_errors.ErrorValidation{
				Code:    "40001",
				Message: "Invalid field format",
			}
		}
	}
	return nil
}

func (v validatorImpl[T]) validateUrl(customValidator models.CustomValidator, parentProperty ...models.ValidatorProperty) error {
	errorCode := "40001"
	fieldValue := customValidator.FieldValue

	err := v.ValidateUrl(fieldValue.Interface())
	if err != nil {
		return v.parsingError(customValidator, errorCode, "url", parentProperty...)
	}
	return nil
}

func (v validatorImpl[T]) parsingErrorWithSuffix(suffix string, customValidator models.CustomValidator, code string, validatorKey string, parentProperty ...models.ValidatorProperty) error {
	fieldType := customValidator.FieldType
	label := v.getLabel(fieldType, parentProperty...)

	if len(v.customValidation) > 0 {
		customValidation := v.customValidation[0]
		customMessage := customValidation[fieldType.Name+":"+validatorKey]
		var indexIfUseArray []any
		if len(parentProperty) > 0 {
			keyLabel := ""
			for _, parent := range parentProperty {
				if keyLabel != "" {
					keyLabel += "."
				}
				keyLabel += parent.FieldName
				if parent.Array {
					keyLabel += ".*"
					indexIfUseArray = append(indexIfUseArray, parent.IdxArray)
				}
			}
			if keyLabel != "" {
				keyLabel += "." + fieldType.Name
			}
			customMessageParent := customValidation[keyLabel+":"+validatorKey]
			if customMessageParent != "" {
				if len(indexIfUseArray) > 0 {
					customMessage = fmt.Sprintf(customMessageParent, indexIfUseArray...)
				} else {
					customMessage = customMessageParent
				}

				//
			}
		}
		if customMessage != "" {
			snapCode := code
			if len(code) > 4 {
				snapCode = code[:3] + v.serviceCode + code[3:5]
			}
			return &snap_validator_errors.ErrorValidation{
				Code:      code,
				SnapCode:  snapCode,
				Message:   customMessage,
				FieldName: fieldType.Name,
			}
		}
	}

	return snap_validator_errors.NewErrorSnap(code, label+" "+suffix, fieldType.Name, v.serviceCode)
	//return error_snap.NewErrorSnap(code, label, fieldType.Name, v.serviceCode)
}

func (v validatorImpl[T]) parsingError(customValidator models.CustomValidator, code string, validatorKey string, parentProperty ...models.ValidatorProperty) error {
	fieldType := customValidator.FieldType
	label := v.getLabel(fieldType, parentProperty...)

	if len(v.customValidation) > 0 {
		customValidation := v.customValidation[0]
		customMessage := customValidation[fieldType.Name+":"+validatorKey]
		var indexIfUseArray []any
		if len(parentProperty) > 0 {
			keyLabel := ""
			for _, parent := range parentProperty {
				if keyLabel != "" {
					keyLabel += "."
				}
				keyLabel += parent.FieldName
				if parent.Array {
					keyLabel += ".*"
					indexIfUseArray = append(indexIfUseArray, parent.IdxArray)
				}
			}
			if keyLabel != "" {
				keyLabel += "." + fieldType.Name
			}
			customMessageParent := customValidation[keyLabel+":"+validatorKey]
			if customMessageParent != "" {
				if len(indexIfUseArray) > 0 {
					customMessage = fmt.Sprintf(customMessageParent, indexIfUseArray...)
				} else {
					customMessage = customMessageParent
				}
			}
		}
		if customMessage != "" {
			snapCode := code
			if len(code) > 4 {
				snapCode = code[:3] + v.serviceCode + code[3:5]
			}
			return &snap_validator_errors.ErrorValidation{
				Code:      code,
				SnapCode:  snapCode,
				Message:   customMessage,
				FieldName: fieldType.Name,
			}
		}
	}

	return snap_validator_errors.NewErrorSnap(code, label, fieldType.Name, v.serviceCode)
	//return error_snap.NewErrorSnap(code, label, fieldType.Name, v.serviceCode)
}

func (v validatorImpl[T]) getLabel(fieldType reflect.StructField, parentProperty ...models.ValidatorProperty) string {
	label := fieldType.Name
	jsonTag := fieldType.Tag.Get("json")
	if jsonTag != "" {
		label = jsonTag
	}
	labelWithParent := ""
	if len(parentProperty) > 0 {
		for _, labelParent := range parentProperty {
			if labelParent.JsonName != "" {
				labelWithParent += labelParent.JsonName
				labelWithParent += "."
				if labelParent.Array {
					labelWithParent += fmt.Sprintf("%d.", labelParent.IdxArray)
				}
			} else if labelParent.FieldName != "" {
				labelWithParent += labelParent.FieldName
				labelWithParent += "."
				if labelParent.Array {
					labelWithParent += fmt.Sprintf("%d.", labelParent.IdxArray)
				}
			}

		}
		return labelWithParent + label
	}
	return label
}
