package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

var clientLogger = log.WithField("go", "clients/client.go")

const BODY_JSON = "json"
const BODY_FORM = "form"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ContextKey int

const RequestID ContextKey = 1

// Client Fields
type Client struct {
	BaseURI         string
	Action          string
	Version         string
	Headers         map[string]string
	Method          string
	APIKey          string
	ParamType       string
	Params          map[string]interface{}
	ParamsInterface interface{}
	HTTPClient      HTTPClient
	Username        string
	Password        string
}

// SetBaseURI set BaseURI field
func (c *Client) SetBaseURI(baseURI string) {
	c.BaseURI = baseURI
}

// SetAction set Action field
func (c *Client) SetAction(action string) {
	c.Action = action
}

// SetVersion set Version field
func (c *Client) SetVersion(version string) {
	c.Version = version
}

// SetHeaders set Headers field
func (c *Client) SetHeaders(headers map[string]string) {
	c.Headers = headers
}

// SetMethod set method field
func (c *Client) SetMethod(method string) {
	c.Method = method
}

// SetAPIKey set APIKey field
func (c *Client) SetAPIKey(apiKey string) {
	c.APIKey = apiKey
}

// SetParamType set ParamType field
func (c *Client) SetParamType(paramType string) {
	c.ParamType = paramType
}

// SetParams set Params field
func (c *Client) SetParams(params map[string]interface{}) {
	c.Params = params
}

// SetParamsInterface set Params interface field with parameter interface
func (c *Client) SetParamsInterface(paramsInterface interface{}) {
	c.ParamsInterface = paramsInterface
}

// SetUsername set Username field
func (c *Client) SetUsername(username string) {
	c.Username = username
}

// SetPassword set Password field
func (c *Client) SetPassword(password string) {
	c.Password = password
}

// CreateWithContext create http request with context
func (c *Client) CreateWithContext(ctx context.Context) (int, map[string]interface{}, error) {
	fLog := clientLogger.WithField("func", "CreateWithContext").WithField("RequestId", ctx.Value(RequestID))

	statusCode := 500

	var payload *bytes.Buffer = new(bytes.Buffer)

	if c.Method != "GET" {
		if c.Params != nil && c.ParamType == BODY_JSON {
			err := json.NewEncoder(payload).Encode(c.Params)
			if err != nil {
				fLog.Errorf("Encode got: %v", err.Error())
				return statusCode, nil, err
			}
		}

		if c.Params != nil && c.ParamType == BODY_FORM {
			formData := url.Values{}
			for key, value := range c.Params {
				formData.Add(key, value.(string))
			}
			payload = bytes.NewBufferString(formData.Encode())
		}

		if c.ParamsInterface != nil {
			err := json.NewEncoder(payload).Encode(c.ParamsInterface)
			if err != nil {
				fLog.Errorf("Encode got: %v", err.Error())
				return statusCode, nil, err
			}
		}
	}

	URIAction := ""

	if c.Version != "" {
		URIAction = c.BaseURI + "/" + "v" + c.Version + "/" + c.Action
	} else {
		URIAction = c.BaseURI + "/" + c.Action
	}

	request, err := http.NewRequest(c.Method, URIAction, payload)

	if c.Method == "GET" {
		if c.Params != nil {
			q := request.URL.Query()
			for key, value := range c.Params {
				q.Add(key, value.(string))
			}
			request.URL.RawQuery = q.Encode()
		}
	}

	if err != nil {
		fLog.Errorf("NewRequest got: %v", err.Error())
		return statusCode, nil, err
	}

	// set header params
	request.Header.Set("x-api-key", c.APIKey)
	request.Header.Set("accept", "application/json")
	for key, value := range c.Headers {
		request.Header.Set(key, value)
	}

	if c.Username != "" {
		request.SetBasicAuth(c.Username, c.Password)
	}

	response, err := c.HTTPClient.Do(request.WithContext(ctx))
	fLog.Trace(URIAction)

	if err != nil {
		fLog.Errorf("HTTPClient.Do got: %v", err.Error())
		return statusCode, nil, errors.New("error dial another platform")
	}
	defer response.Body.Close()

	statusCode = response.StatusCode

	bodyResponse, err := io.ReadAll(response.Body)
	if err != nil {
		fLog.Errorf("ioutil.ReadAll got: %v", err.Error())
		return statusCode, nil, errors.New("error read body from another platform")
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyResponse, &result)

	if err != nil {
		fLog.Errorf("Unmarshal got: %v", err.Error())
		return statusCode, nil, errors.New("error Unmarshal Json")
	}

	return statusCode, result, nil
}

// CreateMultiPartWithContext create multipart request with context
func (c *Client) CreateMultiPartWithContext(ctx context.Context) (int, map[string]interface{}, error) {
	fLog := clientLogger.WithField("func", "CreateMultiPartWithContext").WithField("RequestId", ctx.Value(RequestID))

	statusCode := http.StatusInternalServerError

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	var err error

	for key, param := range c.Params {
		var fw io.Writer

		if mFileHeader, ok := param.(*multipart.FileHeader); ok {
			if mFileHeader == nil {
				continue
			}
			fw, err = createFormFile(writer, key, mFileHeader.Filename, mFileHeader.Header["Content-Type"][0])

			if err != nil {
				return statusCode, nil, err
			}

			mFile, err := mFileHeader.Open()

			if err != nil {
				return statusCode, nil, err
			}
			mFile.Close()

			_, err = io.Copy(fw, mFile)

			if err != nil {
				return statusCode, nil, err
			}

		} else {
			if err = writer.WriteField(key, param.(string)); err != nil {
				return statusCode, nil, err
			}
		}

	}

	err = writer.Close()
	if err != nil {
		return statusCode, nil, err
	}

	URIAction := ""

	if c.Version != "" {
		URIAction = c.BaseURI + "/" + "v" + c.Version + "/" + c.Action
	} else {
		URIAction = c.BaseURI + "/" + c.Action
	}

	request, err := http.NewRequest(c.Method, URIAction, body)

	if err != nil {
		fLog.Errorf("NewRequest got: %v", err.Error())
		return statusCode, nil, err
	}

	// set header params
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("x-api-key", c.APIKey)
	request.Header.Set("accept", "application/json")

	for key, value := range c.Headers {
		request.Header.Set(key, value)
	}

	response, err := c.HTTPClient.Do(request.WithContext(ctx))
	fLog.Trace(URIAction)

	if err != nil {
		fLog.Errorf("HTTPClient.Do got: %v", err.Error())
		return statusCode, nil, errors.New("error dial another platform")
	}
	defer response.Body.Close()

	statusCode = response.StatusCode

	bodyResponse, err := io.ReadAll(response.Body)
	if err != nil {
		fLog.Errorf("ioutil.ReadAll got: %v", err.Error())
		return statusCode, nil, errors.New("error read body from another platform")
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyResponse, &result)

	if err != nil {
		fLog.Errorf("Unmarshall got: %v", err.Error())
		return statusCode, nil, errors.New("error nnmarshal json")
	}

	return statusCode, result, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createFormFile(w *multipart.Writer, fieldname, filename, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}
