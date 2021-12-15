package tritonhttp

import (
	"log"
	"net"
	"os"
	"strings"
	"time"
)

/*
For a connection, keep handling requests until
	1. a timeout occurs or
	2. client closes connection or
	3. client sends a bad request
*/
func (hs *HttpServer) handleConnection(conn net.Conn, delim *string) {

	log.Println("Accepted new connection.")

	remaining := ""

	// We make a for loop to continuing read message from socket.
	for {

		// We create a buffer to store reading request
		buffer := make([]byte, 1024)

		// We set the timeout as 5 seconds from current time.
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))

		// We start reading the request from socket to buffer
		bytes, err := conn.Read(buffer[0:])

		if err != nil {
			if len(remaining) != 0 {
				hs.handleBadRequest(conn)
			}
			//fmt.Fprintf(os.Stderr, "Conn::Read: err %v\n", err)
			conn.Close()
			return
		}

		all := string(buffer[:bytes])
		remaining = remaining + all

		//per_request := strings.Split(all, "\r\n\r\n")

		//for _, s := range per_request {

		for strings.Contains(remaining, *delim) {

			idx := strings.Index(remaining, *delim)
			s := remaining[:idx]
			remaining = remaining[idx+4:]

			//log.Println(s)
			//log.Println(len(s))
			if len(s) == 0 {
				break
			}

			// After reciving request, we construct a request struct
			var request *HttpRequestHeader
			request = new(HttpRequestHeader)

			// We split the buffer string by CRLF
			entries := strings.Split(s, "\r\n")

			// We split the first initial line by space (Method, Url, Version)
			Line := strings.Split(entries[0], " ")

			if len(Line) != 3 {
				log.Println("Length not valid")
				hs.handleBadRequest(conn)
				conn.Close()
				return
			}

			// We set the Method, Url, and Version in request struct
			request.Method = Line[0]
			request.Url = Line[1]
			request.Version = Line[2]

			// We then try to fill out the content map in request
			m := make(map[string]string)

			if len(entries) == 1 {
				hs.handleBadRequest(conn)
				conn.Close()
				return
			}

			// We loop through the rest of the entries, and spilt them by ": "
			for _, e := range entries[1:] {

				parts := strings.Split(e, ": ")

				if len(parts) != 2 {
					break
				}

				// We assign the content map in request struct
				m[parts[0]] = parts[1]
				request.Content = m
			}

			// After making the request struct, we try to check whether format is correct
			formatValid := checkValid(request)

			pathValid := false
			path := ""

			// If the format is valid, we try to check whether the url is valid.
			if formatValid == true {
				if request.Url == "/" {
					request.Url = "/index.html"
				}
				pathValid, path = checkPath(request.Url, hs.DocRoot[request.Content["Host"]])
			}

			if formatValid == true {

				// If the format and path are valid, we handle 200
				if pathValid == true {
					// If the request contains connection close, close socket
					if hs.handleResponse(request, conn, path) == true {
						conn.Close()
						return
					}

					// If the format is not valid, we handle 404 not found
				} else {
					hs.handleFileNotFoundRequest(request, conn)
				}

				// Otherwise, we just handle 400 client error
			} else {
				hs.handleBadRequest(conn)
				conn.Close()
				return
			}
		}
	}

	conn.Close()

	// Start a loop for reading requests continuously
	// Set a timeout for read operation
	// Read from the connection socket into a buffer
	// Validate the request lines that were read
	// Handle any complete requests
	// Update any ongoing requests
	// If reusing read buffer, truncate it before next reach.
}

func checkValid(requestHeader *HttpRequestHeader) (result bool) {
	if requestHeader.Method == "" || requestHeader.Url == "" || requestHeader.Version == "" {
		return false
	}

	if requestHeader.Method != "GET" || requestHeader.Url[0] != '/' {
		return false
	}

	entries := strings.Split(requestHeader.Version, "/")
	if len(entries) != 2 || entries[0] != "HTTP" || entries[1] != "1.1" {
		return false
	}

	_, ok := requestHeader.Content["Host"]

	if ok == false {
		return false
	}

	for k, v := range requestHeader.Content {
		if k == "" || v == "" {
			return false
		}
	}

	return true
}

func checkPath(url string, docRoot string) (bool, string) {

	path := docRoot + url

	_, err := os.Stat(path)

	if err != nil {
		return false, ""
	}

	if os.IsNotExist(err) {
		return false, ""
	}

	if docRoot == "" {
		return false, ""
	}

	if url == "" {
		return false, ""
	}

	return true, path
}
