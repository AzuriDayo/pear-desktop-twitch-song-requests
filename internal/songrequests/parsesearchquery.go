package songrequests

import (
	"net/url"
	"strings"
)

func ParseSearchQuery(s string) string {
	s = strings.TrimPrefix(s, "!sr ")
	url, err := url.Parse(s)
	if err != nil {
		return s
	}
	if url.Host == "www.youtube.com" || url.Host == "music.youtube.com" {
		urlpath := strings.TrimPrefix(url.Path, "/")
		if urlpath == "watch" || urlpath == "/watch" {
			vid := url.Query().Get("v")
			if vid != "" {
				return vid
			}
		}
	}
	if url.Host == "youtu.be" {
		urlpath := strings.TrimPrefix(url.Path, "/")
		if urlpath != "" {
			return urlpath
		}
	}
	return s
}
