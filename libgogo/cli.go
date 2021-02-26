package libgogo

import (
	"errors"
	"strconv"
)

/*
 * Returns a Server object with configuration as specified from given args
 * Possible Errors:
 *   Invalid arguments
 *   Missing values
 *   Root does not exist or is not directory
 */
func ParseArgs(args []string) (*Server, error) {
	var err error
	addr      := "127.0.0.1"
	port      := 7000
	host      := "localhost"
	root      := "."
	strict    := false
	blocklist := ""

	for i := 0; i < len(args); i++ {
		switch a := args[i]; a {
		case "--address":
			fallthrough
		case "-a":
			if i + 1 >= len(args) {
				return nil, errors.New("Address not specified" + helpMsg())
			}
			addr = args[i + 1]
			i = i + 1
		case "--port":
			fallthrough
		case "-p":
			if i + 1 >= len(args) {
				return nil, errors.New("Port not specified" + helpMsg())
			}
			port, err = strconv.Atoi(args[i + 1])
			if err != nil {
				return nil, errors.New("Port must be a number" + helpMsg())
			}
			i = i + 1
		case "--host":
			fallthrough
		case "-h":
			if i + 1 >= len(args) {
				return nil, errors.New("Hostname not specified\n" + helpMsg())
			}
			host = args[i + 1]
			i = i + 1
		case "--root":
			fallthrough
		case "-r":
			if i + 1 >= len(args) {
				return nil, errors.New("Root not specified\n" + helpMsg())
			}
			root = args[i + 1]
			i = i + 1
		case "--strict":
			fallthrough
		case "-s":
			strict = true
		case "--block":
			fallthrough
		case "-b":
			blocklist = args[i + 1]
			i = i + 1
		case "--help":
			fallthrough
		case "-?":
			return nil, errors.New(helpMsg())
		default:
			return nil, errors.New("Unrecognized arguments\n" + helpMsg())
		}
	}

	s, err := NewServer(addr, port, host, root, strict, blocklist)
	if err != nil {
		return nil, err
	}

	return s, nil
}

/*
 * Returns help message and documented configuration arguments
 */
func helpMsg() string {
	return "Usage: gogo [-?s] [-a address] [-p port] [-h hostname] [-r root]" +
	       "\nOptions:\n" +
	       "    -?, --help     Print this help message\n" +
	       "    -a, --address  IP address to listen on\n" +
	       "    -p, --port     TCP port to listen on\n" +
	       "    -h, --host     Hostname to identify with\n" +
	       "    -r, --root     Directory to use as root\n" +
	       "    -s, --strict   Do not perform interpolation (host, port, etc.)\n" +
	       "    -b, --block    Name of file containing list of blocked IP addresses"
}
