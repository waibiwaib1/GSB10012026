package fasthttp

import "testing"

type bytePair [2]string

func TestArgsAll(t *testing.T) {
	t.Parallel()

	var args Args
	args.Add("foo", "bar")
	args.Add("baz", "qux")
	args.AddNoValue("empty")

	assertPairs(t, argsVisitAll(&args), argsAllPairs(&args))

	var got []bytePair
	for key, value := range args.All() {
		got = append(got, bytePair{string(key), string(value)})
		if len(got) == 2 {
			break
		}
	}
	if len(got) != 2 {
		t.Fatalf("unexpected pair count after break: %d", len(got))
	}
}

func TestHeaderAll(t *testing.T) {
	t.Parallel()

	var response ResponseHeader
	response.SetContentType("text/plain")
	response.Set("X-Custom", "value")
	if err := response.AddTrailer("X-Response-Trailer"); err != nil {
		t.Fatal(err)
	}
	var responseCookie Cookie
	responseCookie.SetKey("session")
	responseCookie.SetValue("abc")
	response.SetCookie(&responseCookie)

	assertPairs(t, responseVisitAllPairs(&response), responseAllPairs(&response))
	assertValues(t, responseTrailerVisitAll(&response), responseTrailerAll(&response))
	assertPairs(t, responseCookieVisitAll(&response), responseCookieAll(&response))

	var request RequestHeader
	request.SetMethod("POST")
	request.SetHost("example.com")
	request.SetContentType("application/json")
	request.Set("X-Custom", "value")
	request.SetCookie("session", "abc")
	if err := request.AddTrailer("X-Request-Trailer"); err != nil {
		t.Fatal(err)
	}

	assertPairs(t, requestVisitAllPairs(&request), requestAllPairs(&request))
	assertValues(t, requestTrailerVisitAll(&request), requestTrailerAll(&request))
	assertPairs(t, requestCookieVisitAll(&request), requestCookieAll(&request))
}

func TestHeaderAllBreak(t *testing.T) {
	t.Parallel()

	var response ResponseHeader
	response.SetContentType("text/plain")
	response.Set("X-Custom", "value")
	response.SetConnectionClose()

	var responsePairs []bytePair
	for key, value := range response.All() {
		responsePairs = append(responsePairs, bytePair{string(key), string(value)})
		if len(responsePairs) == 1 {
			break
		}
	}
	if len(responsePairs) != 1 {
		t.Fatalf("unexpected response pair count after break: %d", len(responsePairs))
	}

	var request RequestHeader
	request.SetHost("example.com")
	request.Set("X-Custom", "value")
	request.SetConnectionClose()

	var requestPairs []bytePair
	for key, value := range request.All() {
		requestPairs = append(requestPairs, bytePair{string(key), string(value)})
		if len(requestPairs) == 1 {
			break
		}
	}
	if len(requestPairs) != 1 {
		t.Fatalf("unexpected request pair count after break: %d", len(requestPairs))
	}
}

func assertPairs(t *testing.T, expected, got []bytePair) {
	t.Helper()
	if len(expected) != len(got) {
		t.Fatalf("unexpected pair count: got %d, want %d (%v vs %v)", len(got), len(expected), got, expected)
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("pair %d differs: got %q, want %q", i, got[i], expected[i])
		}
	}
}

func assertValues(t *testing.T, expected, got []string) {
	t.Helper()
	if len(expected) != len(got) {
		t.Fatalf("unexpected value count: got %d, want %d (%q vs %q)", len(got), len(expected), got, expected)
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("value %d differs: got %q, want %q", i, got[i], expected[i])
		}
	}
}

func argsVisitAll(args *Args) []bytePair {
	var pairs []bytePair
	args.VisitAll(func(key, value []byte) { pairs = append(pairs, bytePair{string(key), string(value)}) })
	return pairs
}

func argsAllPairs(args *Args) []bytePair {
	var pairs []bytePair
	for key, value := range args.All() {
		pairs = append(pairs, bytePair{string(key), string(value)})
	}
	return pairs
}

func responseVisitAllPairs(header *ResponseHeader) []bytePair {
	var pairs []bytePair
	header.VisitAll(func(key, value []byte) { pairs = append(pairs, bytePair{string(key), string(value)}) })
	return pairs
}

func responseAllPairs(header *ResponseHeader) []bytePair {
	var pairs []bytePair
	for key, value := range header.All() {
		pairs = append(pairs, bytePair{string(key), string(value)})
	}
	return pairs
}

func requestVisitAllPairs(header *RequestHeader) []bytePair {
	var pairs []bytePair
	header.VisitAll(func(key, value []byte) { pairs = append(pairs, bytePair{string(key), string(value)}) })
	return pairs
}

func requestAllPairs(header *RequestHeader) []bytePair {
	var pairs []bytePair
	for key, value := range header.All() {
		pairs = append(pairs, bytePair{string(key), string(value)})
	}
	return pairs
}

func responseTrailerVisitAll(header *ResponseHeader) []string {
	var values []string
	header.VisitAllTrailer(func(value []byte) { values = append(values, string(value)) })
	return values
}

func responseTrailerAll(header *ResponseHeader) []string {
	var values []string
	for value := range header.AllTrailer() {
		values = append(values, string(value))
	}
	return values
}

func requestTrailerVisitAll(header *RequestHeader) []string {
	var values []string
	header.VisitAllTrailer(func(value []byte) { values = append(values, string(value)) })
	return values
}

func requestTrailerAll(header *RequestHeader) []string {
	var values []string
	for value := range header.AllTrailer() {
		values = append(values, string(value))
	}
	return values
}

func responseCookieVisitAll(header *ResponseHeader) []bytePair {
	var pairs []bytePair
	header.VisitAllCookie(func(key, value []byte) { pairs = append(pairs, bytePair{string(key), string(value)}) })
	return pairs
}

func responseCookieAll(header *ResponseHeader) []bytePair {
	var pairs []bytePair
	for key, value := range header.AllCookie() {
		pairs = append(pairs, bytePair{string(key), string(value)})
	}
	return pairs
}

func requestCookieVisitAll(header *RequestHeader) []bytePair {
	var pairs []bytePair
	header.VisitAllCookie(func(key, value []byte) { pairs = append(pairs, bytePair{string(key), string(value)}) })
	return pairs
}

func requestCookieAll(header *RequestHeader) []bytePair {
	var pairs []bytePair
	for key, value := range header.AllCookie() {
		pairs = append(pairs, bytePair{string(key), string(value)})
	}
	return pairs
}
