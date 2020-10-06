// Copyright 2015 go-fuzz project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package regexp

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/dop251/goja"
)

func Fuzz(data []byte) int {
	if bytes.Contains(bytes.ToLower(data), []byte("\\p")) {
		// this is not supported
		return 0
	}
	sstr := string(data[:len(data)/2])
	// restrb := data[len(data)/2:]
	restr := string(data[len(data)/2:])

	score := 0
	vm := goja.New()

	execString := func(template string) func(s, r string) {
		return func(s, r string) {
			str := fmt.Sprintf(template, s, r)
			defer func() {
				if r := recover(); r != nil {
					if !strings.Contains(fmt.Sprint(r), "Invalid regular expression") {
						panic(fmt.Sprintf("%s\n'%#v'\n'%s'\n", r, []byte(str), str))
					}
				}
			}()
			vm.RunString(str)
		}
	}

	re2Exec := func(s, r string) {
		re, err := regexp2.Compile(r, regexp2.ECMAScript|regexp2.Multiline)
		if err != nil {
			return
		}
		score = 1
		re.FindStringMatch(s)

		re, err = regexp2.Compile(r, regexp2.ECMAScript)
		if err != nil {
			return
		}
		re.FindStringMatch(s)
	}
	for _, fn := range [...]func(s, r string){
		re2Exec,
		execString(`/%s/gu.exec(%q)`),
		execString(`/%s/g.exec(%q)`),
		execString(`/%s/.exec(%q)`),
		execString(`%q.split(/%s/gu)`),
		execString(`%q.split(/%s/g)`),
		execString(`%q.split(/%s/)`),
	} {
		func() {
			fn(sstr, restr)
			fn(restr, sstr)
		}()
	}

	return score
}
