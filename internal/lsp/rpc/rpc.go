package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
)

type BaseMessage struct {
	Method string `json:"method,omitempty"`
}

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		log.Panicf("encode msg: %v\n", err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("separator not found")
	}

	contentLengthBytes := header[len("Content-Length: "):]

	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, fmt.Errorf("content length: %v", err)
	}

	var baseMsg BaseMessage
	if err = json.Unmarshal(content[:contentLength], &baseMsg); err != nil {
		return "", nil, fmt.Errorf("decode message: %v", err)
	}

	return baseMsg.Method, content[:contentLength], nil
}
