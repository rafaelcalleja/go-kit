package helper

import "github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"

type jxhelper struct {
}

func (j *jxhelper) CheckErr(err error) {
	helper.CheckErr(err)
}

func (j *jxhelper) BehaviorOnFatal(f func(string, int)) {
	helper.BehaviorOnFatal(f)
}

func newJenkinsXHelper() *jxhelper {
	return new(jxhelper)
}
