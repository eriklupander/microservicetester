package main

import (
        "os/exec"
        "log"
        "fmt"
)

func CheckDocker() {
        _, err := exec.LookPath("docker-compose")
        if err != nil {
                log.Fatal("docker-compose not installed, fix!")
        }
        fmt.Printf("docker-compose installed OK\n")
}
