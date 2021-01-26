package error

import (
	"errors"
	"fmt"
	"os"
)

func ExitAndError(errMessage string) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", errors.New(errMessage))
	os.Exit(0)
}
