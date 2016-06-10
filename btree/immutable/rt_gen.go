package btree

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z *Tr) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "u"
	o = append(o, 0x84, 0xa1, 0x75)
	o, err = z.UUID.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "c"
	o = append(o, 0xa1, 0x63)
	o = msgp.AppendInt(o, z.Count)
	// string "r"
	o = append(o, 0xa1, 0x72)
	o, err = z.Root.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "nw"
	o = append(o, 0xa2, 0x6e, 0x77)
	o = msgp.AppendInt(o, z.NodeWidth)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Tr) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "u":
			bts, err = z.UUID.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "c":
			z.Count, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "r":
			bts, err = z.Root.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "nw":
			z.NodeWidth, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *Tr) Msgsize() (s int) {
	s = 1 + 2 + z.UUID.Msgsize() + 2 + msgp.IntSize + 2 + z.Root.Msgsize() + 3 + msgp.IntSize
	return
}
