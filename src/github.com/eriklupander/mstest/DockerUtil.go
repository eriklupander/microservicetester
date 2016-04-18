package main

import (
        "os/exec"
        "fmt"
        "log"
        "os"
)

func CheckDocker() {
        _, err := exec.LookPath("docker-compose")
        if err != nil {
                log.Fatal("docker-compose not installed, fix!")
        }
        fmt.Printf("docker-compose installed OK\n")
}

func DockerComposeUp(t TestDef) {
        cmd := exec.Command("docker-compose", "-f", t.DockerComposeFile, "up", "-d")
        cmd.Dir = t.DockerComposeRoot
        env := os.Environ()
        env = append(env, fmt.Sprintf("PROJECT_ROOT=%s", t.DockerComposeRoot))
        cmd.Env = env
        //cmd.Stdout = os.Stdout
        //cmd.Stderr = os.Stderr
        cmd.Run()
        fmt.Println("Docker starting up using " + t.DockerComposeRoot + "/" + t.DockerComposeFile + " ...")

}

func DockerComposeDown(t TestDef) {
        fmt.Println("Docker shutting down...")
        cmd := exec.Command("docker-compose", "-f", t.DockerComposeRoot + "/" + t.DockerComposeFile, "down")
        // cmd.Stdout = os.Stdout
        // cmd.Stderr = os.Stderr
        cmd.Run()
}




