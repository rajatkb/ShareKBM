package device

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)

// MousePointer ... struct for data interchange format
type MousePointer struct {
	Curx  int16
	Cury  int16
	MaxX  int16
	MaxY  int16
	Prevx int16
	Prevy int16
}

//CreateMouse creates a MousePointer data strcuture
func CreateMouse(curx int16, cury int16, maxX int16, maxY int16, prevx int16, prevy int16) *MousePointer {
	mpt := new(MousePointer)
	mpt.Curx = curx
	mpt.Cury = cury
	mpt.MaxX = maxX
	mpt.MaxY = maxY
	mpt.Prevx = prevx
	mpt.Prevy = prevy
	return mpt
}

// EncDecMouseptr closure for encoder , decoder
func EncDecMouseptr() (func(*MousePointer) (bytes.Buffer, error), func(*bytes.Buffer) (MousePointer, error)) {
	var dataptr bytes.Buffer
	enc := gob.NewEncoder(&dataptr)
	dec := gob.NewDecoder(&dataptr)

	return func(mdata *MousePointer) (bytes.Buffer, error) {
			err := enc.Encode(*mdata)
			if err != nil {
				fmt.Print(err)
				return dataptr, errors.New("bad data cannot encode")
			}
			return dataptr, nil
		},
		func(mdata *bytes.Buffer) (MousePointer, error) {
			var data MousePointer
			dataptr = *mdata
			err := dec.Decode(&data)
			if err != nil {
				fmt.Print("Bad data cannot encode !!")
				return data, errors.New("bad data cannot decode")
			}
			return data, nil
		}
}
