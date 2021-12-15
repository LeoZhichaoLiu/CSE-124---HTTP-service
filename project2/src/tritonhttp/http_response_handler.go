package tritonhttp

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func (hs *HttpServer) handleBadRequest(conn net.Conn) {
	//panic("todo - handleBadRequest")

	response := HttpResponseHeader{
		Version: "HTTP/1.1",
		Code:    400,
		Message: "Bad Request",
		Content: make(map[string]string),
	}

	response.Content["Date"] = time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST")
	response.Content["Connection"] = "close"
	hs.sendResponse(response, conn)

}

func (hs *HttpServer) handleFileNotFoundRequest(requestHeader *HttpRequestHeader, conn net.Conn) {
	//panic("todo - handleFileNotFoundRequest")

	response := HttpResponseHeader{
		Version: "HTTP/1.1",
		Code:    404,
		Message: "page not found",
		Content: make(map[string]string),
	}

	response.Content["Date"] = time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST")
	hs.sendResponse(response, conn)
}

func (hs *HttpServer) handleResponse(requestHeader *HttpRequestHeader, conn net.Conn, path string) (close bool) {
	//panic("todo - handleResponse")

	close = false
	contentType := getType(requestHeader.Url, hs.MIMEMap)
	last_modified := getModified(path)
	bodyString := getBody(path)

	response := HttpResponseHeader{
		Version: "HTTP/1.1",
		Code:    200,
		Message: "OK",
		body:    bodyString,
		Content: make(map[string]string),
	}

	response.Content["Content-Length"] = strconv.Itoa(len(bodyString))
	response.Content["Content-Type"] = contentType
	response.Content["Last-Modified"] = last_modified
	response.Content["Date"] = time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST")

	if checkClose(requestHeader) {
		response.Content["Connection"] = "close"
		close = true
	}

	hs.sendResponse(response, conn)

	return close
}

func (hs *HttpServer) sendResponse(responseHeader HttpResponseHeader, conn net.Conn) {
	//panic("todo - sendResponse")

	message := responseHeader.Version + " " + strconv.Itoa(responseHeader.Code) + " " + responseHeader.Message + "\r\n"
	for k, v := range responseHeader.Content {
		message = message + k + ": " + v + "\r\n"
	}
	message = message + "\r\n"
	if responseHeader.body != "" {
		message = message + responseHeader.body
	}

	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString(message)

	if err != nil {
		return
	}
	err = writer.Flush()

	// Send headers
	// Send file if required
	// Hint - Use the bufio package to write response
}

func checkClose(request *HttpRequestHeader) bool {
	val, ok := request.Content["Connection"]

	if ok == true && val == "close" {
		return true
	}
	return false
}

func getModified(path string) string {

	file, err := os.Stat(path)

	if err != nil {
		log.Panicln(err)
	}
	modifiedtime := file.ModTime()

	return modifiedtime.Format("Mon, 02 Jan 2006 15:04:05 MST")
}

func getType(path string, MIMEMap map[string]string) string {

	entries := strings.Split(path, ".")

	if len(entries) == 1 {
		return "application/octet-stream"
	}

	extension := "." + entries[1]

	val, ok := MIMEMap[extension]
	contentType := "application/octet-stream"

	if ok == true {
		contentType = val
	}

	return contentType
}

func getBody(path string) string {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	bodyString := string(content)

	/*
		f, err := os.Open(path)
		if err != nil {
			log.Panicln(err)
		}

		defer f.Close()

		bodyString := ""
		// read the file line by line using scanner
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			bodyString = bodyString + scanner.Text() + "\n"
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}*/

	return bodyString
}
