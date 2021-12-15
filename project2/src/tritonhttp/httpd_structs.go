package tritonhttp

type HttpServer struct {
	ServerPort string
	DocRoot    map[string]string
	MIMEPath   string
	MIMEMap    map[string]string
}

type HttpResponseHeader struct {
	// Add any fields required for the response here
	Version string
	Code    int
	Message string
	body    string
	Content map[string]string
}

type HttpRequestHeader struct {
	// Add any fields required for the request here
	Valid   bool
	Method  string
	Url     string
	Version string
	Content map[string]string
}
