//go:build windows

package main

import title "github.com/lxi1400/GoTitle"

func setTitle(s string) {
	title.SetTitle(s)
}
