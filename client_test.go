package go_ernie

import "testing"

func TestNewClient(t *testing.T) {
	client, err := NewClient("akk", "akkk")
	if err != nil {
		t.Error(err)
	}
	if client == nil {
		t.Error("client is nil")
	}
}
