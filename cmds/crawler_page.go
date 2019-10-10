package cmds

import (
	"bytes"
	"github.com/pubgo/g/errors"
	"github.com/pubgo/g/gcmds"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
	"time"
)
import "github.com/tidwall/gjson"

func init() {
	//var retryAt = gotry.RetryAt
	var url = ""
	var page = false
	var output = "term"

	type _Page struct {
		Page     int
		PageSize int
	}

	rootCmd := gcmds.Default()
	// https://turing.yuanben.site/extract_rules/api/rule/list?page={{.Page}}&limit={{.PageSize}}
	rootCmd.AddCommand(func(cmd *cobra.Command) *cobra.Command {
		cmd.Flags().StringVar(&url, "url", url, "抓取的url")
		cmd.Flags().StringVar(&output, "output", output, "输出类型,默认终端")
		cmd.Flags().BoolVar(&page, "page", page, "是否分页?")

		return cmd
	}(&cobra.Command{
		Use:   "crawler",
		Short: "mycli crawler",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer errors.RespErr(&err)

			tmpl, err := template.New("main").Parse(url)
			errors.PanicM(err, "url template 解析失败")

			_output := os.Stdout
			if output != "term" {
				f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0755)
				errors.PanicM(err, "打开文件错误")
				_output = f
			}
			defer _output.Close()

			pageSize := 50
			_output.WriteString("[")
			for page := 0; ; page++ {
				var _url = bytes.NewBufferString("")
				errors.PanicM(tmpl.Execute(_url, &_Page{
					Page:     page,
					PageSize: pageSize,
				}), "模板适配")

				resp, err := http.Get(_url.String())
				errors.PanicM(err, "获取数据失败")
				_body, err := ioutil.ReadAll(resp.Body)
				errors.Panic(err)
				if len(gjson.GetBytes(_body, "data.items").Array()) == 0 {
					break
				}
				_output.Write(_body)
				_output.WriteString(",")
			}
			_output.WriteString("]")

			time.Sleep(time.Second)
			return nil
		}},
	))
}
