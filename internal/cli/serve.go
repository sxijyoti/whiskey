package cli

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/build"
	"github.com/sxijyoti/whiskey/internal/devserver"
	"github.com/sxijyoti/whiskey/internal/watcher"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve [site-root]",
	Short: "Build and serve a Whiskey site",
	Args:  cobra.MaximumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {

		root := "site"

		if len(args) == 1 {
			root = args[0]
		}

		if err := build.BuildSite(root); err != nil {
			return err
		}

		reloader := devserver.NewReloader()

		w, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		defer w.Close()

		if err := watcher.WatchRecursive(
			w,
			filepath.Join(root, "content"),
		); err != nil {
			return err
		}

		if err := watcher.WatchRecursive(
			w,
			filepath.Join(root, "layouts"),
		); err != nil {
			return err
		}

		if err := watcher.WatchRecursive(
			w,
			filepath.Join(root, "static"),
		); err != nil {
			return err
		}

		debouncer := watcher.NewDebouncer(
			750 * time.Millisecond,
		)

		go func() {
			for {
				select {

				case event := <-w.Events:

					if event.Op&(fsnotify.Write|
						fsnotify.Create|
						fsnotify.Remove|
						fsnotify.Rename) == 0 {
						continue
					}

					fmt.Printf(
						"[watch] changed: %s\n",
						event.Name,
					)

					debouncer.Run(func() {

						fmt.Println(
							"[watch] rebuilding...",
						)

						if err := build.BuildSite(root); err != nil {
							fmt.Println(err)
							return
						}

						fmt.Println(
							"[watch] build complete",
						)

						reloader.Broadcast()
					})

				case err := <-w.Errors:
					fmt.Println(err)
				}
			}
		}()

		http.HandleFunc(
			"/__reload",
			reloader.Handler,
		)

		http.Handle(
			"/",
			devserver.NewFileHandler("dist"),
		)

		fmt.Printf(
			"Serving %s at http://localhost:%d\n",
			root,
			port,
		)

		return http.ListenAndServe(
			fmt.Sprintf(":%d", port),
			nil,
		)
	},
}

func init() {
	serveCmd.Flags().IntVarP(
		&port,
		"port",
		"p",
		8080,
		"server port",
	)

	rootCmd.AddCommand(
		serveCmd,
	)
}