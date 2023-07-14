package libgogo

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// The Server object
type Server struct {
	sock      net.Listener // The actual socket
	Addr      string       // The address to listen on
	Lport     int          // The port to listen on
	host      string       // The hostname to use for interpolation
	port      int          // The port to use for interpolation
	root      string       // The root directory to serve
	strict    bool         // To interpolate or not
	blocklist []string     // IPs to block
	logger    *log.Logger  // The logger
}

// Returns a new Server object
// Possible Errors:
//   Root doesn't exist or is not a directory
func NewServer(a string, lp int, h string, p int, r string, s bool, bl string, lf string) (*Server, error) {
	root, err := filepath.Abs(r)
	if err != nil {
		return nil, err
	}
	rInfo, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !rInfo.IsDir() {
		return nil, err
	}

	blist := []string{}
	if bl != "" {
		file, err := os.Open(bl)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		blist = blocklistFromFile(file)
	}

	logfile := os.Stdout
	if lf != "" {
		logfile, err = os.OpenFile(lf, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
	}

	return &Server{
		sock:      nil,
		Addr:      a,
		Lport:     lp,
		host:      h,
		port:      p,
		root:      root,
		strict:    s,
		blocklist: blist,
		logger:    log.New(logfile, "", log.Flags()),
	}, nil
}

// Returns a slice of all lines in provided file. If any errors are
// encountered an empty slice is returned.
func blocklistFromFile(file *os.File) []string {
	bl := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		bl = append(bl, scanner.Text())
	}
	if scanner.Err() != nil {
		return []string{}
	}

	return bl
}

// Dispatches TCP connections to the Gopher handler
// Possible Errors:
//   TCP socket errors
func (s *Server) Go() error {
	var err error
	s.sock, err = net.Listen("tcp4", s.Addr + ":" + strconv.Itoa(s.Lport))
	if err != nil {
		return err
	}

	for {
		conn, err := s.sock.Accept()
		if err != nil {
			return err
		}
		if s.isBlocked(conn.RemoteAddr().String()) {
			conn.Close()
			continue
		}

		go s.handle(conn)
	}

	return nil
}

func (s *Server) handle(client net.Conn) {
	defer client.Close()

	req := make([]byte, 64)
	n, err := client.Read(req)
	if err != nil {
		s.logger.Printf("Error: %s", err.Error())
		return
	}

	// Trim newlines or CRLFs
	req = bytes.TrimSpace(req[:n])

	s.logger.Printf("%s: %s", client.RemoteAddr().String(), string(req))
	_, err = client.Write(s.render(req))
	if err != nil {
		s.logger.Printf("Error: %s", err.Error())
		return
	}
}

func httpRedirect(url string) []byte {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta http-equiv="refresh" content="1; url=URL">
  <title>redirect</title>
</head>
<body>
  <p>
    <a href="URL">URL</a>
  </p>
</body>
</html>
`
	return []byte(strings.Replace(html, "URL", url, -1))
}

// Returns appropriate response to given request according to Gopher protocol
func (s *Server) render(req []byte) []byte {
	if strings.HasPrefix(string(req), "URL:") {
		return httpRedirect(string(req)[4:])
	}

	var res []byte
	path := filepath.Join(s.root, string(req))

	// Forbid leaving root dir
	if len(path) < len(s.root) {
		return s.render([]byte("/"))
	}
	if path[:len(s.root)] != s.root {
		return s.render([]byte("/"))
	}

	info, err := os.Stat(path)
	if err != nil {
		res = []byte(notFound(string(req)))
	} else if info.IsDir() {
		res, err = os.ReadFile(filepath.Join(path, "gophermap"))
	} else {
		res, err = os.ReadFile(path)
	}
	if err != nil {
		res = []byte(notFound(string(req)))
	}

	if strings.HasPrefix(http.DetectContentType(res), "text/") {
		return []byte(s.normalize(string(res)))
	}

	return res
}

// Returns the given response with \tags replaced with values from the server.
// If the strict flag is set, the only interpolation that takes place is
// replacing Unix newlines with CRLF.
// Interpolation Tags:
//   \host -> server's hostname
//   \port -> server's port number
func (s *Server) normalize(res string) string {
	if !s.strict {
		res = strings.Replace(res, "\\host", s.host, -1)
		res = strings.Replace(res, "\\port", strconv.Itoa(s.port), -1)
	}

	return strings.Replace(res, "\n", "\r\n", -1)
}

// Returns true if the given address is in the Server's blocklist
func (s *Server) isBlocked(addr string) bool {
	for i := 0; i < len(s.blocklist); i++ {
		if addr[:len(s.blocklist[i])] == s.blocklist[i] {
			return true
		}
	}

	return false
}

// Returns standard "not found" response
func notFound(req string) string {
	msg := fmt.Sprintf("3'%s' does not exist (no handler found)", req)
	return strings.Join([]string{msg, "error.host", "1\n"}, "\t")
}
