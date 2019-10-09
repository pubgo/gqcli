package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
	"strings"
)

type sHandler struct {
}

// http://localhost:8080/?s=?__and__b=?__or__c=?;[1,2,3]
func (f *sHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	fmt.Println(r.URL.RawQuery)
	w.Write([]byte(r.URL.RawQuery))
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "ss",
		Short: "mycli simple server",
		RunE: func(cmd *cobra.Command, args []string) error {

			//cwd, err := os.Getwd()
			//if err != nil {
			//	return err
			//}

			listener, err := net.Listen("tcp", ":8080")
			if err != nil {
				return err
			}

			log.Println("Listening on", listener.Addr())
			//log.Fatal(http.Serve(listener, http.FileServer(http.Dir(cwd))))
			log.Fatal(http.Serve(listener, &sHandler{}))
			return nil
		},
	})
}
