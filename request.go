package vcr

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
)

type vcrRequest struct {
	// Header is intentionally not included and is not used for episode matching
	Method string
	URL    string
	Body   string
}

func newVCRRequest(request *http.Request, filterMap map[string]string) *vcrRequest {
	var body []byte
	if request.Body != nil {
		body, _ = ioutil.ReadAll(request.Body)
		request.Body.Close()
		request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		body = replaceBodyPattern(body, filterMap)
	}

	return &vcrRequest{
		Method: request.Method,
		URL:    request.URL.String(),
		Body:   string(body),
	}
}

func replaceBodyPattern(body []byte, patternFilters map[string]string) []byte {
	newBody := make([]byte, len(body))
	copy(newBody, body)
	for pattern, replacement := range patternFilters {
		r := regexp.MustCompile(pattern)
		newBody = r.ReplaceAllLiteral(newBody, []byte(replacement))
	}

	return newBody
}
