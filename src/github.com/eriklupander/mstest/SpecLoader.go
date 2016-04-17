package main

import (
        "fmt"
        "gopkg.in/yaml.v2"
        "os"
        "io/ioutil"
)

func LoadSpecification() (TestDef) {
        spec := parseSpecFile()
        dir, _ := os.Getwd()
        dat, _ := ioutil.ReadFile(dir + "/" + spec)
        var t TestDef
        yaml.Unmarshal([]byte(dat), &t)

        fmt.Println("Loaded specification '" + t.Title + "'")

        return t
}

func parseSpecFile() (string) {
        if len(os.Args) == 1 {
                fmt.Errorf("No command line arguments, exiting...\n")
                panic("Cannot start, no YAML specification supplied as command-line argument")
        }
        var s, sep string
        for i := 1; i < len(os.Args); i++ {
                s += sep + os.Args[i]
                sep = " "
        }
        if s == "" {
                panic(fmt.Sprintf("Specified simulation file '%s' is not a .yml file", s))
        }
        return s
}
