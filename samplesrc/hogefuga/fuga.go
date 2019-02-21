package hogefuga

// このソースは、カバレッジ測定の対象外にしたいという設定。

// Answer 生命、宇宙、そして万物についての究極の疑問の答えに対応する文字列
func (t *Type) Answer() string {
	return t.foo[42]
}

// TheOne the one に対応する文字列
func (t *Type) TheOne() string {
	return t.foo[1]
}
