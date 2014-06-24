// +build darwin

package darwin

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

type CpuusageGenerator struct {
}

var cpuusageLogger = logging.GetLogger("metrics.cpuusage")

var iostatFieldToMetricName = []string{"user", "system", "idle"}

// Generate returns current CPU usage of the host.
// Keys below are expected:
// - cpu.user.percentage
// - cpu.system.percentage
// - cpu.idle.percentage
func (g *CpuusageGenerator) Generate() (metrics.Values, error) {

	// $ iostat -n0 -c2
	//         cpu     load average
	//    us sy id   1m   5m   15m
  //    13  7 81  1.93 2.23 2.65
  //    13  7 81  1.93 2.23 2.65
	iostatBytes, err := exec.Command("iostat", "-n0", "-c2").Output()
	if err != nil {
		cpuusageLogger.Errorf("Failed to invoke iostat: %s", err)
		return nil, err
	}

	iostat := string(iostatBytes)
	lines := strings.Split(iostat, "\n")
	if len(lines) != 5 {
		return nil, fmt.Errorf("iostat result malformed: [%q]", iostat)
	}

	fields := strings.Fields(lines[3])
	if len(fields) < len(iostatFieldToMetricName) {
		return nil, fmt.Errorf("iostat result malformed: [%q]", iostat)
	}

	cpuusage := make(map[string]float64, len(iostatFieldToMetricName))

	for i, n := range iostatFieldToMetricName {
		value, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return nil, err
		}

		cpuusage["cpu."+n+".percentage"] = value
	}

	return metrics.Values(cpuusage), nil
}