package msgraph

import (
	"bytes"
	"encoding/json"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/sirupsen/logrus"
	"westpac.co.nz/msgraph/pkg/helpers"
)

// BaseResourceAPI GraphBaseResourceAPI
type BaseResource struct {
	UserAgent string
	Version   string

	HTTPClient *http.Client
}

// {
// 	"error": {
// 	  "code": "InvalidAuthenticationToken",
// 	  "message": "CompactToken parsing failed with error code: 80049217",
// 	  "innerError": {
// 		"date": "2020-11-05T08:17:45",
// 		"request-id": "833ad94d-cf61-4b1e-ae40-4c947d4a1348",
// 		"client-request-id": "833ad94d-cf61-4b1e-ae40-4c947d4a1348"
// 	  }
// 	}
//   }

// GraphAPIErrorResponse GraphAPIErrorResponse
type GraphAPIErrorResponse struct {
	Error GraphAPIErrorObject `json:"error"`
}

//GraphAPIErrorObject GraphAPIErrorObject
type GraphAPIErrorObject struct {
	Code       string                   `json:"code"`
	Message    string                   `json:"message"`
	InnerError GraphAPIInnerErrorObject `json:"innerError"`
}

//GraphAPIInnerErrorObject GraphAPIInnerErrorObject
type GraphAPIInnerErrorObject struct {
	Date            string `json:"date"`
	RequestID       string `json:"request-id"`
	ClientRequestID string `json:"client-request-id"`
}

// AzureGraphAPIURL the graph API endpoint
const AzureGraphAPIURL = "https://graph.microsoft.com/"

func (b BaseResource) List(r ResourceAPI, context cli.Context, args cli.Args) []Resource {
	path := r.CreateRequestPath(context, args)
	params := r.CreateQueryParams(context, args)

	request := b.newRequest("GET", path, params, nil)
	body := b.do(request)

	return r.ConvertToResourceSlice(body)
}

func (b BaseResource) newRequest(method, path string, queryParams url.Values, body interface{}) *http.Request {

	rel := &url.URL{Path: path, RawQuery: queryParams.Encode()}
	baseURL, _ := url.Parse(AzureGraphAPIURL)
	u := baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		helpers.ErrorHandlerFatal("JSON encoding for request body failed:", err)
	}

	req, err := http.NewRequest(method, u.String(), buf)
	helpers.ErrorHandlerFatal("Request construction failed:", err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", b.UserAgent)

	return req
}

func (b BaseResource) do(req *http.Request) []byte {

	if log.GetLevel() == log.TraceLevel {
		dump, err := httputil.DumpRequestOut(req, true)
		helpers.ErrorHandlerFatal("Request dump failed:", err)
		log.Tracef("REQUEST: %s", string(dump))
	}

	resp, err := b.HTTPClient.Do(req)
	helpers.ErrorHandlerFatal("Request execution failed:", err)

	if log.GetLevel() == log.TraceLevel {
		dump, err := httputil.DumpResponse(resp, true)
		helpers.ErrorHandlerFatal("Request dump failed:", err)
		log.Tracef("RESPONSE: %q", dump)
	}

	log.Trace("Status Response:", resp.Status)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	helpers.ErrorHandlerFatal("Reading response body failed", err)
	log.Trace(string(body))

	if resp.StatusCode != 200 {
		var graphErr GraphAPIErrorResponse
		err = json.Unmarshal(body, &graphErr)
		helpers.ErrorHandlerFatal("JSON unmarshalling of response body failed:", err)

		log.Fatalf("Error from GraphAPI: %+v", graphErr.Error)
	}
	return body
}

func CreateURLFilterParams(criteria *Criteria) url.Values {
	params := url.Values{}
	params.Add("$filter", (*criteria).String())
	return params
}
