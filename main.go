package main

import (
	"net"
	"fmt"
	"bufio"
	"strings"
	"sync"
)

type Library struct{
	Dependencies map[string]bool
	Parents map[string]bool
}

func main() {
	var libs = make(map[string]Library)
	var mutex = &sync.Mutex{}
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting on 8080")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("error accepting connecting")
		}
		go handleConnection(conn, libs, mutex)
	}
}

func validateMessage(msg string) bool {
	pieces := strings.Split(msg, "|")
	if (len(pieces) != 3) {
		return false;
	}
	method := pieces[0]
	if (method != "REMOVE" && method != "INDEX" && method != "QUERY") {
		return false;
	}
	lib := pieces[1]
	if (len(lib) < 1) {
		return false;
	}
	return true;
}

func handleConnection(conn net.Conn, libs map[string]Library, mutex *sync.Mutex) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		if (!validateMessage(message)) {
			conn.Write([]byte("ERROR\n"))
		} else {
			pieces := strings.Split(message, "|")
			method := pieces[0]
			lib := pieces[1]
			dependencies := pieces[2][:len(pieces[2]) - 1]
			mutex.Lock()
			switch method {

			case "REMOVE":
				if library, ok := libs[lib]; ok {
					if len(library.Parents) == 0 {
						for key, _ := range library.Dependencies {
							if dependentLibrary, ok := libs[key]; ok {
								delete(dependentLibrary.Parents, lib)
							}
						}
						delete(libs, lib)
						conn.Write([]byte("OK\n"))
					} else {
						conn.Write([]byte("FAIL\n"))
					}
				} else {
					conn.Write([]byte("OK\n"))
				}
			case "INDEX":
				// WHAT happens if new dependencies aren't there, but already indexed with old list
				// of dependencies that DO exist - do we remove??
				var splitDeps []string
				if (len(dependencies) > 0) {
					splitDeps = strings.Split(dependencies, ",")
				} else {
					splitDeps = make([]string, 0)
				}
				missingDep := false
				for _, dep := range splitDeps {
					if _, ok := libs[dep]; !ok {
						missingDep = true
					}
				}
				if (missingDep) {
					conn.Write([]byte("FAIL\n"))
				} else {

					depsMap := make(map[string]bool)
					for _, dep := range splitDeps {
						depsMap[dep] = true;
					}

					if oldLibrary, ok := libs[lib]; ok {

						for key, _ := range oldLibrary.Dependencies {
							if dependentLibrary, ok := libs[key]; ok {
								delete(dependentLibrary.Parents, lib)
							}
						}

						libs[lib] = Library{
							Parents: make(map[string]bool),
							Dependencies: depsMap,
						}
						for _, dep := range splitDeps {
							if depLibrary, ok := libs[dep]; ok {
								depLibrary.Parents[lib] = true
							}
						}

					} else {
						libs[lib] = Library{
							Parents: make(map[string]bool),
							Dependencies: depsMap,
						}
						for _, dep := range splitDeps {
							if depLibrary, ok := libs[dep]; ok {
								depLibrary.Parents[lib] = true
							}
						}
					}
					conn.Write([]byte("OK\n"))
				}


			case "QUERY":
				if _, ok := libs[lib]; ok {
					conn.Write([]byte("OK\n"))
				} else {
					conn.Write([]byte("FAIL\n"))
				}
			}
			mutex.Unlock()
		}
	}
}
