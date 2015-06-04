package main

import (
    "bytes"
    "flag"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "net/http"
    "os"
    "strings"
)


func fatalError(err error, message string) {
    if err != nil {
        fmt.Println(message)
        fmt.Println(err)
        panic(err)
    }
}

func icanhazip() string {
    // this requires Go 1.4 due to the TLS implementation
    resp, err := http.Get("https://icanhazip.com/")
    fatalError(err, "Unable to fetch icanhazip.")

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    fatalError(err, "Unable to read response body from icanhazip.")

    // icanhazip does nothing but return us an IP
    ip := string(body)

    // remove the trailing newline
    ip = strings.TrimRight(ip, "\n")

    return ip
}

func authenticate(username string, apiKey string) string {
    client := &http.Client{}

    body := fmt.Sprintf(
        `{
            "auth": {
                "RAX-KSKEY:apiKeyCredentials": {
                    "username": "%s",
                    "apiKey": "%s"
                }
            }
        }`,
        username, apiKey)
    req, err := http.NewRequest(
        "POST",
        "https://identity.api.rackspacecloud.com/v2.0/tokens",
        bytes.NewBuffer([]byte(body)))
    req.Header.Set("Content-Type", "application/json")
    fatalError(err, "Unable to build authentication request.")
    resp, err := client.Do(req)
    fmt.Println("HTTP", resp.StatusCode)
    if resp.StatusCode == 401 {
        fatalError(err, "Unable to authenticate with Rackspace.")
    }

    text, err := ioutil.ReadAll(resp.Body)
    fatalError(err, "Unable to read authentication response.")

    var data map[string]interface{}
    if err := json.Unmarshal(text, &data); err != nil {
        panic(err)
    }
    access := data["access"]
    token := access.(map[string]interface{})["token"]
    return token.(map[string]interface{})["id"].(string)
}

func main() {
    ip := icanhazip()

    fmt.Println(ip)

    usernamePtr := flag.String("username", "", "Rackspace account username")
    apiKeyPtr := flag.String("api-key", "", "Rackspace API key")

    flag.Parse()

    if *usernamePtr == "" || *apiKeyPtr == "" {
        fmt.Println("Rackspace -username and -api-key required.")
        fmt.Println("See -help for more info.")
        os.Exit(1)
    }

    token := authenticate(*usernamePtr, *apiKeyPtr)
}
