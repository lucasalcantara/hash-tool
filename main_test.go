package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSuccessRequests(t *testing.T) {
	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "server 1")
	}))
	defer server.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "server 2")
	}))

	defer server2.Close()

	// when
	hashes := getHashResponses([]string{server.URL, server2.URL}, 2)

	// then
	expectedServeHash := "a664d1c39f85cc989ba372f42f77104b"
	if hashes[server.URL] != expectedServeHash {
		t.Logf("hash does not match. expected %s actual %s for url %s", expectedServeHash, hashes[server.URL], server.URL)
		t.Fail()
	}

	expectedServe2Hash := "b275ddd2ece1a8d3a38d296f506083a3"
	if hashes[server2.URL] != expectedServe2Hash {
		t.Logf("hash does not match. expected %s actual %s for url %s", expectedServe2Hash, hashes[server2.URL], server2.URL)
		t.Fail()
	}
}

func TestSuccessRequest_WhenUrlDoesNotHaveHttpPrefix(t *testing.T) {
	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "server 1")
	}))
	defer server.Close()

	// when
	hashes := getHashResponses([]string{strings.Replace(server.URL, "http://", "", -1)}, 2)

	// then
	expectedServeHash := "a664d1c39f85cc989ba372f42f77104b"
	if hashes[server.URL] != expectedServeHash {
		t.Logf("hash does not match. expected %s actual %s for url %s", expectedServeHash, hashes[server.URL], server.URL)
		t.Fail()
	}
}

func TestNoDuplicationRequests(t *testing.T) {
	// given
	countRequest := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		countRequest++
	}))
	defer server.Close()

	// when
	getHashResponses([]string{server.URL, server.URL}, 2)

	// then
	if countRequest != 1 {
		t.Logf("duplicate requests for the same url")
		t.Fail()
	}
}

func TestWhenRequestFails(t *testing.T) {
	// given
	url := "someWeirdUrl"
	// when
	hashes := getHashResponses([]string{url}, 2)

	// then
	expectedServeHash := ""
	if hashes[url] != expectedServeHash {
		t.Logf("hash does not match. expected %s actual %s for url %s", expectedServeHash, hashes[url], url)
		t.Fail()
	}
}
