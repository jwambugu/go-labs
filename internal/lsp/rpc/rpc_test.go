package rpc_test

import (
	"go-labs/internal/lsp/rpc"
	"testing"
)

type EncodingExample struct {
	Method bool
}

func TestEncodeMessage(t *testing.T) {
	expected := "Content-Length: 15\r\n\r\n{\"Method\":true}"

	actual := rpc.EncodeMessage(EncodingExample{Method: true})
	if expected != actual {
		t.Fatalf("expected: %s, actual %s", expected, actual)
	}
}

func TestDecodeMessage(t *testing.T) {
	incomingMessage := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"

	method, content, err := rpc.DecodeMessage([]byte(incomingMessage))
	if err != nil {
		t.Fatal(err)
	}

	if contentLength := len(content); contentLength != 15 {
		t.Fatalf("expected 16, got %d", contentLength)
	}

	if method != "hi" {
		t.Fatalf("expected hi, got %s", method)
	}
}
