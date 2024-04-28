package gowfs

import (
    "encoding/base64"
    "fmt"
    "net/http"
)
import "errors"
import "time"
import "net/url"

type Configuration struct {
    Addr                  string // host:port
    BasePath              string // initial base path to be appended
    User                  string // user.name to use to connect
    PassWord              string // password to connect
    WebHdfsVer            string // webhdfs version
    ConnectionTimeout     time.Duration
    DisableKeepAlives     bool
    DisableCompression    bool
    ResponseHeaderTimeout time.Duration
    MaxIdleConnsPerHost   int
}

func NewConfiguration() *Configuration {
    return &Configuration{
        ConnectionTimeout:     time.Second * 17,
        DisableKeepAlives:     false,
        DisableCompression:    true,
        ResponseHeaderTimeout: time.Second * 17,
    }
}

func (conf *Configuration) GetNameNodeUrl() (*url.URL, error) {
    if &conf.Addr == nil {
        return nil, errors.New("Configuration namenode address not set.")
    }

    var urlStr string = fmt.Sprintf("https://%s%s%s", conf.Addr, conf.WebHdfsVer, conf.BasePath)

    u, err := url.Parse(urlStr)
    if err != nil {
        return nil, err
    }

    return u, nil
}

func (conf *Configuration) GetBasicAuthInfoHeader() (*http.Header, error) {
    if len(conf.User) == 0 {
        return nil, errors.New("user should not be empty")
    }
    if len(conf.PassWord) == 0 {
        return nil, errors.New("password should not be empty")
    }
    auth := fmt.Sprintf("%s:%s", conf.User, conf.PassWord)
    basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
    header := make(http.Header)
    header.Set("Authorization", basicAuth)

    return &header, nil
}
