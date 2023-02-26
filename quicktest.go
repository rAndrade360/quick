package quick

import (
	"bytes"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type QuickTestReturn interface {
	Body() []byte
	BodyStr() string
	StatusCode() int
	Response() *http.Response
}

type qTest struct {
	body       []byte
	bodyStr    string
	statusCode int
	response   *http.Response
}

// QuickTest: This Method is a helper function to make tests with quick more quickly
// Required Params: Method (GET, POST, PUT, DELETE...), URI (only the path. Example: /test/:myParam)
// Optional Param: Body (If you don't want to define one, just ignore)
func (q Quick) QuickTest(method, URI string, headers map[string]string, body ...[]byte) (QuickTestReturn, error) {

	rand.Seed(time.Now().UnixNano())
	min := 3000
	max := 9999
	randPort := rand.Intn(max-min+1) + min

	port := strconv.Itoa(randPort)

	var buffBody []byte
	if len(body) > 0 {
		buffBody = body[0]
	}

	port = ConcatStr(":", port)
	URI = ConcatStr("http://0.0.0.0", port, URI)

	req, err := http.NewRequest(method, URI, io.NopCloser(bytes.NewBuffer(buffBody)))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if err != nil {
		return nil, err
	}

	go q.Listen(port)

	// This is a wait time to start the server in go routine
	time.Sleep(time.Millisecond * 100)

	c := http.DefaultClient

	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &qTest{
		body:       b,
		bodyStr:    string(b),
		statusCode: resp.StatusCode,
		response:   resp,
	}, nil
}

func (qt *qTest) Body() []byte {
	return qt.body
}

func (qt *qTest) BodyStr() string {
	return qt.bodyStr
}

func (qt *qTest) StatusCode() int {
	return qt.statusCode
}

func (qt *qTest) Response() *http.Response {
	return qt.response
}