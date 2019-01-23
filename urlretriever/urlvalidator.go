package urlretriever

import (
  "regexp"
)

func reSubMatchMap(r *regexp.Regexp, str string) (map[string]string) {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}
	
	return subMatchMap
}

func IsValidHotstarUrl(videoUrl string) (bool, string) {
    var urlRegex = regexp.MustCompile(`((https|http)?://)?(www\.)?hotstar\.com/(?:.+?[/-])+(?P<videoId>\d{10})`)
    if urlRegex.MatchString(videoUrl) {
        match := reSubMatchMap(urlRegex, videoUrl)
        return true, match["videoId"]
    }

    return false, ""
}
