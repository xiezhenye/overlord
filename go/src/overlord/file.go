package overlord
import (
	"io"
	"io/ioutil"
	"os"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"bytes"
)

func (self *Client) GetFileReader(group, item, tail string) (io.Reader, error) {
	resp, err := self.Request("GET", "files", group, item, tail, nil, nil, nil)
	if resp != nil {
		return resp.Body, err
	} else {
		return nil, err
	}
}

func (self *Client) GetFileContent(group, item, tail string) ([]byte, error) {
	reader, err := self.GetFileReader(group, item, tail)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func (self *Client) GetFileTo(group, item, tail string, localPath string) error {
	return getRequestToFile(func() (io.Reader, error) {
		return self.GetFileReader(group, item, tail)
	}, localPath)
}

func (self *Client) GetFileRangeReader(group, item, tail string, start, length int64) (io.Reader, error) {
	headers := make(map[string]string)
	headers["Range"] = rangeHeader(start, length)
	resp, err := self.Request("GET", "files", group, item, tail, nil, nil, headers)
	if resp != nil {
		return resp.Body, err
	} else {
		return nil, err
	}
}

func (self *Client) GetFileRangeContent(group, item, tail string, start, length int64) ([]byte, error) {
	reader, err := self.GetFileRangeReader(group, item, tail, start, length)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func (self *Client) GetFileRangeTo(group, item, tail string, start, length int64, localPath string) error {
	return getRequestToFile(func() (io.Reader, error) {
		return self.GetFileRangeReader(group, item, tail, start, length)
	}, localPath)
}

func (self *Client) PostFileReader(group, item, tail string, reader io.Reader) error {
	_, err := self.Request("POST", "files", group, item, tail, nil, reader, nil)
	return err
}

func (self *Client) PostFileContent(group, item, tail string, content []byte) error {
	reader := bytes.NewReader(content)
	return self.PostFileReader(group, item, tail, reader)
}

func (self *Client) PostFileFrom(group, item, tail string, path string) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	return self.PostFileReader(group, item, tail, reader)
}

func (self *Client) PutFileReader(group, item, tail string, reader io.Reader) error {
	_, err := self.Request("PUT", "files", group, item, tail, nil, reader, nil)
	return err
}

func (self *Client) PutFileContent(group, item, tail string, content []byte) error {
	reader := bytes.NewReader(content)
	return self.PutFileReader(group, item, tail, reader)
}

func (self *Client) PutFileFrom(group, item, tail string, path string) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	return self.PutFileReader(group, item, tail, reader)
}

func (self *Client) DeleteFile(group, item, tail string) error {
	_, err := self.Request("DELETE", "files", group, item, tail, nil, nil, nil)
	return err

}


type FileInfo struct {
	Size   int64
	Mode   os.FileMode // uint32
	Mtime  time.Time
}

func (self *Client) GetFileInfo(group, item, tail string) (ret *FileInfo, err error) {
	resp, err := self.Request("HEAD", "files", group, item, tail, nil, nil, nil)
	if resp.StatusCode != http.StatusOK {
		err = RespToError(resp)
		return
	}
	ret = &FileInfo{}
	ret.Size, err = strconv.ParseInt(resp.Header.Get("X-Servant-File-Size"), 10, 64)
	if err != nil {
		return
	}
	format := "2006-01-02 15:04:05.999999999 -0700 MST" // see time.Time.String()
	ret.Mtime, err = time.Parse(format, resp.Header.Get("X-Servant-File-Mtime"))
	if err != nil {
		return
	}
	m, err := strconv.ParseUint(resp.Header.Get("X-Servant-File-Mode"), 10, 32)
	if err != nil {
		return
	}
	ret.Mode = os.FileMode(m)
	return
}


func rangeHeader(start, length int64) string {
	return fmt.Sprintf("bytes=%d-%d", start, start + length - 1)
}

func getRequestToFile(f func()(io.Reader, error), path string) error {
	dst, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	reader, err := f()
	if err != nil {
		dst.Close()
		return err
	}
	_, err = io.Copy(dst, reader)
	dst.Close()
	return err
}


