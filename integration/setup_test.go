package integration

// setupTest removes test packages 1-4 to ensure that our testing environment is clean.
func setupTest() {
	client, _ := MakeTCPPackageIndexClient(8080)
	client.Send("REMOVE|testpackage1|")
	client.Send("REMOVE|testpackage2|")
	client.Send("REMOVE|testpackage3|")
	client.Close()
}

// teardownTest removes test packages 1-4 to ensure that our testing environment is clean.
func teardownTest() {
	client, _ := MakeTCPPackageIndexClient(8080)
	client.Send("REMOVE|testpackage3|")
	client.Send("REMOVE|testpackage2|")
	client.Send("REMOVE|testpackage1|")
	client.Close()
}
