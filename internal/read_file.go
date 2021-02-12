package internal

import (
	"io/ioutil"
)

/*
read file.

*/

func ReadFile(filePath string) (*[]byte, error) {
	bodyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return &bodyBytes, nil
}
