package deezy

import "testing"

func TestIsChannelOpen(t *testing.T) {
	t.Setenv("DEEZY_PEER", "024bfaf0cabe7f874fd33ebf7c6f4e5385971fc504ef3f492432e9e3ec77e1b5cf")
	got := IsChannelOpen()
	want := true

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestIsNoChannelOpen(t *testing.T) {
	// need to add mock for ListChannels call to return []
	t.Setenv("DEEZY_PEER", "invalidpeer")
	got := IsChannelOpen()
	want := false

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}