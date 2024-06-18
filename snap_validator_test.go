package snap_validator

import (
	"errors"
	"fmt"
	"github.com/apelweb15/snap-validator/snap_validator_errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Request struct {
	TotalAmount         TotalAmount  `json:"totalAmount" snapValidator:"required"`
	PartnerServiceId    string       `json:"partnerServiceId" snapValidator:"required|max_length:8|min_length:8|numeric"`
	CustomerNo          string       `json:"customerNo" snapValidator:"required|min_length:8|max_length:20|numeric"`
	ExpiredDate         string       `json:"expiredDate" snapValidator:"required|max_length:25|min_length:25|iso_date|after_time_now"`
	TrxId               string       `json:"trxId" snapValidator:"required|min_length:1|max_length:64|alpha_numeric_symbol"`
	VirtualAccountEmail string       `json:"virtualAccountEmail" snapValidator:"email|max_length:50"`
	BillDetails         []BillDetail `json:"billDetails"`
}

type BillDetail struct {
	BillCode        string             `json:"billCode" snapValidator:"required|max_length:2|numeric" label:"{dynamic}"`
	BillNo          string             `json:"billNo" snapValidator:"max_length:18|numeric" label:"{dynamic}"`
	BillName        string             `json:"billName" snapValidator:"max_length:20|string" label:"{dynamic}"`
	BillShortName   string             `json:"billShortName" snapValidator:"max_length:10|string" label:"{dynamic}"`
	BillDescription MessageDescription `json:"billDescription"`
	BillAmount      TotalAmount        `json:"billAmount" snapValidator:"required"`
}

type MessageDescription struct {
	Indonesia string `json:"indonesia" validateSnap:"required|max_length:18|string" label:"{dynamic}"`
	English   string `json:"english" validateSnap:"required|max_length:18|string" label:"{dynamic}"`
}

type TotalAmount struct {
	Value    string `json:"value" snapValidator:"required|amount" label:"{dynamic}"`
	Currency string `json:"currency" snapValidator:"required|in_data:[IDR]" label:"{dynamic}"`
}

func TestValidator(t *testing.T) {
	req := &Request{
		PartnerServiceId: "   1114",
		CustomerNo:       "1231231231",
		TotalAmount: TotalAmount{
			Value:    "25000.00",
			Currency: "IDR",
		},
		ExpiredDate:         "2024-12-31T23:59:59+07:00",
		TrxId:               "ishdfi-8u8u9kk",
		VirtualAccountEmail: "acd@ema",
	}
	v := New()
	err := v.ValidateStruct(*req, "25")
	if err != nil {
		fmt.Println("err:", err)
	}
	var errorRes *snap_validator_errors.ErrorValidation
	if errors.As(err, &errorRes) {
		fmt.Println("Error Message: ", errorRes.Message)
		fmt.Println("Error Code: ", errorRes.Code)
		fmt.Println("Error SnapCode: ", errorRes.SnapCode)
	}
	assert.Error(t, err)
}

func TestRequired(t *testing.T) {
	v := New()
	//validatorService := validator.NewService[any]()
	errString := v.validator.ValidateRequired("")
	fmt.Println(errString)
	assert.Error(t, errString, "String required successfully", errString)

	errStructPointer := v.validator.ValidateRequired(&Request{})
	fmt.Println(errStructPointer)
	assert.Error(t, errStructPointer, "Struct pointer required successfully")

	errStruct := v.validator.ValidateRequired(Request{})
	fmt.Println(errStruct)
	assert.Error(t, errStruct, "Struct required successfully")

	errSlice := v.validator.ValidateRequired([]int{})
	fmt.Println(errSlice)
	assert.Error(t, errSlice, "Slice required successfully")

	errNumber := v.validator.ValidateRequired(0)
	fmt.Println(errNumber)
	assert.Error(t, errNumber, "Number required successfully")
}
