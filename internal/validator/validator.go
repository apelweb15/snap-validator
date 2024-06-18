package validator

import "github.com/apelweb15/snap-validator/internal/models"

// Validator /*Available Validation
/*
required = [string, struct, slice, numeric = 0]
min_length = [string]
max_length = [string]
iso_date = [string]
after_time_now = [string] comparing date must more than large than now
alpha_numeric = [string]
alpha_numeric_symbol = [string] comparing alphanumeric with dash and underscore
numeric = [string]
string = [string]
amount = [string] = ISO4217
email = [string]
in_data = [string] = comparing value with defined list data
*/
type Validator[T any] interface {
	ValidateStructSnap(any interface{}, parentProperty ...models.ValidatorProperty) error
	ValidateStructSnapServiceCode(any interface{}, serviceCode string, parentProperty ...models.ValidatorProperty) error
	ValidateRequired(data interface{}) error
	ValidateMaxLength(length int, data interface{}) error
	ValidateMinLength(length int, data interface{}) error
	ValidateIsoDate(data interface{}) error
	ValidateAfterTimeNow(data interface{}) error
	ValidateAlphaNum(data interface{}) error
	ValidateAlphaNumSymbol(data interface{}) error
	ValidateNumeric(data interface{}) error
	ValidateString(data interface{}) error
	ValidateAmount(data interface{}) error
	ValidateInData(allowed []any, data interface{}) error
	ValidateEmail(data interface{}) error
	ValidateUrl(data interface{}) error
}
