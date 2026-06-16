package test

import (
	"bufio"
	"net"
	"redis-mini/internal/handlers"
	"reflect"
	"testing"
	"time"
)

// writePayload spins up one end of an in-memory connection, writes the raw
// RESP bytes to it, then closes it. HandleResp reads from the other end.
// net.Pipe gives us two connected net.Conns with no real sockets involved,
// so these tests need no listener and no ports.
func writePayload(t *testing.T, payload string) net.Conn {
	t.Helper()
	server, client := net.Pipe()

	go func() {
		// net.Pipe is synchronous: this Write blocks until HandleResp's
		// bufio.Reader pulls the bytes, so the goroutine is required.
		_, _ = client.Write([]byte(payload))
		_ = client.Close()
	}()

	return server
}

func TestHandleResp(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    []string
		wantErr bool
	}{
		{
			name:    "single element PING",
			payload: "*1\r\n$4\r\nPING\r\n",
			want:    []string{"PING"},
		},
		{
			name:    "two element GET",
			payload: "*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n",
			want:    []string{"GET", "mykey"},
		},
		{
			name:    "three element SET (the original test case)",
			payload: "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n",
			want:    []string{"SET", "mykey", "myvalue"},
		},
		{
			name: "multi-digit bulk string length",
			// "hello world" is 11 chars — exercises the $1X two-digit length
			// path that the single-ReadByte version got wrong.
			payload: "*3\r\n$3\r\nSET\r\n$3\r\nmsg\r\n$11\r\nhello world\r\n",
			want:    []string{"SET", "msg", "hello world"},
		},
		{
			name: "multi-digit array length (12 elements)",
			// 12 elements forces the array count past one digit too.
			payload: "*12\r\n" +
				"$2\r\nc0\r\n$2\r\nc1\r\n$2\r\nc2\r\n$2\r\nc3\r\n" +
				"$2\r\nc4\r\n$2\r\nc5\r\n$2\r\nc6\r\n$2\r\nc7\r\n" +
				"$2\r\nc8\r\n$2\r\nc9\r\n$3\r\nc10\r\n$3\r\nc11\r\n",
			want: []string{
				"c0", "c1", "c2", "c3", "c4", "c5",
				"c6", "c7", "c8", "c9", "c10", "c11",
			},
		},
		{
			name:    "empty bulk string",
			payload: "*2\r\n$5\r\nHELLO\r\n$0\r\n\r\n",
			want:    []string{"HELLO", ""},
		},
		{
			name: "value containing spaces and punctuation",
			// The byte-count framing means the data is read verbatim — no
			// tokenizing on spaces, so embedded spaces must survive intact.
			payload: "*3\r\n$3\r\nSET\r\n$4\r\nnote\r\n$24\r\nbuy milk, eggs, & bread!\r\n",
			want:    []string{"SET", "note", "buy milk, eggs, & bread!"},
		},
		{
			name:    "non-numeric array length",
			payload: "*abc\r\n",
			// readLine -> "abc", strconv.Atoi fails. The parser must surface
			// this rather than treating it as length 0.
			wantErr: true,
		},
		{
			name:    "non-numeric bulk string length",
			payload: "*1\r\n$xyz\r\n",
			// Array count parses fine (1), but the bulk length "xyz" is junk.
			wantErr: true,
		},
		{
			name:    "truncated: claims 3 elements, supplies 1",
			payload: "*3\r\n$3\r\nSET\r\n",
			// After SET, the loop tries to read element 2 and hits EOF.
			// ReadByte returns io.EOF, which must propagate as an error.
			wantErr: true,
		},
		{
			name:    "bulk length longer than the data provided",
			payload: "*1\r\n$10\r\nhi\r\n",
			// Claims 10 bytes, only 2 are present before EOF. io.ReadFull
			// returns io.ErrUnexpectedEOF.
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := writePayload(t, tt.payload)
			defer conn.Close()

			// Guard against a parser bug that blocks forever waiting on bytes
			// (e.g. a wrong length causing io.ReadFull to hang). Without a
			// deadline a broken parser would stall the whole test run.
			_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))

			got, _, err := handlers.HandleResp(bufio.NewReader(conn))

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got nil (args=%v)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleResp() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
