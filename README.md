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
Usage: gogopher [-?s] [-h hostname] [-p port] [-r root] <address>:<port>
Options:
	-?, --help     Print this help message
	-h, --host     Hostname to use for interpolation
	-p, --port     TCP port to use for interpolation
	-r, --root     Directory to use as root
	-s, --strict   Do not perform interpolation (host, port, etc.)
	-b, --block    Name of file containing list of blocked IP addresses
	-l, --log      Name of file to direct logs (default is stdout)
```
