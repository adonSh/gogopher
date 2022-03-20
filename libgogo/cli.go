package libgogo

import (
	"errors"
	"strconv"
	"strings"
)

// Returns a Server object with configuration as specified from given args
// Possible Errors:
//   Invalid arguments
//   Missing values
//   Root does not exist or is not a directory
func ParseArgs(args []string) (*Server, error) {
	var err error
	addr      := "127.0.0.1"
	lport     := 7000
	host      := "localhost"
	port      := lport
	root      := "."
	strict    := false
	blocklist := ""
	logfile   := ""

	for i := 0; i < len(args); i++ {
		switch a := args[i]; a {
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
		case "--log":
			fallthrough
		case "-l":
			logfile = args[i + 1]
			i = i + 1
		case "--help":
			fallthrough
		case "-?":
			return nil, errors.New(helpMsg())
		default:
			ap := strings.Split(a, ":")
			addr = ap[0]
			if len(ap) > 1 {
				lport, err = strconv.Atoi(ap[1])
				if err != nil {
					return nil, errors.New("Port must be a number" + helpMsg())
				}
			}
		}
	}

	s, err := NewServer(addr, lport, host, port, root, strict, blocklist, logfile)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Returns help message and documented configuration arguments
func helpMsg() string {
	return "Usage: gogopher [-?s] [-h hostname] [-p port] [-r root] <address>:<port>" +
	       "\nOptions:\n" +
	       "    -?, --help     Print this help message\n" +
	       "    -h, --host     Hostname to use for interpolation\n" +
	       "    -p, --port     TCP port to use for interpolation\n" +
	       "    -r, --root     Directory to use as root\n" +
	       "    -s, --strict   Do not perform interpolation (host, port, etc.)\n" +
	       "    -b, --block    Name of file containing list of blocked IP addresses\n" +
	       "    -l, --log      Name of file to direct logs (default is stdout)"
}
