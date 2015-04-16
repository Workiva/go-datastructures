package palm

import "sync"

var keyBundlePool = sync.Pool{
	New: func() interface{} {
		return &keyBundle{}
	},
}
