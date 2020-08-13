package test

import (
	"testing"

	device "sharekbm/device"
)

func TestEncoderDecoder(t *testing.T) {
	pos := device.CreateMouse(0, 0, 100, 100, 10, 10)
	enc, dec := device.EncDecMouseptr()
	b, error1 := enc(pos)
	if error1 != nil {
		t.Errorf("Error parsing pos to binary error : %s", error1.Error())
	}
	posr, error2 := dec(&b)
	if error2 != nil {
		t.Errorf("Error parsing posr from binary error : %s", error2.Error())
	}
	if posr.MaxX != pos.MaxX {
		t.Errorf("Error parsing:  value %+v , decoder found %+v ", pos, posr)
	}

}
