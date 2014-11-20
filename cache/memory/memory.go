package memory

import (
	"time"

	"github.com/dtynn/dout/cache"
)

const (
	actionSet int = iota
	actionSetex
	actionGet
	actionDel
)

var nonDead int64 = -1
var nonData = &data{}

type action struct {
	action int
	key    interface{}
	data   *data
	ch     chan *data
}

type data struct {
	val  interface{}
	dead int64
}

func isNonData(d *data) bool {
	return d == nonData
}

type safeMap struct {
	m       map[interface{}]*data
	actions chan *action
}

func New() *safeMap {
	s := &safeMap{
		m:       map[interface{}]*data{},
		actions: make(chan *action),
	}
	go func(s *safeMap) {
		s.run()
	}(s)
	return s
}

func (this *safeMap) Set(key, val interface{}) error {
	a := &action{
		action: actionSet,
		key:    key,
		data: &data{
			val:  val,
			dead: nonDead,
		},
	}
	this.actions <- a
	return nil
}

func (this *safeMap) Setex(key, val interface{}, expire int64) error {
	dead := time.Now().Add(time.Duration(expire) * time.Second).Unix()
	a := &action{
		action: actionSetex,
		key:    key,
		data: &data{
			val:  val,
			dead: dead,
		},
	}
	this.actions <- a
	return nil
}

func (this *safeMap) Get(key interface{}) (interface{}, error) {
	ch := make(chan *data)
	a := &action{
		action: actionGet,
		key:    key,
		ch:     ch,
	}
	this.actions <- a

	d := <-ch

	if isNonData(d) {
		return nil, cache.CacheNotFound
	}

	now := time.Now().Unix()
	if d.dead != nonDead && d.dead < now {
		return nil, cache.CacheNotFound
	}

	return d.val, nil
}

func (this *safeMap) Del(key interface{}) error {
	a := &action{
		action: actionDel,
		key:    key,
	}
	this.actions <- a
	return nil
}

func (this *safeMap) run() {
	for a := range this.actions {
		switch a.action {
		case actionSet, actionSetex:
			this.m[a.key] = a.data
		case actionGet:
			if d, has := this.m[a.key]; has {
				a.ch <- d
			} else {
				a.ch <- nonData
			}
		case actionDel:
			delete(this.m, a.key)
		}
	}
}
