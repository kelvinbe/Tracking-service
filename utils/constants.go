package utils

import "os"

var AFRICAS_TALKING_SEND_URL = (func (env string)string {
	if env == "production" {
		return "https://api.africastalking.com/version1/messaging"
	}
	return "https://api.sandbox.africastalking.com/version1/messaging"
})(os.Getenv("APP_ENV"))

const ACTIVATION_CODE = "777"
const DEACTIVATION_CODE = "444"