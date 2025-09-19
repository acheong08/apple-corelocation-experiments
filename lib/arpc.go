package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

func readPascalString(r io.Reader, remaining *int64) (string, error) {
	lengthBytes := make([]byte, 2)
	if _, err := r.Read(lengthBytes); err != nil {
		return "", err
	}
	*remaining -= 2
	length := binary.BigEndian.Uint16(lengthBytes)
	if int64(length) > *remaining {
		return "", fmt.Errorf("pascal string length %d exceeds remaining buffer size %d", length, *remaining)
	}
	str := make([]byte, length)
	if _, err := r.Read(str); err != nil {
		return "", err
	}
	*remaining -= int64(length)
	return string(str), nil
}

type ArpcRequest struct {
	Version       string
	Locale        string
	AppIdentifier string
	OsVersion     string
	FunctionId    int
	Payload       []byte
}

func (a *ArpcRequest) Deserialize(data []byte) error {
	r := bytes.NewReader(data)
	remaining := int64(len(data))

	versionBytes := make([]byte, 2)
	if _, err := r.Read(versionBytes); err != nil {
		return err
	}
	remaining -= 2
	version := binary.BigEndian.Uint16(versionBytes)

	locale, err := readPascalString(r, &remaining)
	if err != nil {
		return err
	}
	appIdentifier, err := readPascalString(r, &remaining)
	if err != nil {
		return err
	}
	osVersion, err := readPascalString(r, &remaining)
	if err != nil {
		return err
	}
	unknownBytes := make([]byte, 4)
	if _, err := r.Read(unknownBytes); err != nil {
		return err
	}
	remaining -= 4
	unknown := int(binary.BigEndian.Uint32(unknownBytes))

	payloadLenBytes := make([]byte, 4)
	if _, err := r.Read(payloadLenBytes); err != nil {
		return err
	}
	remaining -= 4
	payloadLen := int(binary.BigEndian.Uint32(payloadLenBytes))

	if int64(payloadLen) > remaining {
		return fmt.Errorf("payload length %d exceeds remaining buffer size %d", payloadLen, remaining)
	}
	payload := make([]byte, payloadLen)
	if _, err := r.Read(payload); err != nil {
		return err
	}

	*a = ArpcRequest{
		Version:       fmt.Sprintf("%d", version),
		Locale:        locale,
		AppIdentifier: appIdentifier,
		OsVersion:     osVersion,
		FunctionId:    unknown,
		Payload:       payload,
	}
	return nil
}

func writePascalString(w io.Writer, s string) error {
	length := uint16(len(s))
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}

func (a *ArpcRequest) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	version, err := strconv.Atoi(a.Version)
	if err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, uint16(version)); err != nil {
		return nil, err
	}

	if err := writePascalString(buf, a.Locale); err != nil {
		return nil, err
	}
	if err := writePascalString(buf, a.AppIdentifier); err != nil {
		return nil, err
	}
	if err := writePascalString(buf, a.OsVersion); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(a.FunctionId)); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(len(a.Payload))); err != nil {
		return nil, err
	}

	if _, err := buf.Write(a.Payload); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
