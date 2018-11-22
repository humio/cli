package command

import "log"

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Map(vs []string, f func(string) testCase) []testCase {
	vsm := make([]testCase, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
