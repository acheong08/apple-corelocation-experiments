package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func readPascalString(r io.Reader) (string, error) {
	length := make([]byte, 2)
	_, err := r.Read(length)
	if err != nil {
		return "", err
	}
	str := make([]byte, binary.BigEndian.Uint16(length))
	_, err = r.Read(str)
	return string(str), err
}

type ArpcWrapper struct {
	Version       string
	Locale        string
	AppIdentifier string
	OsVersion     string
	FunctionId    int
	Payload       []byte
}

func ParseArpcRequest(data []byte) (ArpcWrapper, error) {
	r := io.NewSectionReader(
		bytes.NewReader(data), 0, int64(len(data)),
	)
	versionBytes := make([]byte, 2)
	if _, err := r.Read(versionBytes); err != nil {
		return ArpcWrapper{}, err
	}
	version := binary.BigEndian.Uint16(versionBytes)

	locale, err := readPascalString(r)
	if err != nil {
		return ArpcWrapper{}, err
	}
	appIdentifier, err := readPascalString(r)
	if err != nil {
		return ArpcWrapper{}, err
	}
	osVersion, err := readPascalString(r)
	if err != nil {
		return ArpcWrapper{}, err
	}
	unknownBytes := make([]byte, 4)
	if _, err := r.Read(unknownBytes); err != nil {
		return ArpcWrapper{}, err
	}
	unknown := int(binary.BigEndian.Uint32(unknownBytes))

	payloadLenBytes := make([]byte, 4)
	if _, err := r.Read(payloadLenBytes); err != nil {
		return ArpcWrapper{}, err
	}
	payloadLen := int(binary.BigEndian.Uint32(payloadLenBytes))

	payload := make([]byte, payloadLen)
	if _, err := r.Read(payload); err != nil {
		return ArpcWrapper{}, err
	}

	return ArpcWrapper{
		Version:       fmt.Sprintf("%d", version),
		Locale:        locale,
		AppIdentifier: appIdentifier,
		OsVersion:     osVersion,
		FunctionId:    unknown,
		Payload:       payload,
	}, nil
}
