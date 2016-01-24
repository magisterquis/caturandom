# caturandom
Splegs out the contents of `/dev/urandom` (or any other file) to connecting
clients.

Meant to be used as a catch-all for ports behind which there's no listening
service.  Anything the client sends will be logged.

Installation
------------
```bash
go get github.com/kd5pbo/caturandom
go install github.com/kd5pbo/caturandom
```
Compiled binaries available upon request.

Usage
-----
```
Usage: caturandom [options]

Listens on the specified address and sends contents of specified file to
connected clients.  Logs what clients send.

Options:
  -a string
    	Listen address (default ":7347")
  -f string
    	File from which to send contents (default "/dev/urandom")
```
