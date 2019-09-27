package cmds

import (
	"fmt"
	"github.com/pubgo/g/errors"
	"github.com/pubgo/g/gotry"
	"github.com/spf13/cobra"
	"time"
)
import "github.com/google/uuid"

func init() {
	var retryAt = gotry.RetryAt
	var n = 1

	rootCmd.AddCommand(func(cmd *cobra.Command) *cobra.Command {
		cmd.Flags().IntVar(&n, "n", n, "生成UUID的数量")

		return cmd
	}(&cobra.Command{
		Use:   "uuid",
		Short: "mycli uuid",
		RunE: func(cmd *cobra.Command, args []string) error {
			for i := n; i > 0; i-- {
				retryAt(time.Second, func(at time.Duration) {
					_uuid, err := uuid.NewRandom()
					errors.PanicM(err, "uuid gen error")
					fmt.Println(_uuid.Version(), _uuid.ID(), _uuid.ClockSequence(), _uuid.String())
				})
			}

			for i := n; i > 0; i-- {
				retryAt(time.Second, func(at time.Duration) {
					_uuid, err := uuid.NewUUID()
					errors.PanicM(err, "uuid gen error")
					fmt.Println(_uuid.Version(), _uuid.ID(), _uuid.ClockSequence(), _uuid.String())
				})
			}

			time.Sleep(time.Second)
			return nil
		}},
	))
}
