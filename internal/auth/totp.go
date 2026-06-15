package auth

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateTOTPKey(
	username string,
) (*otp.Key, error) {

	return totp.Generate(
		totp.GenerateOpts{
			Issuer:      "CLI Auth",
			AccountName: username,
		},
	)
}

func ValidateTOTP(
	code string,
	secret string,
) bool {

	return totp.Validate(
		code,
		secret,
	)
}
