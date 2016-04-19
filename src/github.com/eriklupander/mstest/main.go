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
        "strconv"
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
        defer DockerComposeDown(t)
        consoleRow+=2

        // 4. Then wait for specified microservices
        awaitServicesHasStarted(t)

        // 5. When all are started, store OAuth token
        consoleRow++
        Cprint(consoleRow, 0, "Getting OAuth token ...")
        StoreOAuthToken(t)
        Cprint(consoleRow, 0, "Getting OAuth token ... OK")

        // 6. execute list of endpoint HTTP calls.
        runEndpoints(t)

        // 6. Shut down
        time.Sleep(time.Second * 3)

}

func awaitServicesHasStarted(t TestDef) {
        wg := sync.WaitGroup{}

        fmt.Println("Waiting for all microservices to start...")
        consoleRow++
        fmt.Println("")
        consoleRow++

        for _, service := range t.Services {
                wg.Add(1)
                consoleRow++
                go PollService(service, &wg, len(t.Services), consoleRow)
        }
        wg.Wait()

        // Fix cursor position after all services have started
        consoleRow+=3
        fmt.Printf("\033[%d;0H", consoleRow)   // Move cursor to row
}

func runEndpoints(t TestDef) {
        consoleRow+=3
        fmt.Printf("\033[%d;0H", consoleRow)   // Move cursor to row
        wg2 := sync.WaitGroup{}
        for _, endpoint := range t.Endpoints {
                consoleRow++
                wg2.Add(1)
                time.Sleep(time.Millisecond * 300) // Stagger service calls slightly.
                go invokeEndpoint(endpoint, &wg2, consoleRow)
        }
        wg2.Wait()
        consoleRow+=2
        Cprint(consoleRow, 0, "All done.\n")
}

func invokeEndpoint(endpoint Endpoint, wg2 *sync.WaitGroup, consoleRow int) {

        req := BuildHttpRequest(endpoint.Url, endpoint.Method)
        var DefaultTransport http.RoundTripper = &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        if endpoint.Auth_method == "TOKEN" {
                req.Header.Add("Authorization", "Bearer " + TOKEN)
        } else if endpoint.Auth_method == "BASIC" {
                log.Fatalln("BASIC auth not supported yet for endpoint invocations.")
        }
        retries := 10
        for i := 0; i < retries; i++ {
                if i == 0 {
                        Cprint(consoleRow, 0, endpoint.Url)
                        Cprint(consoleRow, 70, "... testing                          ")
                }

                resp, err := DefaultTransport.RoundTrip(req)

                if err != nil {
                       // fmt.Println(endpoint.Url + " failed with error: " + err.Error())
                        Cprint(consoleRow, 0, endpoint.Url)
                        Cprint(consoleRow, 70, "... failed. Retrying " + strconv.Itoa(i) + "/" + strconv.Itoa(retries) + " ...            ")
                } else if resp.StatusCode > 299 {
                       // fmt.Println(endpoint.Url + " failed with status: " + string(resp.StatusCode) + " " + resp.Status)
                        Cprint(consoleRow, 0, endpoint.Url)
                        Cprint(consoleRow, 70, "... failed. Retrying " + strconv.Itoa(i) + "/" + strconv.Itoa(retries) + " ...             ")
                } else {
                       // fmt.Println(endpoint.Url + " " + string(resp.StatusCode) + " " + resp.Status)
                        Cprint(consoleRow, 0, endpoint.Url)
                        Cprint(consoleRow, 70, "... OK                                          ")
                        wg2.Done()
                        break
                }
                time.Sleep(time.Second * 3)

                if (i == retries - 1) {
                        wg2.Done()
                        Cprint(consoleRow, 0, endpoint.Url)
                        Cprint(consoleRow, 70, "... All attempts failed, something is broken.                                 ")
                }
        }

}


