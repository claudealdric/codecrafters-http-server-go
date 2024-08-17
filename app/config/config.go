package config

import "flag"

var Directory string

func ParseFlags() {
	flag.StringVar(
		&Directory,
		"directory",
		"",
		"Specify the directory where files are stored",
	)
	flag.Parse()
}
