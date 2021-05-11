package cmd

import "strings"

type stringSlice []string

func (ss *stringSlice) String() string {
	return strings.Join(*ss, ",")
}

func (ss *stringSlice) Set(val string) error {
	*ss = append(*ss, val)
	return nil
}
