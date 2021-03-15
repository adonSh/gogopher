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
Options:
  -?, --help     Print this help message
  -a, --address  IP address to listen on
  -p, --port     TCP port to listen on
  -h, --host     Hostname to identify with
  -r, --root     Directory to use as root
  -s, --strict   Do not perform interpolation (host, port, etc.)
  -b, --block    Name of file containing list of blocked IP addresses
```
