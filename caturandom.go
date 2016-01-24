package main

/*
 * caturandom.go
 * Spews out /dev/urandom to connecting clients, and prints what they send
 * By J. Stuart McMurray
 * Created 20150124
 * Last Modified 20150124
 */

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	var (
		randF = flag.String(
			"f",
			"/dev/urandom",
			"File from which to send contents",
		)
		addr = flag.String(
			"a",
			":7347",
			"Listen address",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %v [options]

Listens on the specified address and sends contents of specified file to
connected clients.  Logs what clients send.

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Better logging */
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	/* Listen on the given address */
	l, err := net.Listen("tcp", *addr)
	if nil != err {
		log.Fatalf("Unable to listen on %v: %v", *addr, err)
	}
	log.Printf("Listening on %v", l.Addr())
	log.Printf("Will send contents of %v to connected clients", *randF)

	/* Handle connecting clients */
	for {
		c, err := l.Accept()
		if nil != err {
			log.Fatalf(
				"Unable to accept connections on %v: %v",
				l.Addr(),
				err,
			)
		}
		go handle(c, *randF)
	}
}

/* handle sends the contents of the file named fn to c, and logs data from
c. */
func handle(c net.Conn, fn string) {
	log.Printf("%v Start", c.RemoteAddr())
	defer c.Close()

	/* Try to open the file */
	f, err := os.Open(fn)
	if nil != err {
		log.Printf("Unable to open %v: %v", fn, err)
		return
	}

	/* Print what's sent to us */
	go logRecv(c)

	/* Send bytes to them */
	n, err := io.Copy(c, f)
	if nil != err && !strings.HasSuffix(
		err.Error(),
		"write: broken pipe",
	) {
		log.Printf(
			"%v Error after sending %v bytes (%T): %v",
			c.RemoteAddr(),
			n,
			err,
			err,
		)
	} else {
		log.Printf("%v Sent %v bytes", c.RemoteAddr(), n)
	}
}

/* logRecv logs messages received on c */
func logRecv(c net.Conn) {
	var (
		b   = make([]byte, 2048)
		n   int
		err error
	)
	for {
		/* Try to read */
		n, err = c.Read(b)
		/* Log what we have */
		if 0 != n {
			log.Printf("%v %q", c.RemoteAddr(), b[:n])
		}
		/* Exit on error */
		if nil != err {
			if io.EOF != err && !strings.HasSuffix(
				err.Error(),
				"read: connection reset by peer",
			) {
				log.Printf(
					"%v Receive error: %v",
					c.RemoteAddr(),
					err,
				)
			}
			return
		}
	}
}
