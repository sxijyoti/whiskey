package cli

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/build"
	"github.com/sxijyoti/whiskey/internal/devserver"
	"github.com/sxijyoti/whiskey/internal/watcher"
	"github.com/sxijyoti/whiskey/internal/source"
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

		if err := build.IncrementalBuild(root); err != nil {
			return err
		}

		reloader := devserver.NewReloader()

		source.Poller{
			Interval: 30 * time.Second,
		}.Start(func() {

			if err := build.IncrementalBuild(
				root,
			); err != nil {

				fmt.Println(err)
				return
			}

			reloader.Broadcast()
		})

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
			filepath.Join(
				"themes",
				"default",
				"layouts",
			),
		); err != nil {
			return err
		}

		if err := watcher.WatchRecursive(
			w,
			filepath.Join(
				"themes",
				"default",
				"static",
			),
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

				case event, ok := <-w.Events:

					if !ok {
						return
					}

					if watcher.IgnoreFile(
						event.Name,
					) {
						continue
					}

					if event.Op&(fsnotify.Write|
						fsnotify.Create|
						fsnotify.Remove|
						fsnotify.Rename) == 0 {
						continue
					}

					// watch newly-created directories
					if event.Op&fsnotify.Create != 0 {

						info, err := os.Stat(
							event.Name,
						)

						if err == nil &&
							info.IsDir() {

							_ = watcher.WatchRecursive(
								w,
								event.Name,
							)
						}
					}

					fmt.Printf(
						"[watch] queued: %s\n",
						filepath.Base(
							event.Name,
						),
					)

					debouncer.Run(
						func() {

							start := time.Now()

							fmt.Println(
								"[watch] rebuilding...",
							)

							if err := build.IncrementalBuild(root); err != nil {

								fmt.Printf(
									"[watch] build failed: %v\n",
									err,
								)

								return
							}

							fmt.Printf(
								"[watch] build complete (%v)\n",
								time.Since(start).
									Round(time.Millisecond),
							)

							reloader.Broadcast()
						},
					)

				case err, ok := <-w.Errors:

					if !ok {
						return
					}

					fmt.Printf(
						"[watch] error: %v\n",
						err,
					)
				}
			}
		}()

		http.HandleFunc(
			"/__reload",
			reloader.Handler,
		)

		http.Handle(
			"/",
			devserver.NewFileHandler(
				"dist",
			),
		)

		fmt.Printf(
			"Serving %s at http://localhost:%d\n",
			root,
			port,
		)

		return http.ListenAndServe(
			fmt.Sprintf(
				":%d",
				port,
			),
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
