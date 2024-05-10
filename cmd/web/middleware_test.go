package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH MyHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		//Do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH MyHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		//Do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.handler, but is %T", v))
	}
}
