package utils

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func LoadSecrets(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return errors.New("invalid .secret format")
		}
		key := parts[0]
		value := parts[1]
		os.Setenv(key, value) // set as environment variable
	}

	return scanner.Err()
}
