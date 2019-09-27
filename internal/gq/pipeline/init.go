package pipeline

import (
	"fmt"
	"github.com/pubgo/g/errors"
	"github.com/storyicon/graphquery/kernel/pipeline"
	"github.com/storyicon/graphquery/kernel/selector"
)

func Init() {
	errors.Panic(pipeline.RegistProcessor("test",
		func(sec selector.Selection, strings []string) (selection selector.Selection, e error) {
			fmt.Println(sec.Text())
			return sec, nil
		}, 1))
}
