package main

import (
        "fmt"
        //"os/exec"
        "time"
       // "os"
        "net/http"
        "crypto/tls"
        "log"
        //"io/ioutil"
        "sync"
       // "os"
)

var consoleRow = 0

func main() {
        CallClear()
        fmt.Println("Starting up...")
        consoleRow++

        // 1. Check required command-line tools are present
        // docker-compose etc.
        CheckDocker()
        consoleRow++

        // 2. Load service specification from yaml
        t := LoadSpecification()
        consoleRow++

        // 3. Start using docker-compose up -d
        DockerComposeUp(t)
        consoleRow+=2

        // 4. Then wait for specified microservices

        wg := sync.WaitGroup{}

        fmt.Println("Waiting for all microservices to start...")
        consoleRow++
        fmt.Println("")
        consoleRow++

        for _, service := range t.Services {
                wg.Add(1)
                consoleRow++
                go pollService(service, &wg, len(t.Services), consoleRow)
        }

        wg.Wait()

        // Fix cursor position after all services have started
        consoleRow+=3
        fmt.Printf("\033[%d;0H", consoleRow)   // Move cursor to row

        // 5. When all are started, store OAuth token
        StoreOAuthToken(t)

        // 6. execute list of endpoint HTTP calls.
        for _, endpoint := range t.Endpoints {
                invokeEndpoint(endpoint)
        }

        // 6. Shut down
        time.Sleep(time.Second * 1)
        DockerComposeDown(t)
}

func invokeEndpoint(endpoint Endpoint) {
        req := buildHttpRequest(endpoint.Url)
        var DefaultTransport http.RoundTripper = &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        if endpoint.Auth_method == "TOKEN" {
                fmt.Println("Setting auth token: " + TOKEN)
                req.Header.Add("Authorization", "Bearer " + TOKEN)
        } else if endpoint.Auth_method == "NONE" {
                // Why??
        }

        resp, err := DefaultTransport.RoundTrip(req)
        if err != nil {
                log.Fatalln(endpoint.Url + " failed with error: " + err.Error())
        } else if resp.StatusCode > 299 {
                log.Fatalln(endpoint.Url + " failed with status: " + string(resp.StatusCode) + " " + resp.Status)
        }
        fmt.Println(endpoint.Url + " " + string(resp.StatusCode) + " " + resp.Status)
}



func pollService(service string, wg *sync.WaitGroup, total int, consoleRow int) {

        Cprint(consoleRow, 0, service)
        Cprint(consoleRow, 60, "... waiting ")

        req := buildHttpRequest(service)
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

func buildHttpRequest(url string) *http.Request {

        var req *http.Request
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
                log.Fatal(err)
        }
        return req
}


