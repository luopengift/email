package pop3

import (
	"io"
	"strings"
	"testing"
)

type readWriteFaker struct {
	io.Reader
	io.Writer
	io.Closer
}

type fakeCloser struct{}

func (f fakeCloser) Close() error { return nil }

type fakeWriter struct {
	buffer *[]byte
}

func (f fakeWriter) Write(p []byte) (int, error) {
	*f.buffer = append(*f.buffer, p...)
	return len(p), nil
}
func (f fakeWriter) Flush() error { return nil }

func TestCmdOK(t *testing.T) {
	responseOK := "+OK\r\n"
	command := "STAT"

	var buffer []byte
	buffer = make([]byte, 0)
	var fake readWriteFaker
	fake.Reader = strings.NewReader(responseOK)
	fake.Writer = &fakeWriter{buffer: &buffer}
	fake.Closer = &fakeCloser{}

	connection := NewConnection(fake)

	response, err := connection.Cmd(command)
	if err != nil {
		t.Fatalf("Got:\n%v\nExpected: nil\n", err)
	}

	if response != "" {
		t.Fatalf("Got:\n%v\nExpected: \"\"\n", response)
	}

	expectedCommand := command + "\r\n"
	if string(buffer) != expectedCommand {
		t.Fatalf("Got:\n%s\nExpected: %s\n", string(buffer), expectedCommand)
	}
}

func TestCmdOKWithArgs(t *testing.T) {
	responseOK := "+OK\r\n"
	command := "LIST %d"

	var buffer []byte
	buffer = make([]byte, 0)
	var fake readWriteFaker
	fake.Reader = strings.NewReader(responseOK)
	fake.Writer = &fakeWriter{buffer: &buffer}
	fake.Closer = &fakeCloser{}

	connection := NewConnection(fake)

	response, err := connection.Cmd(command, 1)
	if err != nil {
		t.Fatalf("Got:\n%v\nExpected: nil\n", err)
	}

	if response != "" {
		t.Fatalf("Got:\n%v\nExpected: \"\"\n", response)
	}

	expectedCommand := "LIST 1\r\n"
	if string(buffer) != expectedCommand {
		t.Fatalf("Got:\n%s\nExpected: %s\n", string(buffer), expectedCommand)
	}
}

func TestCmdOKWithMessage(t *testing.T) {
	responseOK := "+OK 5 messages:\r\n"

	var buffer []byte
	buffer = make([]byte, 0)
	var fake readWriteFaker
	fake.Reader = strings.NewReader(responseOK)
	fake.Writer = &fakeWriter{buffer: &buffer}
	fake.Closer = &fakeCloser{}

	connection := NewConnection(fake)

	response, err := connection.Cmd("STAT")
	if err != nil {
		t.Fatalf("Got:\n%v\nExpected: nil\n", err)
	}

	expected := "5 messages:"
	if response != expected {
		t.Fatalf("Got:\n%v\nExpected: %s\n", response, expected)
	}
}

func TestCmdERR(t *testing.T) {
	responseErr := "-ERR An error occured\r\n"

	var buffer []byte
	buffer = make([]byte, 0)
	var fake readWriteFaker
	fake.Reader = strings.NewReader(responseErr)
	fake.Writer = &fakeWriter{buffer: &buffer}
	fake.Closer = &fakeCloser{}

	connection := NewConnection(fake)

	_, err := connection.Cmd("STAT")
	if err == nil {
		t.Fatal("Expected an error")
	}

	expectedError := "An error occured"
	if err.Error() != expectedError {
		t.Fatalf("Got:\n%v\nExpected: %s\n", err.Error(), expectedError)
	}
}

func TestReadMultiLines(t *testing.T) {
	multiLineResponse := "line 1\r\nline 2\r\n."

	var fake readWriteFaker
	fake.Reader = strings.NewReader(multiLineResponse)
	fake.Writer = &fakeWriter{}
	fake.Closer = &fakeCloser{}

	connection := NewConnection(fake)

	lines, err := connection.ReadMultiLines()
	if err != nil {
		t.Fatalf("Got:\n%v\nExpected: nil\n", err)
	}

	if len(lines) != 2 {
		t.Fatalf("Got %d lines.\nExpected %d lines.\n", len(lines), 2)
	}

	if lines[0] != "line 1" {
		t.Fatalf("Got: %s\nExpected %s\n", lines[0], "line 1")
	}
}

func TestReadMultiLinesRemovesLeadingDot(t *testing.T) {
	multiLineResponse := ".line 1\r\n."

	var fake readWriteFaker
	fake.Reader = strings.NewReader(multiLineResponse)
	fake.Writer = &fakeWriter{}
	fake.Closer = &fakeCloser{}

	connection := NewConnection(fake)

	lines, err := connection.ReadMultiLines()
	if err != nil {
		t.Fatalf("Got:\n%v\nExpected: nil\n", err)
	}

	if lines[0] != "line 1" {
		t.Fatalf("Got: %s\nExpected %s\n", lines[0], "line 1")
	}
}
