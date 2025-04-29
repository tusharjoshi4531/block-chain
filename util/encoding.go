package util

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

func EncoderGobEncodables(w io.Writer, vals ...any) error {
	enc := gob.NewEncoder(w)
	for _, val := range vals {
		if err := enc.Encode(val); err != nil {
			return err
		}
	}
	return nil
}

func DecodeGobDecodable(r io.Reader, vals ...any) error {
	dec := gob.NewDecoder(r)
	for _, val := range vals {
		if err := dec.Decode(val); err != nil {
			return err
		}
	}
	return nil
}

type Encoder interface {
	Encode(io.Writer) error
}

type Decoder interface {
	Decode(io.Reader) error
}

func ToEncoderSlice[T Encoder](items []T) []Encoder {
	result := make([]Encoder, len(items))
	for i, v := range items {
		result[i] = v
	}
	return result
}

func EncodeToBytesUsingEncoder(encode func(io.Writer) error) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := encode(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func EncodeToBytes(item Encoder) ([]byte, error) {
	return EncodeToBytesUsingEncoder(item.Encode)
}

func EncodeSlice(w io.Writer, items []Encoder) error {
	len := len(items)
	enc := gob.NewEncoder(w)
	if err := enc.Encode(len); err != nil {
		return err
	}

	for _, item := range items {
		if err := item.Encode(w); err != nil {
			return err
		}
	}

	return nil
}

func EncodeSliceToBytes(items []Encoder) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := EncodeSlice(buf, items); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeSlice[T Decoder](r io.Reader, factory func() T) ([]T, error) {
	var length int
	if err := gob.NewDecoder(r).Decode(&length); err != nil {
		return nil, err
	}
	fmt.Println(length)

	items := make([]T, 0, length)
	for i := 0; i < length; i++ {
		item := factory()
		if err := item.Decode(r); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
