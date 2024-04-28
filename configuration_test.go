package gowfs

import "testing"

func Test_GetNameNodeUrl(t *testing.T) {
    conf := Configuration{Addr: "localhost:8080", BasePath: "/test/gofs", User: "vvivien", PassWord: "123456"}
    u, err := conf.GetNameNodeUrl()

    if err != nil {
        t.Fatal(err)
    }

    if u.Scheme != "http" {
        t.Errorf("Expecting url.Scheme http, but got %s", u.Scheme)
    }

    if u.Host != "localhost:8080" {
        t.Errorf("Expecting url.Host locahost:8080, but got %s", u.Host)
    }

    if u.Path != conf.WebHdfsVer+conf.BasePath {
        t.Errorf("Expecting url.Path %s, but got %s", conf.WebHdfsVer+conf.BasePath, u.Path)
    }
}

func TestConfiguration_GetBasicAuthInfoHeader(t *testing.T) {
    conf := Configuration{Addr: "localhost:8080", BasePath: "/test/gofs", User: "vvivien", PassWord: "123456"}
    header, err := conf.GetBasicAuthInfoHeader()
    if err != nil {
        t.Fatal(err)
    }
    t.Logf("header:%s", header)
}
