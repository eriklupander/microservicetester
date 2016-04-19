package main

import (
        "net/http"
        "log"
)

func BuildHttpRequest(url string, method string) *http.Request {

        var req *http.Request
        req, err := http.NewRequest(method, url, nil)
        if err != nil {
                log.Fatal(err)
        }
        return req
}
