package main

import "fmt"
import "strings"
import "strconv"
import "io/ioutil"
import "time"

// while true; do egrep 'xvda2|sda2' /proc/diskstats|head -1|awk '{print $14}'; sleep 1; done

const (
	// https://www.kernel.org/doc/Documentation/iostats.txt
	DISKSTATS = "/proc/diskstats"
)

func mustParseUint64(value string) uint64 {
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Could not parse value %v as uint64", value))
	}
	return result
}

type diskStats struct {
	Inflight       uint64
	WeightedIoTime uint64
}

func parseDiskStats(rawStats string) map[string]*diskStats {
	allStats := make(map[string]*diskStats)

	for _, line := range strings.Split(rawStats, "\n") {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		name := fields[2]

		if strings.HasPrefix(name, "ram") || strings.HasPrefix(name, "loop") {
			continue
		}

		allStats[name] = &diskStats{
			Inflight:       mustParseUint64(fields[11]),
			WeightedIoTime: mustParseUint64(fields[13]),
		}
	}

	return allStats
}

func count(name string, value uint64) string {
	return fmt.Sprintf("count#%v=%v", name, value)
}

func emitStats(disk string, stats *diskStats) {
	s := []string{
		count(fmt.Sprintf("%v.inflight", disk), stats.Inflight),
		count(fmt.Sprintf("%v.weighted-io-time", disk), stats.WeightedIoTime),
	}
	fmt.Println(strings.Join(s, " "))
}

func main() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			data, err := ioutil.ReadFile(DISKSTATS)
			if err != nil {
				panic(fmt.Sprintf("Could not read %v", DISKSTATS))
			}
			stats := parseDiskStats(string(data))
			if stats["sda2"] != nil {
				emitStats("sda2", stats["sda2"])
			}
			if stats["xvda2"] != nil {
				emitStats("xvda2", stats["xvda2"])
			}
		}
	}
}
