package curl

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idoubi/goz"
	"github.com/tidwall/gjson"

	"github.com/dekaiju/go-skeleton/pkg/log"
)

func Get(c *gin.Context, url string, data interface{}, headers map[string]interface{}) (gjson.Result, error) {
	var res gjson.Result
	cli := goz.NewClient()
	resp, err := cli.Get(url, goz.Options{
		Headers: headers,
		Query:   data,
	})
	if err != nil {
		log.WithGinContext(c).WithError(err).Error("http get request failed")
		return res, err
	}
	body, err := resp.GetBody()
	if err != nil {
		log.WithGinContext(c).WithError(err).Error("http get response read failed")
		return res, err
	}
	res = gjson.Parse(body.GetContents())
	return res, nil
}

func PostForm(c *gin.Context, url string, data map[string]interface{}, headers map[string]interface{}) (gjson.Result, error) {
	var res gjson.Result
	cli := goz.NewClient()
	resp, err := cli.Post(url, goz.Options{
		Headers:    headers,
		FormParams: data,
	})
	if err != nil {
		log.WithGinContext(c).WithError(err).Error("http post form request failed")
		return res, err
	}
	body, err := resp.GetBody()
	if err != nil {
		log.WithGinContext(c).WithError(err).Error("http post form response read failed")
		return res, err
	}
	res = gjson.Parse(body.GetContents())
	return res, nil
}

func PostJson(c *gin.Context, url string, data string, headers map[string]string) (gjson.Result, error) {
	var res gjson.Result
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	for index := range headers {
		req.Header.Set(index, headers[index])
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithGinContext(c).WithError(err).Error("http post json request failed")
		return res, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithGinContext(c).WithError(err).Error("http post json response read failed")
		return res, err
	}
	res = gjson.Parse(string(content))
	return res, nil
}
