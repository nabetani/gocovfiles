package samplesrc

import (
	"fmt"
	"sort"
	"strings"
)

// Type 特に意味はない構造体
type Type struct {
	foo map[int]string
}

// New Type を作る
func New() *Type {
	return &Type{foo: map[int]string{0: "hoge", 42: "answer"}}
}

// SumKey キーの合計を計算する。目的はない。
func (t *Type) SumKey() int {
	sum := 0
	for k := range t.foo {
		sum += k
	}
	return sum
}

// SumValue 値の文字列を連結する
func (t *Type) SumValue() string {
	keys := []int{}
	for k := range t.foo {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	b := strings.Builder{}
	for _, k := range keys {
		b.WriteString(t.foo[k])
	}
	return b.String()
}

// Delete key-value ペアをひとつ削除する
func (t *Type) Delete(key int) {
	delete(t.foo, key)
}

// Load key から value を引く
func (t *Type) Load(key int) (string, bool) {
	value, ok := t.foo[key]
	return value, ok
}

// MustLoad Load key から value を引く。キーがなかったら panic
func (t *Type) MustLoad(key int) string {
	value, ok := t.foo[key]
	if !ok {
		panic(fmt.Sprintf("no such key %v", key))
	}
	return value
}

// Store key-value ペアをひとつ保存する
func (t *Type) Store(key int, value string) {
	t.foo[key] = value
}
