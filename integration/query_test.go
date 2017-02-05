package integration

import (
	"testing"
)

// SPEC:
// For `QUERY` commands, the server returns `OK\n` if the package is indexed.
// It returns `FAIL\n` if the package isn't indexed.

//Tests package query should return FAIL when package has not been indexed.
func TestQueryNotIndexed(t *testing.T) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	defer client.Close()
	if (err != nil) {
		t.Error(err)
	}
	respCode, err := client.Send("QUERY|testpackage1|")
	if (err != nil || respCode != FAIL) {
		t.Error("Package query should FAIL if package hasn't been indexed")
	}
	teardownTest()
}

//Tests package removal should return OK when package was not indexed.
func TestQueryIndexed(t *testing.T) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	defer client.Close()
	if (err != nil) {
		t.Error(err)
	}
	client.Send("INDEX|testpackage1|")
	respCode, err := client.Send("QUERY|testpackage1|")
	if (err != nil || respCode != OK) {
		t.Error("Package query should return OK if package is indexed")
	}
	teardownTest()
}

// Takes one sample Query message and runs a benchmark on it.
func BenchmarkQueryMessage(b *testing.B) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	defer client.Close()
	if (err != nil) {
		b.Error(err)
	}
	// Reset timer after setup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Send("QUERY|testpackage1|")
	}

	teardownTest()
}
