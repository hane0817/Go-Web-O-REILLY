package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	}
	tracer.Trace("こんにちは，traceパッケージ")
	if buf.String() != "こんにちは，traceパッケージ\n" {
		t.Errorf("'%s'という誤った文字列が出力されました", buf.String()) //もともとはt.Errorだが後ろにfをつけることで書式指定できる
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("データ")
}
