package main

import (
	"errors"
	"fmt"
)

var LoginError = errors.New("login error")

func errThrower() error {
	return fmt.Errorf("error %w", LoginError)
}

func main() {
	err := errThrower()
	if errors.Is(err, LoginError) {
		fmt.Println(err)
	}
}
