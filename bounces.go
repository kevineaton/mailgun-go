package mailgun

import (
	"strconv"
	"time"
	"github.com/mbanzon/simplehttp"
)

type Bounce struct {
	CreatedAt string `json:"created_at"`
	Code      string `json:"code"`
	Address   string `json:"address"`
	Error     string `json:"error"`
}

type bounceEnvelope struct {
	TotalCount int      `json:"total_count"`
	Items      []Bounce `json:"items"`
}

type singleBounceEnvelope struct {
	Bounce Bounce `json:"bounce"`
}

func (i Bounce) GetCreatedAt() (t time.Time, err error) {
	t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 MST", i.CreatedAt)
	return
}

func (m *mailgunImpl) GetBounces(limit, skip int) (int, []Bounce, error) {
	r := simplehttp.NewSimpleHTTPRequest("GET", generateApiUrl(m, bouncesEndpoint))
	if limit != -1 {
		r.AddParameter("limit", strconv.Itoa(limit))
	}
	if skip != -1 {
		r.AddParameter("skip", strconv.Itoa(skip))
	}

	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	var response bounceEnvelope
	err := r.MakeJSONRequest(&response)
	if err != nil {
		return -1, nil, err
	}

	return response.TotalCount, response.Items, nil
}

func (m *mailgunImpl) GetSingleBounce(address string) (Bounce, error) {
	r := simplehttp.NewSimpleHTTPRequest("GET", generateApiUrl(m, bouncesEndpoint) + "/" + address)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	var response singleBounceEnvelope
	err := r.MakeJSONRequest(&response)
	if err != nil {
		return Bounce{}, err
	}

	return response.Bounce, nil
}

func (m *mailgunImpl) AddBounce(address, code, error string) error {
	r := simplehttp.NewSimpleHTTPRequest("POST", generateApiUrl(m, bouncesEndpoint))

	r.AddFormValue("address", address)
	if code != "" {
		r.AddFormValue("code", code)
	}
	if error != "" {
		r.AddFormValue("error", error)
	}
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	_, err := r.MakeRequest()
	return err
}

func (m *mailgunImpl) DeleteBounce(address string) error {
	r := simplehttp.NewSimpleHTTPRequest("DELETE", generateApiUrl(m, bouncesEndpoint) + "/" + address)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	_, err := r.MakeRequest()
	return err
}
