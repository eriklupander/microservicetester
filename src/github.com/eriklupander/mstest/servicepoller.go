package main

import (
        "sync"
        "net/http"
        "crypto/tls"
        "time"
)

func PollService(service string, wg *sync.WaitGroup, total int, consoleRow int) {

        Cprint(consoleRow, 0, service)
        Cprint(consoleRow, 60, "... waiting ")

        req := BuildHttpRequest(service, "GET")
        var DefaultTransport http.RoundTripper = &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        for {
                resp, err := DefaultTransport.RoundTrip(req)
                if err != nil || resp.StatusCode > 299 {
                        Cprint(consoleRow, 0, service)
                        Cprint(consoleRow, 60, "... waiting ")
                }  else {
                        Cprint(consoleRow, 0, service)
                        Cprint(consoleRow, 60, "done                   ")
                        wg.Done()
                        return
                }
                time.Sleep(time.Second * 1)
        }
}




