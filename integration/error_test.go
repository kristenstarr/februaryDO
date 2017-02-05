package integration

import (
	"testing"
)

// SPEC:
// If the server doesn't recognize the command or if there's any problem with the message sent
// by the client it should return `ERROR\n`.

//Tests that ERROR should be received for all invalid sent messages.
func TestErrorMessages(t *testing.T) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	defer client.Close()
	if (err != nil) {
		t.Error(err)
	}
	respCode, err := client.Send("BadMessage")
	if (err != nil || respCode != ERROR) {
		t.Error("ERROR should be returned for bad syntax")
	}
	respCode, err = client.Send("QUERY|")
	if (err != nil || respCode != ERROR) {
		t.Error("ERROR should be returned for missing library")
	}
	respCode, err = client.Send("QUERY|LIB")
	if (err != nil || respCode != ERROR) {
		t.Error("ERROR should be returned for missing third argument")
	}
	respCode, err = client.Send("BAD|lib|")
	if (err != nil || respCode != ERROR) {
		t.Error("ERROR should be returned for bad verb choice")
	}
	respCode, err = client.Send("INDEX|&badchar|")
	if (err != nil || respCode != ERROR) {
		t.Error("ERROR should be returned for bad character")
	}
	teardownTest()
}

// Takes one sample Error message and runs a benchmark on it.
func BenchmarkErrorMessages(b *testing.B) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	defer client.Close()
	if (err != nil) {
		b.Error(err)
	}
	// Reset timer after setup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Send("QUERY|testpackage1")
	}

	teardownTest()
}
