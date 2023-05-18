package helpers

import (
	"bufio"
	"os"
	"strings"
)

func ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	str = strings.TrimSuffix(str, "\n")
	return str, nil
}
