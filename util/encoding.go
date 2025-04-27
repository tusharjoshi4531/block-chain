package util

import (
	"bytes"
	"encoding/gob"
	"io"
)

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

func DecodeSlice[T Decoder](r io.Reader) ([]T, error) {
	var length int
	if err := gob.NewDecoder(r).Decode(&length); err != nil {
		return nil, err
	}

	items := make([]T, 0, length)
	for i := 0; i < length; i++ {
		var item T
		if err := item.Decode(r); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
