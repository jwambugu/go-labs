package config

import (
	"fmt"
	"os"
)

func GetPWD() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get pwd: %v", err)
	}
	return pwd, nil
}
