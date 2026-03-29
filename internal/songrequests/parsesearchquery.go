package songrequests

import (
	"net/url"
	"strings"
)

func ParseSearchQuery(s string) (string, bool) {
	s = strings.TrimPrefix(s, "!sr ")
	url, err := url.Parse(s)
	if err != nil {
		return s, false
	}
	if url.Host == "www.youtube.com" || url.Host == "music.youtube.com" {
		urlpath := strings.TrimPrefix(url.Path, "/")
		if urlpath == "watch" {
			vid := url.Query().Get("v")
			if vid != "" {
				return vid, true
			}
		}
	}
	if url.Host == "youtu.be" {
		urlpath := strings.TrimPrefix(url.Path, "/")
		urlpath = strings.Split(urlpath, "/")[0]
		i := strings.Index(urlpath, "#")
		if i != -1 {
			urlpath = urlpath[:i]
		}
		i = strings.Index(urlpath, "?")
		if i != -1 {
			urlpath = urlpath[:i]
		}
		if urlpath != "" {
			return urlpath, true
		}
	}
	return s, false
}
