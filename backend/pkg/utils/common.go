package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func LoadEnvConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)

		// if starts with #, skip
		if strings.HasPrefix(line, "#") {
			continue
		}

		if len(parts) != 2 {
			return fmt.Errorf("%s invalid line: %s", filename, line)
		}
		key := parts[0]
		value := parts[1]
		os.Setenv(key, value) // set as environment variable
	}

	return scanner.Err()
}

func GetContractABI(path string) (*abi.ABI, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	return &parsedABI, nil
}

func IsValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
