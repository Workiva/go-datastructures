package palm

import (
	"log"
	"sync/atomic"
)

func init() {
	log.Println(`LOG HATES THIS.`)
}

type actionBundles []*actionBundle

type actionBundle struct {
	key    Key
	index  uint64
	action action
	node   *node
}

type actions []action

type action interface {
	operation() operation
	getKey() (Key, uint64) // returns nil if operation complete
	addResult(index uint64, result Key)
	len() uint64
}

type insertAction struct {
	keys        Keys
	count, done uint64
	completer   chan Keys
}

func (ia *insertAction) complete() {
	ia.completer <- ia.keys
	close(ia.completer)
}

func (ia *insertAction) operation() operation {
	return add
}

func (ia *insertAction) getKey() (Key, uint64) {
	index := atomic.AddUint64(&ia.count, 1)
	index-- // 0-index
	if index >= uint64(len(ia.keys)) {
		return nil, 0
	}

	return ia.keys[index], index
}

func (ia *insertAction) addResult(index uint64, result Key) {
	i := atomic.AddUint64(&ia.done, 1)
	i--
	if i >= uint64(len(ia.keys)) {
		return
	}
	ia.keys[index] = result
	if i == uint64(len(ia.keys))-1 {
		ia.complete()
	}
}

func (ia *insertAction) len() uint64 {
	return uint64(len(ia.keys))
}

func newInsertAction(keys Keys) *insertAction {
	return &insertAction{
		keys:      keys,
		completer: make(chan Keys),
	}
}

type getAction struct {
	*insertAction
}

func (ga *getAction) operation() operation {
	return get
}

func newGetAction(keys Keys) *getAction {
	return &getAction{
		&insertAction{
			keys:      keys,
			completer: make(chan Keys),
		},
	}
}
