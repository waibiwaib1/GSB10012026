package fasthttp

import (
	"bufio"
	"bytes"
	"iter"
	"reflect"
	"testing"
)

func collectSeq2(seq iter.Seq2[[]byte, []byte]) [][2]string {
	var got [][2]string
	for key, value := range seq {
		got = append(got, [2]string{string(key), string(value)})
	}
	return got
}

func collectSeq(seq iter.Seq[[]byte]) []string {
	var got []string
	for value := range seq {
		got = append(got, string(value))
	}
	return got
}

func findHeader(headers [][2]string, key string) (string, bool) {
	for _, header := range headers {
		if header[0] == key {
			return header[1], true
		}
	}
	return "", false
}

func TestArgsAll(t *testing.T) {
	t.Parallel()

	var a Args
	a.Set("foo", "bar")
	a.Set("baz", "qux")

	got := collectSeq2(a.All())
	want := [][2]string{{"foo", "bar"}, {"baz", "qux"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected args %v. Expected %v", got, want)
	}

	visited := 0
	for range a.All() {
		visited++
		break
	}
	if visited != 1 {
		t.Fatalf("unexpected number of args before break: %d. Expected 1", visited)
	}
}

func TestResponseHeaderAll(t *testing.T) {
	t.Parallel()

	var h ResponseHeader
	r := bytes.NewBufferString("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 123\r\nSet-Cookie: aa=bb; path=/foo/bar\r\nSet-Cookie: ccc\r\nTrailer: Foo, Bar\r\n\r\n")
	if err := h.Read(bufio.NewReader(r)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	headers := collectSeq2(h.All())
	if value, ok := findHeader(headers, HeaderContentLength); !ok || value != "123" {
		t.Fatalf("unexpected content-length %q, %v. Expected 123", value, ok)
	}
	if value, ok := findHeader(headers, HeaderContentType); !ok || value != "text/plain" {
		t.Fatalf("unexpected content-type %q, %v. Expected text/plain", value, ok)
	}
	if value, ok := findHeader(headers, HeaderTrailer); !ok || value != "Foo, Bar" {
		t.Fatalf("unexpected trailer %q, %v. Expected Foo, Bar", value, ok)
	}

	var cookies [][2]string
	for key, value := range h.All() {
		if string(key) == HeaderSetCookie {
			cookies = append(cookies, [2]string{string(key), string(value)})
		}
	}
	wantCookies := [][2]string{{HeaderSetCookie, "aa=bb; path=/foo/bar"}, {HeaderSetCookie, "ccc"}}
	if !reflect.DeepEqual(cookies, wantCookies) {
		t.Fatalf("unexpected cookies %v. Expected %v", cookies, wantCookies)
	}

	if got := collectSeq(h.AllTrailer()); !reflect.DeepEqual(got, []string{"Foo", "Bar"}) {
		t.Fatalf("unexpected trailers %v. Expected [Foo Bar]", got)
	}
	if got := collectSeq2(h.AllCookie()); !reflect.DeepEqual(got, [][2]string{{"aa", "aa=bb; path=/foo/bar"}, {"ccc", "ccc"}}) {
		t.Fatalf("unexpected cookies %v", got)
	}
}

func TestRequestHeaderAll(t *testing.T) {
	t.Parallel()

	var h RequestHeader
	r := bytes.NewBufferString("GET / HTTP/1.1\r\nHost: aa.com\r\nXX: YYY\r\nXX: ZZ\r\nCookie: a=b; c=d\r\nTrailer: Foo, Bar\r\n\r\n")
	if err := h.Read(bufio.NewReader(r)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	headers := collectSeq2(h.All())
	if value, ok := findHeader(headers, HeaderHost); !ok || value != "aa.com" {
		t.Fatalf("unexpected host %q, %v. Expected aa.com", value, ok)
	}
	if value, ok := findHeader(headers, HeaderCookie); !ok || value != "a=b; c=d" {
		t.Fatalf("unexpected cookie header %q, %v. Expected a=b; c=d", value, ok)
	}
	if value, ok := findHeader(headers, HeaderTrailer); !ok || value != "Foo, Bar" {
		t.Fatalf("unexpected trailer %q, %v. Expected Foo, Bar", value, ok)
	}

	var xx []string
	for key, value := range h.All() {
		if string(key) == "Xx" {
			xx = append(xx, string(value))
		}
	}
	if !reflect.DeepEqual(xx, []string{"YYY", "ZZ"}) {
		t.Fatalf("unexpected XX headers %v. Expected [YYY ZZ]", xx)
	}

	if got := collectSeq(h.AllTrailer()); !reflect.DeepEqual(got, []string{"Foo", "Bar"}) {
		t.Fatalf("unexpected trailers %v. Expected [Foo Bar]", got)
	}
	if got := collectSeq2(h.AllCookie()); !reflect.DeepEqual(got, [][2]string{{"a", "b"}, {"c", "d"}}) {
		t.Fatalf("unexpected cookies %v. Expected a=b and c=d", got)
	}
}

func TestHeaderIteratorsBreak(t *testing.T) {
	t.Parallel()

	var response ResponseHeader
	response.Add("First", "1")
	response.Add("Second", "2")
	response.AddTrailer("Foo")
	response.AddTrailer("Bar")
	var responseCookie Cookie
	responseCookie.SetKey("a")
	responseCookie.SetValue("b")
	response.SetCookie(&responseCookie)

	visited := 0
	for range response.All() {
		visited++
		break
	}
	if visited != 1 {
		t.Fatalf("unexpected number of response headers before break: %d. Expected 1", visited)
	}

	visited = 0
	for range response.AllTrailer() {
		visited++
		break
	}
	if visited != 1 {
		t.Fatalf("unexpected number of response trailers before break: %d. Expected 1", visited)
	}

	visited = 0
	for range response.AllCookie() {
		visited++
		break
	}
	if visited != 1 {
		t.Fatalf("unexpected number of response cookies before break: %d. Expected 1", visited)
	}

	var request RequestHeader
	request.SetHost("example.com")
	request.Add("Second", "2")
	request.AddTrailer("Foo")
	request.AddTrailer("Bar")
	request.SetCookie("a", "b")
	request.SetCookie("c", "d")

	visited = 0
	for range request.All() {
		visited++
		break
	}
	if visited != 1 {
		t.Fatalf("unexpected number of request headers before break: %d. Expected 1", visited)
	}

	visited = 0
	for range request.AllTrailer() {
		visited++
		break
	}
	if visited != 1 {
		t.Fatalf("unexpected number of request trailers before break: %d. Expected 1", visited)
	}

	visited = 0
	for range request.AllCookie() {
		visited++
		break
	}
	if visited != 1 {
		t.Fatalf("unexpected number of request cookies before break: %d. Expected 1", visited)
	}
}
