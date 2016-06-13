package btree

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z ID) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBytes(o, []byte(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ID) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var tmp []byte
		tmp, bts, err = msgp.ReadBytesBytes(bts, []byte((*z)))
		(*z) = ID(tmp)
	}
	if err != nil {
		return
	}
	o = bts
	return
}

func (z ID) Msgsize() (s int) {
	s = msgp.BytesPrefixSize + len([]byte(z))
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Key) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "u"
	o = append(o, 0x83, 0xa1, 0x75)
	o = msgp.AppendBytes(o, []byte(z.UUID))
	// string "v"
	o = append(o, 0xa1, 0x76)
	o, err = msgp.AppendIntf(o, z.Value)
	if err != nil {
		return
	}
	// string "p"
	o = append(o, 0xa1, 0x70)
	o = msgp.AppendBytes(o, z.Payload)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Key) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			{
				var tmp []byte
				tmp, bts, err = msgp.ReadBytesBytes(bts, []byte(z.UUID))
				z.UUID = ID(tmp)
			}
			if err != nil {
				return
			}
		case "v":
			z.Value, bts, err = msgp.ReadIntfBytes(bts)
			if err != nil {
				return
			}
		case "p":
			z.Payload, bts, err = msgp.ReadBytesBytes(bts, z.Payload)
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

func (z *Key) Msgsize() (s int) {
	s = 1 + 2 + msgp.BytesPrefixSize + len([]byte(z.UUID)) + 2 + msgp.GuessSize(z.Value) + 2 + msgp.BytesPrefixSize + len(z.Payload)
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Keys) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for xvk := range z {
		if z[xvk] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z[xvk].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Keys) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var xsz uint32
	xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(xsz) {
		(*z) = (*z)[:xsz]
	} else {
		(*z) = make(Keys, xsz)
	}
	for bzg := range *z {
		if msgp.IsNil(bts) {
			bts, err = msgp.ReadNilBytes(bts)
			if err != nil {
				return
			}
			(*z)[bzg] = nil
		} else {
			if (*z)[bzg] == nil {
				(*z)[bzg] = new(Key)
			}
			bts, err = (*z)[bzg].UnmarshalMsg(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z Keys) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for bai := range z {
		if z[bai] == nil {
			s += msgp.NilSize
		} else {
			s += z[bai].Msgsize()
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Node) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "id"
	o = append(o, 0x84, 0xa2, 0x69, 0x64)
	o = msgp.AppendBytes(o, []byte(z.ID))
	// string "il"
	o = append(o, 0xa2, 0x69, 0x6c)
	o = msgp.AppendBool(o, z.IsLeaf)
	// string "cv"
	o = append(o, 0xa2, 0x63, 0x76)
	o = msgp.AppendArrayHeader(o, uint32(len(z.ChildValues)))
	for cmr := range z.ChildValues {
		o, err = msgp.AppendIntf(o, z.ChildValues[cmr])
		if err != nil {
			return
		}
	}
	// string "ck"
	o = append(o, 0xa2, 0x63, 0x6b)
	o = msgp.AppendArrayHeader(o, uint32(len(z.ChildKeys)))
	for ajw := range z.ChildKeys {
		if z.ChildKeys[ajw] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.ChildKeys[ajw].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Node) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "id":
			{
				var tmp []byte
				tmp, bts, err = msgp.ReadBytesBytes(bts, []byte(z.ID))
				z.ID = ID(tmp)
			}
			if err != nil {
				return
			}
		case "il":
			z.IsLeaf, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				return
			}
		case "cv":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.ChildValues) >= int(xsz) {
				z.ChildValues = z.ChildValues[:xsz]
			} else {
				z.ChildValues = make([]interface{}, xsz)
			}
			for cmr := range z.ChildValues {
				z.ChildValues[cmr], bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
			}
		case "ck":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.ChildKeys) >= int(xsz) {
				z.ChildKeys = z.ChildKeys[:xsz]
			} else {
				z.ChildKeys = make(Keys, xsz)
			}
			for ajw := range z.ChildKeys {
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					z.ChildKeys[ajw] = nil
				} else {
					if z.ChildKeys[ajw] == nil {
						z.ChildKeys[ajw] = new(Key)
					}
					bts, err = z.ChildKeys[ajw].UnmarshalMsg(bts)
					if err != nil {
						return
					}
				}
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

func (z *Node) Msgsize() (s int) {
	s = 1 + 3 + msgp.BytesPrefixSize + len([]byte(z.ID)) + 3 + msgp.BoolSize + 3 + msgp.ArrayHeaderSize
	for cmr := range z.ChildValues {
		s += msgp.GuessSize(z.ChildValues[cmr])
	}
	s += 3 + msgp.ArrayHeaderSize
	for ajw := range z.ChildKeys {
		if z.ChildKeys[ajw] == nil {
			s += msgp.NilSize
		} else {
			s += z.ChildKeys[ajw].Msgsize()
		}
	}
	return
}
