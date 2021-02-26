package libgogo

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Server struct {
	sock   net.Listener
	addr   string
	port   int
	host   string
	root   string
	strict bool
}

/*
 * Returns a new Server object
 * Possible Errors:
 *   Root doesn't exist or is not a directory
 */
func NewServer(a string, p int, h string, r string, s bool) (*Server, error) {
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

	return &Server{
		sock:   nil,
		addr:   a,
		port:   p,
		host:   h,
		root:   root,
		strict: s,
	}, nil
}

/*
 * Dispatches TCP requests to the Gopher handler
 * Possible Errors:
 *   TCP socket errors
 */
func (s *Server) Go() error {
	var err error
	s.sock, err = net.Listen("tcp4", s.addr + ":" + strconv.Itoa(s.port))
	if err != nil {
		return err
	}

	for {
		conn, err := s.sock.Accept()
		if err != nil {
			log.Printf("Error: %s", err.Error())
			continue
		}

		// kinda quick and dirty :/
		go func() {
			defer conn.Close()
			req := make([]byte, 64)
			n, err := conn.Read(req)
			if err != nil {
				log.Printf("Error: %s", err.Error())
				return
			}

			/* Trim newlines or CRLFs */
			req = req[:n-1]
			if n > 1 {
				if string(req[n-2]) == "\r" {
					req = req[:n-2]
				}
			}

			log.Printf("%s: %s", conn.RemoteAddr().String(), string(req))
			_, err = conn.Write(s.handle(req))
			if err != nil {
				log.Printf("Error: %s", err.Error())
				return
			}
		}()
	}

	return nil
}

/*
 * Returns appropriate response to given request according to Gopher protocol
 */
func (s *Server) handle(req []byte) []byte {
	var res []byte
	path := filepath.Join(s.root, string(req))

	/* Forbid leaving root dir */
	if len(path) < len(s.root) {
		return s.handle([]byte("/"))
	}
	if path[:len(s.root)] != s.root {
		return s.handle([]byte("/"))
	}

	info, err := os.Stat(path)
	if err != nil {
		res = []byte(notFound(string(req)))
	} else if info.IsDir() {
		res, err = ioutil.ReadFile(filepath.Join(path, "gophermap"))
	} else {
		res, err = ioutil.ReadFile(path)
	}
	if err != nil {
		res = []byte(notFound(string(req)))
	}

	if http.DetectContentType(res)[:5] == "text/" {
		return []byte(s.interpolate(string(res)))
	}

	return res
}

/*
 * Returns the given response with \tags replaced with values from the server.
 * If the strict flag is set, the only interpolation that takes place is
 * replacing Unix newlines with CRLF.
 * Interpolation Tags:
 *   \host -> server's hostname
 *   \port -> server's port number
 */
func (s *Server) interpolate(res string) string {
	if !s.strict {
		res = strings.Replace(res, "\\host", s.host, -1)
		res = strings.Replace(res, "\\port", strconv.Itoa(s.port), -1)
	}

	return strings.Replace(res, "\n", "\r\n", -1)
}

/*
 * Returns standard "not found" response
 */
func notFound(req string) string {
	return "3'" + req +
	       "' does not exist (no handler found)\t\terror.host\t1\r\n"
}
