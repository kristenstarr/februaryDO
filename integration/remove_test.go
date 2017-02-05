package integration

import (
	"testing"
)

// SPEC:
// For `REMOVE` commands, the server returns `OK\n` if the package could be removed from the index.
// It returns `FAIL\n` if the package could not be removed from the index because some other indexed
// package depends on it. It returns `OK\n` if the package wasn't indexed.

//Tests package removal should return OK when package was not indexed.
func TestRemovalNotPresent(t *testing.T) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	defer client.Close()
	if (err != nil) {
		t.Error(err)
	}
	respCode, err := client.Send("REMOVE|testpackage1|")
	if (err != nil || respCode != OK) {
		t.Error("Package remove should return true when package is not indexed")
	}
	teardownTest()
}

//Tests package removal should succeed when no other packages depend on it.
func TestRemovalSuccess(t *testing.T) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	defer client.Close()
	if (err != nil) {
		t.Error(err)
	}
	respCode, err := client.Send("INDEX|testpackage1|")
	if (err != nil || respCode != OK) {
		t.Error("Package should be successfully indexed")
	}
	respCode, err = client.Send("REMOVE|testpackage1|")
	if (err != nil || respCode != OK) {
		t.Error("Package remove should succeed when no other packages depend on it")
	}
	teardownTest()
}

//Tests package removal should fail when other packages depend on it.
func TestRemovalFailure(t *testing.T) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	defer client.Close()
	if (err != nil) {
		t.Error(err)
	}
	respCode, err := client.Send("INDEX|testpackage1|")
	if (err != nil || respCode != OK) {
		t.Error("Package should be successfully indexed")
	}
	respCode, err = client.Send("INDEX|testpackage2|testpackage1")
	if (err != nil || respCode != OK) {
		t.Error("Package should be successfully indexed since dependency is present")
	}
	respCode, err = client.Send("REMOVE|testpackage1|")
	if (err != nil || respCode != FAIL) {
		t.Error("Package remove should fail as another package depends on it")
	}
	teardownTest()
}

// Takes one sample Remove message and runs a benchmark on it.
func BenchmarkRemoveNonIndexedMessage(b *testing.B) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	defer client.Close()
	if (err != nil) {
		b.Error(err)
	}
	// Reset timer after setup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Send("REMOVE|testpackage1|")
	}

	teardownTest()
}
