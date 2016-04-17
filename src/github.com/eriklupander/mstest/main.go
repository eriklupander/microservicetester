package main

import (
        "fmt"
        "os/exec"
        "time"
       // "os"
        "net/http"
        "crypto/tls"
        "log"
        //"io/ioutil"
        "sync"

"strings"
        "io/ioutil"
        "encoding/base64"
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
        DockerComposeUp()
        consoleRow++

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



        // 5. When all are started, execute list of HTTP calls.
        wg.Wait()

        // Fix cursor position
        consoleRow+=3
        fmt.Printf("\033[%d;0H", consoleRow)   // Move cursor to row

        // Fix OAUTH token -d grant_type=password -d client_id=acme -d scope=webshop -d username=user -d password=password

        body := "grant_type=" + t.OAuth.Grant_type + "&client_id=" + t.OAuth.Client_id + "&scope=" + t.OAuth.Scope + "&username=" + t.OAuth.Username + "&password=" + t.OAuth.Password
        reader := strings.NewReader(body)

        //fmt.Println("OAuth data: " + t.OAuth.Url + " data: " + body)
        postReq, err := http.NewRequest("POST", t.OAuth.Url, reader)
        postReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        //b64 := basicAuth(t.OAuth.Client_id, t.OAuth.Client_password)
        //fmt.Println("Basic: " + b64)
        //postReq.Header.Add("Authorization", "Basic " + b64)
        postReq.SetBasicAuth(t.OAuth.Client_id, t.OAuth.Client_password)
        if err != nil {
            log.Fatal("Error constructing OAuth POST")
        }
        var DefaultTransport http.RoundTripper = &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        resp, err := DefaultTransport.RoundTrip(postReq)
        if err != nil {
                fmt.Println(resp.Body)
                fmt.Println(err)
                log.Fatalf("OAuth request failed")
        }  else {
              respBody, _ := ioutil.ReadAll(resp.Body)
              fmt.Println(string(respBody))
        }

        // 6. Shut down
        time.Sleep(time.Second * 1)
}

func basicAuth(username string, password string) string {
        auth := username + ":" + password
        return base64.StdEncoding.EncodeToString([]byte(auth))
}


var l sync.Mutex

func cprint(row int, col int, text string) {
       l.Lock()
                fmt.Printf("\033[%d;%dH", row, col)
                fmt.Print(text)
       l.Unlock()
}

func pollService(service string, wg *sync.WaitGroup, total int, consoleRow int) {

        cprint(consoleRow, 0, service)
        cprint(consoleRow, 60, "... waiting ")

        req := buildHttpRequest(service)
        var DefaultTransport http.RoundTripper = &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        for {
                resp, err := DefaultTransport.RoundTrip(req)
                if err != nil || resp.StatusCode > 299 {
                        cprint(consoleRow, 0, service)
                        cprint(consoleRow, 60, "... waiting ")
                }  else {
                        cprint(consoleRow, 0, service)
                        cprint(consoleRow, 60, "done                   ")
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

func DockerComposeUp() {
        cmd := exec.Command("docker-compose", "up", "-d")
        // cmd.Stdout = os.Stdout
        // cmd.Stderr = os.Stderr
        cmd.Run()
        fmt.Println("Docker starting up...")
}

