package main

import (
        //"fmt"
        "os"
        "os/exec"
        "runtime"
        //"time"
        "sync"
        "fmt"
)

var clear map[string]func() //create a map for storing clear funcs

func init() {
        clear = make(map[string]func()) //Initialize it
        clear["linux"] = func() {
                cmd := exec.Command("clear") //Linux example, its tested
                cmd.Stdout = os.Stdout
                cmd.Run()
        }
        clear["darwin"] = func() {
                cmd := exec.Command("clear") //darwin example, its tested
                cmd.Stdout = os.Stdout
                cmd.Run()
        }
        clear["windows"] = func() {
                cmd := exec.Command("cls") //Windows example it is untested, but I think its working
                cmd.Stdout = os.Stdout
                cmd.Run()
        }
}

func CallClear() {
        value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
        if ok { //if we defined a clear func for that platform:
                value()  //we execute it
        } else { //unsupported platform
                panic("Your platform is unsupported! I can't clear terminal screen :(")
        }
}

var l sync.Mutex

func Cprint(row int, col int, text string) {
        l.Lock()
        fmt.Printf("\033[%d;%dH", row, col)
        fmt.Print(text)
        l.Unlock()
}