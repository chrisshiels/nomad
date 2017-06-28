// 'web.go'.
// Chris Shiels.


package main


import (
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "net"
    "net/http"
    "os"
    "strings"
    "time"
)


var version string
const (
    defaultport = 7000
    defaultdateendpoint = "date.service.consul"
    defaulttimeendpoint = "time.service.consul"
)


const exitsuccess = 0
const exitfailure = 1


func splithostport(s string) (host string, port string) {
    var i int
    if i = strings.LastIndex(s, ":"); i == -1 {
        return s, ""
    }

    return s[0:i], s[i + 1:]
}


func lookupendpoint(endpoint string) (host string, port int, err error) {
    _, addrs, err := net.LookupSRV("", "", endpoint)
    if err != nil {
        return "", 0, err
    }

    return addrs[0].Target, int(addrs[0].Port), nil
}


func getendpoint(endpoint string, path string) (bytes []byte, err error) {
    host, port, err := lookupendpoint(endpoint)
    if err != nil {
        return nil, err
    }

    resp, err := http.Get(fmt.Sprintf("http://%s:%d/%s", host, port, path))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    bytes, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return bytes, err
}


func getdate(endpoint string) (date string,
                               hostname string,
                               version string,
                               err error) {
    bytes, err := getendpoint(endpoint, "/date")
    if err != nil {
        return "", "", "", err
    }

    type message struct {
        Date string `json:"date"`
        Hostname string `json:"hostname"`
        Version string  `json:"version"`
    }
    var m message
    err = json.Unmarshal(bytes, &m)
    if err != nil {
        return "", "", "", err
    }

    return m.Date, m.Hostname, m.Version, nil
}


func gettime(endpoint string) (time string,
                               hostname string,
                               version string,
                               err error) {
    bytes, err := getendpoint(endpoint, "/time")
    if err != nil {
        return "", "", "", err
    }

    type message struct {
        Time string `json:"time"`
        Hostname string `json:"hostname"`
        Version string `json:"version"`
    }
    var m message
    err = json.Unmarshal(bytes, &m)
    if err != nil {
        return "", "", "", err
    }

    return m.Time, m.Hostname, m.Version, nil
}


type app struct {
    name string
    port int
    dateendpoint string
    timeendpoint string
    version string
}


func newapp(name string,
            port int,
            dateendpoint string,
            timeendpoint string,
            version string) *app {
    return &app{ name: name,
                 port: port,
                 dateendpoint: dateendpoint,
                 timeendpoint: timeendpoint,
                 version: version }
}


func (a *app) home() (textbody string, err error) {
    var buffer bytes.Buffer

    hostname, err := os.Hostname()
    if err == nil {
        fmt.Fprintf(&buffer, "%s %s:\n", hostname, a.version)
    } else {
        fmt.Fprintf(&buffer, "unknown: %s\n", a.version)
    }

    date, hostname, version, err := getdate(a.dateendpoint)
    if err == nil {
        fmt.Fprintf(&buffer, "%s - %s %s\n", date, hostname, version)
    } else {
        fmt.Fprintf(&buffer, "unknown - %s\n", err)
    }

    time, hostname, version, err := gettime(a.timeendpoint)
    if err == nil {
        fmt.Fprintf(&buffer, "%s - %s %s\n", time, hostname, version)
    } else {
        fmt.Fprintf(&buffer, "unknown - %s\n", err)
    }

    return buffer.String(), nil
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

        textbody, err := f()
        if err != nil {
            status = http.StatusInternalServerError
            http.Error(responsewriter, err.Error(), status)
        } else {
            status = http.StatusOK
            responsewriter.Header().Set("Content-Type", "text/plain")
            fmt.Fprint(responsewriter, textbody)
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
                   len(textbody),                    // Bytes sent.
                   referrer,                         // Referrer.
                   request.Header["User-Agent"][0])  // User-agent.
    }
}


func (a *app) start() error {
    fmt.Printf("Starting %s version %s on port %d.\n",
               a.name, a.version, a.port)

    http.HandleFunc("/", a.do(a.home))
    http.HandleFunc("/status", a.do(a.status))

    return http.ListenAndServe(fmt.Sprintf(":%d", a.port), nil);
}


func main() {
    flag.Usage = func() {
        fmt.Println("Usage:  web -p port -d dateendpoint -t timeendpoint")
        flag.PrintDefaults()
    }
    flagp := flag.Int("p", defaultport, "Port number")
    flagdateendpoint :=
        flag.String("dateendpoint", defaultdateendpoint, "Date endpoint")
    flagtimeendpoint :=
        flag.String("timeendpoint", defaulttimeendpoint, "Time endpoint")

    // Note flag.Parse() will also handle '-h' and '--help' and will exit with
    // exit status 2.
    flag.Parse()


    if err := newapp("web",
                     *flagp,
                     *flagdateendpoint,
                     *flagtimeendpoint,
                     version).start(); err != nil {
        fmt.Fprintf(os.Stderr, "web: %s\n", err)
        os.Exit(exitfailure)
    }

    os.Exit(exitsuccess)
}
