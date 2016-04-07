package overlord

import (
	"net/http"
	"crypto/sha1"
	"encoding/hex"
	"time"
	"strconv"
	"fmt"
	"io"
	"net/url"
)

const ServantErrHeader = "X-Servant-Err"

type Client struct {
	Host  string
	User  string
	Key   string
}

type Error struct {
	Code     int
	Message  string
}


func (self Error) Error() string {
	return fmt.Sprintf("%d: %s", self.Code, self.Message)
}

func RespToError(resp *http.Response) Error {
	msg := resp.Header.Get(ServantErrHeader)
	return Error{Code:resp.StatusCode, Message: msg}
}

func (self *Client) authHeader(method, uri string) string {
	// Authorization: user ts sha1(user + key + ts + method + uri)
	tsStr := strconv.FormatInt(time.Now().Unix(), 10)
	strToHash := self.User + self.Key + tsStr + method + uri
	sha1Sum := sha1.Sum([]byte(strToHash))
	hexHash := hex.EncodeToString(sha1Sum[:])
	return fmt.Sprintf("%s %s %s", self.User, tsStr, hexHash)
}

func (self *Client) httpRequest(method, uri string, body io.Reader) (*http.Request, error) {
	ret, err:= http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	ret.Header.Set("Authorization", self.authHeader(method, uri))
	return ret, nil
}


func (self *Client) buildUri(resource, group, item, tail string, params map[string]string) string {
	uri := fmt.Sprintf("/%s/%s/%s%s", url.QueryEscape(resource), url.QueryEscape(group), url.QueryEscape(item), tail)
	values := url.Values{}
	if params != nil && len(params) > 0 {
		for k, v := range (params) {
			values.Add(k, v)
		}
		uri += "?" + values.Encode()
	}
	return uri
}

func (self *Client) Request(method, resource, group, item, tail string,
	params map[string]string, body io.Reader, headers map[string]string) (ret *http.Response, err error) {
	uri := self.buildUri(resource, group, item, tail, params)
	req, err := self.httpRequest(method, uri, body)
	if err != nil {
		return
	}
	client := http.Client{}
	if headers != nil {
		for k, v := range (headers) {
			req.Header.Add(k, v)
		}
	}
	ret, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if ret.StatusCode != http.StatusOK {
		err = RespToError(ret)
	}
	return ret, err
}
