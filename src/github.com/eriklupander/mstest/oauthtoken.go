package main

import (
        "strings"
        "net/http"
        "log"
        "crypto/tls"
        "fmt"
        "io/ioutil"
        "encoding/json"
)

var TOKEN string

func StoreOAuthToken(t TestDef) {
        body := "grant_type=" + t.OAuth.Grant_type + "&client_id=" + t.OAuth.Client_id + "&scope=" + t.OAuth.Scope + "&username=" + t.OAuth.Username + "&password=" + t.OAuth.Password
        reader := strings.NewReader(body)
        postReq, err := http.NewRequest("POST", t.OAuth.Url, reader)
        postReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
                log.Fatalf("OAuth request failed")
        }  else {
                respBody, _ := ioutil.ReadAll(resp.Body)
                fmt.Println(string(respBody))
                var auth interface{}
                err = json.Unmarshal(respBody, &auth)
                m := auth.(map[string]interface{})
                fmt.Println("Token: " + m[t.OAuth.Token_key].(string))
                TOKEN = m[t.OAuth.Token_key].(string)
                if err != nil {
                        fmt.Println("Error unmarshalling OAuth Token JSON")
                }
        }
}
