package integration

import (
	"testing"
)

// SPEC :
// For `INDEX` commands, the server returns `OK\n` if the package can be indexed.
// It returns `FAIL\n` if the package cannot be indexed because some of its dependencies aren't
// indexed yet and need to be installed first. If a package already exists, then its list
// of dependencies is updated to the one provided with the latest command.

//Tests indexing a Package with no dependencies, and then ensuring that Query for this
//Package returns true.
func TestSimpleIndexAndQuery(t *testing.T) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	if (err != nil) {
		t.SkipNow()
	}
	defer client.Close()
	respCode, err := client.Send("INDEX|testpackage1|")
	if (err != nil) {
		t.Error(err)
	}
	if (respCode != OK) {
		t.Error("Package should be indexed when no dependencies are required.")
	}

	respCode, err = client.Send("QUERY|testpackage1|")
	if (err != nil || respCode != OK) {
		t.Error("Once package has been indexed, it should be query-able")
	}
	teardownTest()
}

//Tests indexing a Package whose dependencies have not been indexed, ensuring that index
//fails and query also fails.
func TestIndexWithMissingDependencies(t *testing.T) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	if (err != nil) {
		t.SkipNow()
	}
	defer client.Close()
	respCode, err := client.Send("INDEX|testpackage1|missingdep")
	if (err != nil) {
		t.Error(err)
	}
	if (respCode != FAIL) {
		t.Error("Package should not be indexed when dependencies are missing.")
	}

	respCode, err = client.Send("QUERY|testpackage1|")
	if (err != nil || respCode != FAIL) {
		t.Error("Since package was not indexed, it should not be query-able")
	}
	teardownTest()
}

// Tests that indexing a package fails when dependencies are missing, and then succeeds once
// dependency packages have been indexed.
func TestIndexWithPresentDependencies(t *testing.T) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	if (err != nil) {
		t.SkipNow()
	}
	defer client.Close()

	//ensure that package cannot be indexed when dependencies are missing
	respCode, err := client.Send("INDEX|testpackage3|testpackage1,testpackage2")
	if (err != nil || respCode != FAIL) {
		t.Error("Initial package addition should fail")
	}

	//index dependencies
	respCode, err = client.Send("INDEX|testpackage1|")
	if (err != nil || respCode != OK) {
		t.Error("Initial package addition should succeed")
	}

	respCode, err = client.Send("INDEX|testpackage2|")
	if (err != nil || respCode != OK) {
		t.Error("Initial package addition should succeed")
	}

	//now our package index should succeed, since dependencies have been indexed
	respCode, err = client.Send("INDEX|testpackage3|testpackage1,testpackage2")
	if (err != nil || respCode != OK) {
		t.Error("Package addition should succeed once dependencies are indexed")
	}
	teardownTest()
}

// Tests that indexing a package that has already been indexed updates dependencies -
// we test this be attempting to remove package that are no longer dependencies.
func TestIndexUpdatesDependencies(t *testing.T) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	if (err != nil) {
		t.SkipNow()
	}
	defer client.Close()

	//index dependencies
	respCode, err := client.Send("INDEX|testpackage1|")
	if (err != nil || respCode != OK) {
		t.Error("Initial package addition should succeed")
	}

	respCode, err = client.Send("INDEX|testpackage2|")
	if (err != nil || respCode != OK) {
		t.Error("Initial package addition should succeed")
	}

	//index our package depending on testpackage1 and testpackage2
	respCode, err = client.Send("INDEX|testpackage3|testpackage1,testpackage2")
	if (err != nil || respCode != OK) {
		t.Error("Package addition should succeed once dependencies are indexed")
	}

	//attempt to remove testpackage1 should fail as testpackage3 depends on it.
	respCode, err = client.Send("REMOVE|testpackage1|")
	if (err != nil || respCode != FAIL) {
		t.Error("Dependency removal should fail")
	}

	//re-index our package depending on only testpackage2
	respCode, err = client.Send("INDEX|testpackage3|testpackage2")
	if (err != nil || respCode != OK) {
		t.Error("Package addition should succeed once dependencies are indexed")
	}

	//now removal of testpackage1 should succeed, as testpackage3 no longer depends on it.
	respCode, err = client.Send("REMOVE|testpackage1|")
	if (err != nil || respCode != OK) {
		t.Error("Dependency removal should succeed now since nothing depends on it")
	}
	teardownTest()
}

// Takes one sample Index message and runs a benchmark on it.
func BenchmarkIndexMessage(b *testing.B) {
	setupTest()
	client, err := MakeTCPPackageIndexClient(8080)
	if (err != nil) {
		b.SkipNow()
	}
	defer client.Close()

	// Reset timer after setup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Send("INDEX|testpackage1|")
	}

	teardownTest()
}
