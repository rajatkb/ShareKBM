package device

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)

// KeyBoard ... struct for data interchange format
type KeyBoard struct {
	Keys []int8
}

//EncDecKeyboard ... function for encoding decoding keyboard strokes
func EncDecKeyboard() (func(*KeyBoard) (bytes.Buffer, error), func(*bytes.Buffer) (KeyBoard, error)) {
	var dataptr bytes.Buffer
	enc := gob.NewEncoder(&dataptr)
	dec := gob.NewDecoder(&dataptr)

	return func(mdata *KeyBoard) (bytes.Buffer, error) {
			err := enc.Encode(*mdata)
			if err != nil {
				fmt.Print(err)
				return dataptr, errors.New("bad data cannot encode")
			}
			return dataptr, nil
		},
		func(mdata *bytes.Buffer) (KeyBoard, error) {
			var data KeyBoard
			dataptr = *mdata
			err := dec.Decode(&data)
			if err != nil {
				fmt.Print("Bad data cannot encode !!")
				return data, errors.New("bad data cannot decode")
			}
			return data, nil
		}
}
