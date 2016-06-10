package btree

type Item struct {
	Value   interface{}
	Payload []byte
}

type items []*Item

func (its items) split(numParts int) []items {
	parts := make([]items, numParts)
	for i := int64(0); i < int64(numParts); i++ {
		parts[i] = its[i*int64(len(its))/int64(numParts) : (i+1)*int64(len(its))/int64(numParts)]
	}
	return parts
}
