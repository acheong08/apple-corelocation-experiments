package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
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

type ArpcRequest struct {
	Version       string
	Locale        string
	AppIdentifier string
	OsVersion     string
	FunctionId    int
	Payload       []byte
}

func (a *ArpcRequest) Deserialize(data []byte) error {
	r := io.NewSectionReader(
		bytes.NewReader(data), 0, int64(len(data)),
	)
	versionBytes := make([]byte, 2)
	if _, err := r.Read(versionBytes); err != nil {
		return err
	}
	version := binary.BigEndian.Uint16(versionBytes)

	locale, err := readPascalString(r)
	if err != nil {
		return err
	}
	appIdentifier, err := readPascalString(r)
	if err != nil {
		return err
	}
	osVersion, err := readPascalString(r)
	if err != nil {
		return err
	}
	unknownBytes := make([]byte, 4)
	if _, err := r.Read(unknownBytes); err != nil {
		return err
	}
	unknown := int(binary.BigEndian.Uint32(unknownBytes))

	payloadLenBytes := make([]byte, 4)
	if _, err := r.Read(payloadLenBytes); err != nil {
		return err
	}
	payloadLen := int(binary.BigEndian.Uint32(payloadLenBytes))

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
