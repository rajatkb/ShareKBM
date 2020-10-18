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

func (mouseptr *MousePointer) adjustClickPosition(maxX int16, maxY int16) {
	mouseptr.Curx = (mouseptr.Curx / mouseptr.MaxX) * maxX
	mouseptr.Cury = (mouseptr.Cury / mouseptr.MaxY) * maxY
	mouseptr.Prevx = (mouseptr.Prevx / mouseptr.MaxX) * maxX
	mouseptr.Prevy = (mouseptr.Prevy / mouseptr.MaxY) * maxY
	mouseptr.MaxX = maxX
	mouseptr.MaxY = maxY
}

func (mouseptr *MousePointer) eventType() EventType {
	return MouseMove
}

//MouseButton used for denoting mouse button
type MouseButton int8

//Left , Right  , Center
const (
	LeftMouse       MouseButton = 1
	RightMouse      MouseButton = 2
	CenterMouse     MouseButton = 3
	ScrollUpMouse   MouseButton = 4
	ScrollDownMouse MouseButton = 5
)

//MouseClick ... struct for representing mouse click
type MouseClick struct {
	MouseBtn MouseButton
	Press    bool
}

func (mc *MouseClick) eventType() EventType {
	return MouseInteract
}

//CreateMousePointer creates a MousePointer data strcuture
func CreateMousePointer(curx int16, cury int16, maxX int16, maxY int16, prevx int16, prevy int16) *MousePointer {
	mpt := new(MousePointer)
	mpt.Curx = curx
	mpt.Cury = cury
	mpt.MaxX = maxX
	mpt.MaxY = maxY
	mpt.Prevx = prevx
	mpt.Prevy = prevy

	return mpt
}

//CreateMouseClick creates a mouse click object
func CreateMouseClick(mouseButton MouseButton, press bool) *MouseClick {
	mpt := new(MouseClick)
	mpt.MouseBtn = mouseButton
	mpt.Press = press
	return mpt
}

// EncDecMousePointer closure for encoder , decoder
func EncDecMousePointer() (func(*MousePointer) (bytes.Buffer, error), func(*bytes.Buffer) (MousePointer, error)) {
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

//EncDecMouseClick ... encode decodes mouse clicks
func EncDecMouseClick() (func(*MouseClick) (bytes.Buffer, error), func(*bytes.Buffer) (MouseClick, error)) {
	var dataptr bytes.Buffer
	enc := gob.NewEncoder(&dataptr)
	dec := gob.NewDecoder(&dataptr)

	return func(mdata *MouseClick) (bytes.Buffer, error) {
			err := enc.Encode(*mdata)
			if err != nil {
				fmt.Print(err)
				return dataptr, errors.New("bad data cannot encode")
			}
			return dataptr, nil
		},
		func(mdata *bytes.Buffer) (MouseClick, error) {
			var data MouseClick
			dataptr = *mdata
			err := dec.Decode(&data)
			if err != nil {
				fmt.Print("Bad data cannot encode !!")
				return data, errors.New("bad data cannot decode")
			}
			return data, nil
		}
}
