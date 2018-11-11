package goredis

import (
	"github.com/alicebob/miniredis"
	"testing"
)

func TestGetProto(t *testing.T) {
	addr := "127.0.0.1:6379"

	proto := getProto(addr)

	t.Logf("addr proto: %s", proto)
}

func TestNewClient(t *testing.T) {
	addr := "127.0.0.1:6379"
	password := "123"

	client := NewClient(addr, password)

	if client.addr == addr && client.password == password {
		t.Logf("创建client成功: addr: %s, password: %s", client.addr, client.password)
	} else {
		t.Error("创建client失败")
	}
}

func TestGet(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	addr := s.Addr()

	client := NewClient(addr, "")

	conn, err := client.Get()
	defer conn.Close()

	if err != nil {
		t.Errorf("error: %s\n", err)
	} else {

		if conn == nil {
			t.Error("error, cannot connect redis server")
		} else {
			t.Logf("连接成功conn: %s", conn)
		}
	}
}

func TestDo(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	addr := s.Addr()

	client := NewClient(addr, "")

	ret, err := client.Do("PING", nil)
	if err != nil {
		t.Errorf("error: %s\n", err)
	} else {
		t.Logf("PING返回: %s", ret)
	}

	ret, err = client.Do("SET", "msg", "hello")
	if err != nil {
		t.Errorf("error: %s\n", err)
	} else {
		t.Logf("SET 'msg' 返回: %s", ret)
	}

	ret, err = client.Do("GET", "msg")
	if err != nil {
		t.Errorf("error: %s\n", err)
	} else {
		t.Logf("GET 'msg' 返回: %s", ret)
	}

}
