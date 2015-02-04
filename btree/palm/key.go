package palm

import "sort"

func (keys Keys) search(key Key) int {
	return sort.Search(len(keys), func(i int) bool {
		return keys[i].Compare(key) >= 0
	})
}

func (keys *Keys) insert(key Key) Key {
	i := keys.search(key)
	return keys.insertAt(key, i)
}

func (keys *Keys) insertAt(key Key, i int) Key {
	if i == len(*keys) {
		*keys = append(*keys, key)
		return nil
	}

	if (*keys)[i].Compare(key) == 0 { //overwrite case
		oldKey := (*keys)[i]
		(*keys)[i] = key
		return oldKey
	}

	*keys = append(*keys, nil)
	copy((*keys)[i+1:], (*keys)[i:])
	(*keys)[i] = key
	return nil
}

func (keys *Keys) split() (Key, Keys, Keys) {
	i := (len(*keys) / 2) - 1
	middle := (*keys)[i]

	left, right := keys.splitAt(i)
	return middle, left, right
}

func (keys *Keys) splitAt(i int) (Keys, Keys) {
	right := make(Keys, len(*keys)-i-1, cap(*keys))
	copy(right, (*keys)[i+1:])
	for j := i + 1; j < len(*keys); j++ {
		(*keys)[j] = nil
	}
	*keys = (*keys)[:i+1]

	return *keys, right
}

func (keys Keys) last() Key {
	return keys[len(keys)-1]
}

func (keys Keys) first() Key {
	return keys[0]
}

func (keys Keys) needsSplit() bool {
	return cap(keys) == len(keys)
}

func (keys Keys) reverse() Keys {
	reversed := make(Keys, len(keys))
	for i := len(keys) - 1; i >= 0; i-- {
		reversed[len(keys)-1-i] = keys[i]
	}

	return reversed
}

func chunkKeys(keys Keys, numParts int64) []Keys {
	parts := make([]Keys, numParts)
	for i := int64(0); i < numParts; i++ {
		parts[i] = keys[i*int64(len(keys))/numParts : (i+1)*int64(len(keys))/numParts]
	}
	return parts
}
