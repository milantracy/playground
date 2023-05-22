package greetings

import "testing"

func TestHello(t *testing.T) {
	expected := "Hello, world."
	if msg := Hello("world"); msg != expected {
		t.Errorf("Hello() = %q, expect %q", msg, expected)
	}
}
