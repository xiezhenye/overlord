package overlord
import (
	"encoding/json"
	"io/ioutil"
)

type SqlResult []map[string]string

func (self *Client) Sql(group, item string) ([]SqlResult, error) {
	return self.SqlWithParams(group, item, nil)
}

func (self *Client) SqlWithParams(group, item string, params map[string]string) ([]SqlResult, error) {
	resp, err := self.Request("GET", "database", group, item, "", params, nil, nil)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ret []SqlResult
	err = json.Unmarshal(content, ret)
	return ret, err
}
