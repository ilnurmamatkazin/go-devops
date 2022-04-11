package model

import "fmt"

type Build struct {
	version string
	date    string
	commit  string
}

func NewBuild(version, date, commit string) *Build {
	return &Build{
		version: version,
		date:    date,
		commit:  commit,
	}
}

func (b *Build) Print() {
	fmt.Printf("Build version: %s\n", b.version)
	fmt.Printf("Build date: %s\n", b.date)
	fmt.Printf("Build commit: %s\n", b.commit)
}
