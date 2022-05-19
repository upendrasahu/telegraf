package active_monitors

import (
    "github.com/influxdata/telegraf"
    "github.com/influxdata/telegraf/plugins/inputs"

    "fmt"
    "log"
    "net"
    "net/http"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
    "time"
)

type Active_monitors struct {
	Host string
}

var ping_time float64

func (_ *Active_monitors) Description() string {
    return "active_monitors input description"
}

func (_ *Active_monitors) SampleConfig() string {
    return "no config required"
}

func curl_it(host string, port string) {
    timeout := time.Duration(1 * time.Second)
    _, err := net.DialTimeout("tcp", host+":"+port, timeout)
    if err != nil {
        log.Printf("D! [inputs.active_monitors] %s %s %s\n", host, "not responding", err.Error())
    } else {
        log.Printf("D! [inputs.active_monitors] %s %s %s\n", host, "responding on port:", port)
    }
}

func call_ping_it() func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        ping_it("8.8.8.8")
        s := fmt.Sprintf("Hello... ping time = %0.2f\n", ping_time)
        w.Write([]byte(s))
    }
}

func ping_it(host string) {
    ping_time = -1.0
    command := exec.Command("sh", "-c", "timeout 10 ping -c 1 "+host+" 2>&1")
    output, err := command.CombinedOutput()
    if err != nil {
        log.Printf("D! [inputs.active_monitors] error: %v\n", err)
    } else {
        /*
        PING 8.8.8.8 (8.8.8.8) 56(84) bytes of data.
        64 bytes from 8.8.8.8: icmp_seq=1 ttl=48 time=9.79 ms

        --- 8.8.8.8 ping statistics ---
        1 packets transmitted, 1 received, 0% packet loss, time 0ms
        rtt min/avg/max/mdev = 9.794/9.794/9.794/0.000 ms
        */
        for line_number, line := range strings.Split(string(output), "\n") {
            if (line_number == 1){
                tokens := regexp.MustCompile(`[= ]`).Split(line, -1)
                log.Printf("D! [inputs.active_monitors] tokens: %v", tokens)
                ping_time, err = strconv.ParseFloat(tokens[9], 64)
                if (err != nil){
                    log.Printf("D! [inputs.active_monitors] error: %v - %v", ping_time, err)
                    break
                }
                log.Printf("D! [inputs.active_monitors] value: %v", ping_time)

            }
        }
    }
}


func (s *Active_monitors) Gather(acc telegraf.Accumulator) error {
    log.Printf("D! [inputs.active_monitors] Gather() is called")
	fields := map[string]interface{}{}
    ping_it(s.Host)
	fields["ping_time_ms"] = ping_time
	tags := map[string]string{"server": s.Host}
    // acc.AddGauge("active_monitors_ping_time_ms", map[string]interface{}{"gauge": ping_time}, tags, time.Now())
    // call addmetric
	acc.AddFields("active_monitors", fields, tags)
    if (ping_time > 10) {
        //do something
    }
    curl_it(s.Host, "80")
    return nil
}

func setup_server() {
    mux := http.NewServeMux()
    mux.HandleFunc("/ping", call_ping_it())

    server := &http.Server{Handler: mux}
    listener, err := net.Listen("tcp", "127.0.0.1:10101")

    if err != nil {
        log.Printf("D! [inputs.active_monitors] Listener Error: %v", err)
        _, dial_err := net.Dial("tcp", "127.0.0.1:10101")
        if dial_err == nil {
            log.Printf("D! [inputs.active_monitors] Listener connected")
        } else {
            log.Printf("D! [inputs.active_monitors] Listener Dial Error: %v", dial_err)
        }
    }

    go func() {
        err := server.Serve(listener)
        if err != nil && err != http.ErrServerClosed {
            log.Printf("E! Error creating prometheus metric endpoint, err: %s\n", err.Error())
        }
    }()
}

func init() {
    log.Printf("D! [inputs.active_monitors] init() is called")
    setup_server()
    inputs.Add("active_monitors", func() telegraf.Input {
        return &Active_monitors{}
    })
}