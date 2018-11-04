package goredis

import (
	"testing"
)

func TestParseInt(t *testing.T) {

	data1 := []byte("-ERR unknown command 'foobar'")
	data2 := []byte("-1")

	ret, err := parseInt(data1)
	if err != nil {
		t.Logf("error: %q\n", err)
	} else {
		t.Errorf("success: %d\n", ret)
	}

	ret, err = parseInt(data2)
	if err != nil {
		t.Errorf("error: %q\n", err)
	} else {
		t.Logf("success: %d\n", ret)
	}

}
