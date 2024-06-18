package snap_validator_errors

type ErrorValidation struct {
	Code      string
	SnapCode  string
	Message   string
	FieldName string
}

func (e *ErrorValidation) Error() string {
	return e.Message
}
func NewErrorSnap(code string, label string, fieldName string, serviceCode string) *ErrorValidation {
	snapCode := code
	if len(code) > 4 {
		snapCode = code[:3] + serviceCode + code[3:5]
	}
	return &ErrorValidation{
		Code:      code,
		SnapCode:  snapCode,
		Message:   GetSnapMessage(code) + " " + label,
		FieldName: fieldName,
	}
}

func NewError(code string, fieldName string) *ErrorValidation {
	return &ErrorValidation{
		Code:      code,
		SnapCode:  code,
		Message:   GetSnapMessage(code),
		FieldName: fieldName,
	}
}

func GetSnapMessage(code string) string {
	message := ""
	switch code {
	case "20000":
		message = "Successful"
	case "20200":
		message = "Request In Progress"
	case "40000":
		message = "Bad Request"
	case "40001":
		message = "Invalid Field Format"
	case "40002":
		message = "Missing Mandatory Field"
	case "40100":
		message = "Unauthorized"
	case "40101":
		message = "Invalid Token (B2B)"
	case "40102":
		message = "Invalid Customer Token"
	case "40103":
		message = "Token Not Found (B2B)"
	case "40104":
		message = "Customer Token Not Found"
	case "40300":
		message = "Transaction Expired"
	case "40301":
		message = "Feature Not Allowed"
	case "40302":
		message = "Exceeds Transaction Amount Limit"
	case "40303":
		message = "Suspected Fraud"
	case "40304":
		message = "Activity Count Limit Exceeded"
	case "40305":
		message = "Do Not Honor"
	case "40306":
		message = "Feature Not Allowed At This Time"
	case "40307":
		message = "Card Blocked"
	case "40308":
		message = "Card Expired"
	case "40309":
		message = "Dormant Account"
	case "40310":
		message = "Need To Set Token Limit"
	case "40311":
		message = "OTP Blocked"
	case "40312":
		message = "OTP Lifetime Expired"
	case "40313":
		message = "OTP Sent To Cardholer"
	case "40314":
		message = "Insufficient Funds"
	case "40315":
		message = "Transaction Not Permitted"
	case "40316":
		message = "Suspend Transaction"
	case "40317":
		message = "Token Limit Exceeded"
	case "40318":
		message = "Inactive Card/Account/Customer"
	case "40319":
		message = "Merchant Blacklisted"
	case "40320":
		message = "Merchant Limit Exceed"
	case "40321":
		message = "Set Limit Not Allowed"
	case "40322":
		message = "Token Limit Invalid"
	case "40323":
		message = "Account Limit Exceed"
	case "40400":
		message = "Invalid Transaction Status"
	case "40401":
		message = "Transaction Not Found"
	case "40402":
		message = "Invalid Routing"
	case "40403":
		message = "Bank Not Supported By Switch"
	case "40404":
		message = "Transaction Cancelled"
	case "40405":
		message = "Merchant Is Not Registered For Card Registration Services"
	case "40406":
		message = "Need To Request OTP"
	case "40407":
		message = "Journey Not Found"
	case "40408":
		message = "Invalid Merchant"
	case "40409":
		message = "No Issuer"
	case "40410":
		message = "Invalid API Transition"
	case "40411":
		message = "Invalid Card/Account/Customer/Virtual Account"
	case "40412":
		message = "Invalid Bill/Virtual Account"
	case "40413":
		message = "Invalid Amount"
	case "40414":
		message = "Paid Bill"
	case "40415":
		message = "Invalid OTP"
	case "40416":
		message = "Partner Not Found"
	case "40417":
		message = "Invalid Terminal"
	case "40418":
		message = "Inconsistent Request"
	case "40419":
		message = "Invalid Bill/Virtual Account"
	case "40500":
		message = "Requested Function Is Not Supported"
	case "40501":
		message = "Requested Opearation Is Not Allowed"
	case "40900":
		message = "Conflict"
	case "40901":
		message = "Duplicate partnerReferenceNo"
	case "42900":
		message = "Too Many Requests"
	case "50000":
		message = "General Error"
	case "50001":
		message = "Internal Server Error"
	case "50002":
		message = "External Server Error"
	case "50400":
		message = "Timeout"

	default:
		message = "General Error"
	}
	return message

}
