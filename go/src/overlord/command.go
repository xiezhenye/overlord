package overlord

import (
	"io"
	"io/ioutil"
)

func (self *Client) RunCommand(group, name string) ([]byte, error) {
	return self.RunCommandWithInputAndParams(group, name, nil, nil)
}

func (self *Client) RunCommandWithInput(group, name string, input io.Reader) ([]byte, error) {
	return self.RunCommandWithInputAndParams(group, name, input, nil)
}

func (self *Client) RunCommandWithParams(group, name string, params map[string]string) ([]byte, error) {
	return self.RunCommandWithInputAndParams(group, name, nil, params)
}

func (self *Client) RunCommandWithInputAndParams(group, name string, input io.Reader, params map[string]string) ([]byte, error) {
	method := "GET"
	if input != nil {
		method = "POST"
	}
	resp, err := self.Request(method, "commands", group, name, "", params, input, nil)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}
