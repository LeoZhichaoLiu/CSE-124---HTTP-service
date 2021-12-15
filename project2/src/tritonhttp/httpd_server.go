package tritonhttp

import (
	"log"
	"net"
	"flag"
)

/**
	Initialize the tritonhttp server by populating HttpServer structure
**/
func NewHttpdServer(port string, docRoot map[string]string, mimePath string) (*HttpServer, error) {

	// Initialize mimeMap for server to refer
	mimeMap, err := ParseMIME(mimePath)

	// We create a server struct, containg its port, root, and type-extension map.
	Server := &HttpServer{
		ServerPort: port,
		DocRoot:    docRoot,
		MIMEPath:   mimePath,
		MIMEMap:    mimeMap,
	}

	// Return pointer to HttpServer
	return Server, err
}

/**
	Start the tritonhttp server
**/
func (hs *HttpServer) Start() (err error) {

	// Start listening to the server port
	l, err := net.Listen("tcp", hs.ServerPort)

	if err != nil {
		log.Panicln(err)
	}

	defer l.Close()

	delim := flag.String("delimiter", "\r\n\r\n", "Delimiter used to separate request")
	flag.Parse()

	// Accept connection from client
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		// Spawn a go routine to handle request
		go hs.handleConnection(conn, delim)
	}
}
