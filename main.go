package main

import "crypto/tls"
import "fmt"
import "io/ioutil"
import "net/http"
import "os"
import "strings"

import "github.com/certifi/gocertifi"


func main() {
    cert_pool, err := gocertifi.CACerts()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    transport := &http.Transport{
        TLSClientConfig: &tls.Config{RootCAs: cert_pool},
    }
    client := &http.Client{
        Transport: transport,
    }

    resp, err := client.Get("https://icanhazip.com/")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // icanhazip does nothing but return us an IP
    ip := string(body)

    // remove the trailing newline
    ip = strings.TrimRight(ip, "\n")

    fmt.Println(ip)
}
