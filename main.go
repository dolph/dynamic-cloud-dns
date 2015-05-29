package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
)


func icanhazip() string {
    // this requires Go 1.4 due to the TLS implementation
    resp, err := http.Get("https://icanhazip.com/")
    if err != nil {
        fmt.Println("Unable to fetch icanhazip.")
        fmt.Println(err)
        os.Exit(1)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Unable to read response body.")
        fmt.Println(err)
        os.Exit(1)
    }

    // icanhazip does nothing but return us an IP
    ip := string(body)

    // remove the trailing newline
    ip = strings.TrimRight(ip, "\n")

    return ip
}


func main() {
    ip := icanhazip()

    fmt.Println(ip)
}
