package cmds

import (
	"fmt"
	"github.com/pubgo/g/gcmds"
	"github.com/pubgo/mycli/pkg/wordcloud"
	"github.com/spf13/cobra"
	"image/color"
	"time"
)

func init() {
	rootCmd := gcmds.Default()
	rootCmd.AddCommand(func(cmd *cobra.Command) *cobra.Command {
		return cmd
	}(&cobra.Command{
		Use: "wordcloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			startedAt := time.Now().Unix()

			textList := []string{"恭喜", "发财", "万事", "如意"}
			angles := []int{0, 15, -15, 90}
			colors := []*color.RGBA{
				{0x0, 0x60, 0x30, 0xff},
				{0x60, 0x0, 0x0, 0xff},
				// &color.RGBA{0x73, 0x73, 0x0, 0xff},
			}
			render := wordcloud.NewWordCloudRender(60, 8,
				"./fonts/xin_shi_gu_yin.ttf",
				"./imgs/tiger.png",
				textList,
				angles,
				colors,
				"./imgs/tiger_template.png")
			render.Render()

			endAt := time.Now().Unix()
			fmt.Printf("时间消耗:%d\n", endAt-startedAt)
			return nil
		}},
	))
}
