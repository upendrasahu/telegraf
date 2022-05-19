package telegrafFTD

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type TelegrafFTD struct {
}

func (_ *TelegrafFTD) Description() string {
	return "telegrafFTD input description"
}

func (_ *TelegrafFTD) SampleConfig() string {
	return "no config required"
}

func (s *TelegrafFTD) Gather(acc telegraf.Accumulator) error {
	log.Printf("D! [inputs.telegrafFTD] Gather() is called")
	command := exec.Command("sh", "-c", "ps -awx -o comm,%cpu,rss | grep telegraf")
	output, err := command.CombinedOutput()
	if err != nil {
		log.Printf("D! [inputs.telegrafFTD] error in executing command - ps command for telegraf: %v ", err)
	} else { //processing output
		for _, line := range strings.Split(string(output), "\n") {
			var words []string
			for _, word := range strings.Split(line, " ") {
				if word != "" {
					words = append(words, word)
				}
			}
			l := len(words)
			if l == 3 {
				telegraf_cpu, err := strconv.ParseFloat(words[1], 64)
				if err != nil {
					log.Printf("D! [inputs.telegrafFTD] error in parsing - cpu: %v - %v", telegraf_cpu, err)
					continue
				}
				telegraf_rss, err := strconv.ParseInt(words[2], 10, 64)
				if err != nil {
					log.Printf("D! [inputs.telegrafFTD] error in parsing - rss: %v - %v", telegraf_rss, err)
					continue
				}
				telegraf_rss = telegraf_rss * 1024 //Converting to bytes
				acc.AddGauge("telegraf", map[string]interface{}{"gauge": telegraf_cpu}, map[string]string{"telegraf": "cpu"}, time.Now())
				acc.AddGauge("telegraf", map[string]interface{}{"gauge": telegraf_rss}, map[string]string{"telegraf": "rss"}, time.Now())
			}
		}
	}
	return nil
}

func init() {

	log.Printf("D! [inputs.telegrafFTD] init() is called")
	inputs.Add("telegrafFTD", func() telegraf.Input {
		return &TelegrafFTD{}
	})
}