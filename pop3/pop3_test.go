package pop3

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

type fakeAddr struct{}

func (fakeAddr) Network() string { return "" }
func (fakeAddr) String() string  { return "" }

type connFaker struct {
	Buffer *bytes.Buffer
	Writer *bufio.Writer
	io.ReadWriter
}

func (f connFaker) Close() error {
	return nil
}

func (f connFaker) LocalAddr() net.Addr {
	return fakeAddr{}
}

func (f connFaker) RemoteAddr() net.Addr {
	return fakeAddr{}
}

func (f connFaker) SetDeadline(t time.Time) error {
	return nil
}

func (f connFaker) SetReadDeadline(t time.Time) error {
	return nil
}

func (f connFaker) SetWriteDeadline(t time.Time) error {
	return nil
}

func initializeFakeConn(responseLines string) *connFaker {
	basicServer := strings.Join(strings.Split(responseLines, "\n"), "\r\n")
	var commandBuffer bytes.Buffer
	bufferWriter := bufio.NewWriter(&commandBuffer)
	var fake = &connFaker{Buffer: &commandBuffer, Writer: bufferWriter}
	fake.ReadWriter = bufio.NewReadWriter(bufio.NewReader(strings.NewReader(basicServer)), bufferWriter)

	return fake
}

func TestBasic(t *testing.T) {
	basicClient := strings.Join(strings.Split(basicClient, "\n"), "\r\n")

	fakeConn := initializeFakeConn(basicServer)

	c, err := NewClient(fakeConn)
	if err != nil {
		t.Fatalf("NewClient failed: %s", err)
	}

	if err = c.User("uname"); err != nil {
		t.Fatal("User failed: ", err)
	}

	if err = c.Pass("password1"); err == nil {
		t.Fatal("Pass succeeded inappropriately")
	}

	if err = c.Auth("uname", "password2"); err != nil {
		t.Fatal("Auth failed: ", err)
	}

	if err = c.Noop(); err != nil {
		t.Fatal("Noop failed: ", err)
	}

	fakeConn.Writer.Flush()
	if basicClient != fakeConn.Buffer.String() {
		t.Fatalf("Got:\n%s\nExpected:\n%s", fakeConn.Buffer.String(), basicClient)
	}
}

var basicServer = `+OK good morning
+OK send PASS
-ERR [AUTH] mismatched username and password
+OK send PASS
+OK welcome
+OK
`

var basicClient = `USER uname
PASS password1
USER uname
PASS password2
NOOP
`
