// 'date.go'.
// Chris Shiels.


package main


import (
    "encoding/json"
    "flag"
    "fmt"
    "net/http"
    "os"
    "strings"
    "time"
)


var version string
const defaultport = 7001


const exitsuccess = 0
const exitfailure = 1


func splithostport(s string) (host string, port string) {
    var i int
    if i = strings.LastIndex(s, ":"); i == -1 {
        return s, ""
    }

    return s[0:i], s[i + 1:]
}


type app struct {
    name string
    port int
    version string
}


func newapp(name string, port int, version string) *app {
    return &app{ name: name, port: port, version: version }
}


func (a *app) date() (jsonbody string, err error) {
    year, month, day := time.Now().Date()
    date := fmt.Sprintf("%04d%02d%02d", year, month, day)

    hostname, err := os.Hostname()
    if err != nil {
        return "", err
    }

    type message struct {
        Date string `json:"date"`
        Hostname string `json:"hostname"`
        Version string `json:"version"`
    }
    m := message { date, hostname, a.version }
    bytes, err := json.Marshal(m)
    if err != nil {
        return "", err
    }

    return string(bytes), nil
}


func (a *app) status() (jsonbody string, err error) {
    type message struct {
        Ok bool `json:"ok"`
    }
    m := message { true }
    bytes, err := json.Marshal(m)
    if err != nil {
        return "", err
    }

    return string(bytes), nil
}


func (a *app) do(f func() (string, error)) func(responsewriter http.ResponseWriter, request *http.Request) {

    return func(responsewriter http.ResponseWriter,
                request *http.Request) {
        var status int

        jsonbody, err := f()
        if err != nil {
            status = http.StatusInternalServerError
            http.Error(responsewriter, err.Error(), status)
        } else {
            status = http.StatusOK
            responsewriter.Header().Set("Content-Type", "text/json")
            fmt.Fprint(responsewriter, jsonbody)
        }

        remotehost, _ := splithostport(request.RemoteAddr)
        time := time.Now().Format("02/Jan/2006:15:04:05 -0700")

        var referrer string
        if len(request.Header["Referer"]) == 0 {
            referrer = "-"
        } else {
            referrer = request.Header["Referer"][0]
        }

        fmt.Printf("%s %s %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n",
                   remotehost,                       // Remote host.
                   "-",                              // Remote logname
                                                     // (from identd).
                   "-",                              // Remote user.
                   time,                             // Time.
                   request.Method,                   // First line of request.
                   request.URL,
                   request.Proto,
                   status,                           // Status.
                   len(jsonbody),                    // Bytes sent.
                   referrer,                         // Referrer.
                   request.Header["User-Agent"][0])  // User-agent.
    }
}


func (a *app) start() error {
    fmt.Printf("Starting %s version %s on port %d.\n",
               a.name, a.version, a.port)

    http.HandleFunc("/date", a.do(a.date))
    http.HandleFunc("/status", a.do(a.status))

    return http.ListenAndServe(fmt.Sprintf(":%d", a.port), nil);
}


func main() {
    flag.Usage = func() {
        fmt.Println("Usage:  date -p port")
        flag.PrintDefaults()
    }
    flagp := flag.Int("p", defaultport, "Port number")

    // Note flag.Parse() will also handle '-h' and '--help' and will exit with
    // exit status 2.
    flag.Parse()


    if err := newapp("date", *flagp, version).start(); err != nil {
        fmt.Fprintf(os.Stderr, "date: %s\n", err)
        os.Exit(exitfailure)
    }

    os.Exit(exitsuccess)
}
