package urlretriever

import (
	"bufio"
	"fmt"
	"github.com/temp25/hdl/helper"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func GetVideoFormats(masterPlaybackPageContents string, masterPlaybackUrl string) (map[string]interface{}, []int) {

	videoFormats := make(map[string]interface{})
	info := make(map[string]interface{})
	videoFormatKeys := make([]int, 0, len(videoFormats))
	scanner := bufio.NewScanner(strings.NewReader(masterPlaybackPageContents))
	m3u8InfoRegex := regexp.MustCompile(`([\w\-]+)\=([\w\-]+|"[^"]*")`)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#EXT-X-STREAM-INF") {
			line := strings.Replace(line, "#EXT-X-STREAM-INF:", "", -1)
			m3u8Info := m3u8InfoRegex.FindAllStringSubmatch(line, -1)

			for index := range m3u8Info {
				info[m3u8Info[index][1]] = m3u8Info[index][2]
			}

		} else if strings.HasPrefix(line, "master") || strings.HasPrefix(line, "http") {

			averageBandwidthOrBandwidth := func() int {
				var bw string
				if info["AVERAGE-BANDWIDTH"] != nil {
					bw = info["AVERAGE-BANDWIDTH"].(string)
				} else {
					bw = info["BANDWIDTH"].(string)
				}
				var bwInt, _ = strconv.Atoi(bw)
				return bwInt
			}()

			kFactor := averageBandwidthOrBandwidth / 1000

			kForm := fmt.Sprintf("%dk", kFactor)

			if strings.HasPrefix(line, "master") {
				line = strings.Replace(masterPlaybackUrl, "master.m3u8", line, -1)
			}

			info["STREAM-URL"] = line
			info["K-FORM"] = kForm
			key := fmt.Sprintf("hls-%d", kFactor)
			videoFormats[key] = helper.CopyMap(info)
			videoFormatKeys = append(videoFormatKeys, kFactor)

			for k := range info {
				delete(info, k)
			}

		} else {
			//do nothing
		}

	}
	if err := scanner.Err(); err != nil {
		// handle error
		panic(err)
	}

	sort.Ints(videoFormatKeys)

	return videoFormats, videoFormatKeys

}
