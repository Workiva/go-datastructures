package common

type Comparator interface {
	Compare(Comparator) int
}

type Comparators []Comparator
