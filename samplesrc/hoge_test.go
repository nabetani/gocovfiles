package samplesrc

import (
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	o := New()
	type kvType struct {
		k int
		v string
	}
	kvs := []kvType{}
	for k, v := range o.foo {
		kvs = append(kvs, kvType{k: k, v: v})
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].k < kvs[j].k })
	if kvs[0].k != 0 || kvs[0].v != "hoge" {
		t.Errorf("kvs[0] = %v, want {0, hoge}", kvs[0])
	}
	if kvs[1].k != 42 || kvs[1].v != "answer" {
		t.Errorf("kvs[1] = %v, want {42, answer}", kvs[1])
	}
}

func TestSumValue(t *testing.T) {
	o := New()
	if "hogeanswer" != o.SumValue() {
		t.Errorf(`o.SumValue()=%q, want "hogeanswer"`, o.SumValue())
	}
}
