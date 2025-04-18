package pkg

import (
	"errors"
	"log"
	"runtime/debug"
	"strings"
	"unsafe"
)

func Recover() {
	if r := recover(); r != nil {
		log.Printf("panic: %v\nstack: %s", r, debug.Stack())
	}
}

func RecoverWithFunc(f func(r any)) {
	if r := recover(); r != nil {
		f(r)
		log.Printf("panic: %v\nstack: %s", r, debug.Stack())
	}
}

func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func JoinErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	messages := make([]string, len(errs))
	for i, err := range errs {
		messages[i] = err.Error()
	}
	return errors.New(strings.Join(messages, "; "))
}
