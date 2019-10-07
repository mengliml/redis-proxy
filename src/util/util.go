package util

import (
    "net"
    "net/http"
    "io/ioutil"
    "testing"
)

// Make request to the sever and return response body and status code
func GetValueFromServer(path string, s *http.Server, t *testing.T) (string, int) {
    ln, err := net.Listen("tcp", ":0")
    if err != nil {
        t.Fatal(err)
    }
    defer ln.Close()
    go s.Serve(ln) 
    
    res, err := http.Get("http://" + ln.Addr().String() + path)
    if err != nil {
        t.Fatal(err)
    }

    body, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        t.Fatal(err)
    }

    return string(body), res.StatusCode
}