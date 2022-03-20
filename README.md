# GoGopher.Go

A very simple Gopher server written in Go. Although not very featureful, it is
functional and secure to the best of my knowledge.

One feature it does have is hostname and port interpolation. Rather than
entering literal hostnames and ports in all your gophermaps, you can use \host
and \port and the server will substitute the actual hostname and port it is
listening on.

## Building:

Just run `go build`.

## Usage:

```
Usage: gogopher [-?s] [-a address] [-p port] [-h hostname] [-r root]
Usage: gogopher [-?s] [-h hostname] [-p port] [-r root] <address>:<port>" +
Options:
	-?, --help     Print this help message\n" +
	-h, --host     Hostname to use for interpolation\n" +
	-p, --port     TCP port to use for interpolation\n" +
	-r, --root     Directory to use as root\n" +
	-s, --strict   Do not perform interpolation (host, port, etc.)\n" +
	-b, --block    Name of file containing list of blocked IP addresses\n" +
	-l, --log      Name of file to direct logs (default is stdout)"
```
