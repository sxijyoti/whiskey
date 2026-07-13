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
	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/devserver"
	"github.com/sxijyoti/whiskey/internal/parser"
	"github.com/sxijyoti/whiskey/internal/source"
	"github.com/sxijyoti/whiskey/internal/watcher"
)

var port int

func watchIfExists(
	w *fsnotify.Watcher,
	root string,
) error {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil
	}

	return watcher.WatchRecursive(
		w,
		root,
	)
}

var serveCmd = &cobra.Command{
	Use:   "serve [site-root]",
	Short: "Build and serve a Whiskey site",
	Args:  cobra.MaximumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {

		root := siteRoot(args)

		source.Offline = offline

		fmt.Printf("Whiskey %s\n\n", Version)
		fmt.Printf("Serving http://localhost:%d\n\n", port)
		fmt.Println("Initial build...")
		fmt.Println()

		if err := build.IncrementalBuild(root); err != nil {
			return err
		}

		cfg, err := config.Load(root)
		if err != nil {
			return err
		}

		reloader := devserver.NewReloader()

		source.Poller{
			Interval: 30 * time.Second,
		}.Start(func() {

			build.LogNoopBuilds = false
			defer func() {
				build.LogNoopBuilds = true
			}()
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

		if err := watchIfExists(
			w,
			filepath.Join(root, "content"),
		); err != nil {
			return err
		}

		if err := watchIfExists(
			w,
			filepath.Join(root, "layouts"),
		); err != nil {
			return err
		}

		if err := watchIfExists(
			w,
			filepath.Join(
				"themes",
				cfg.Theme,
				"layouts",
			),
		); err != nil {
			return err
		}

		if err := watchIfExists(
			w,
			filepath.Join(
				"themes",
				cfg.Theme,
				"static",
			),
		); err != nil {
			return err
		}

		if err := watchIfExists(
			w,
			filepath.Join(root, "static"),
		); err != nil {
			return err
		}

		configFile := filepath.Join(
			root,
			"whiskey.toml",
		)

		if _, err := os.Stat(configFile); err == nil {

			if err := w.Add(configFile); err != nil {
				return err
			}
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

					if filepath.Ext(event.Name) == ".md" {

						raw, err := os.ReadFile(event.Name)
						if err == nil {

							doc, err := parser.ParseFrontmatter(raw)
							if err == nil && doc.Meta.Draft {
								continue
							}
						}
					}

					name := event.Name
					if rel, err := filepath.Rel(root, event.Name); err == nil {
						name = rel
					}

					fmt.Printf("[watch] %s changed\n", filepath.ToSlash(name))

					changed := event.Name

					debouncer.Run(func(path string) func() {
						return func() {

							if err := build.IncrementalBuild(root); err != nil {

								fmt.Printf(
									"[error] build failed: %v\n",
									err,
								)

								return
							}

							reloader.Broadcast()
						}
					}(changed))

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
				root,
				"dist",
			),
		)

		fmt.Println()
		fmt.Println("Watching for changes...")

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

	serveCmd.Flags().BoolVar(
		&offline,
		"offline",
		false,
		"build using cached remote sources only",
	)

	rootCmd.AddCommand(
		serveCmd,
	)
}
