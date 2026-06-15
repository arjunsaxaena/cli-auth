package auth

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(
	username string,
	otpURL string,
) (string, error) {

	err := os.MkdirAll(
		"qr-codes",
		0755,
	)

	if err != nil {
		return "", err
	}

	filePath := filepath.Join(
		"qr-codes",
		fmt.Sprintf("%s.png", username),
	)

	err = qrcode.WriteFile(
		otpURL,
		qrcode.Medium,
		256,
		filePath,
	)

	if err != nil {
		return "", err
	}

	return filePath, nil
}
