package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

// NewConfigCommand returns the config command which currently supports 'list'
func NewConfigCommand() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Manage todo-cli configuration",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List configured todo storage files",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "raw", Usage: "Print raw JSON config instead of formatted list"},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					// locate config file
					home, err := os.UserHomeDir()
					if err != nil {
						return cli.Exit(fmt.Sprintf("unable to find home directory: %v", err), 1)
					}
					cfgPath := filepath.Join(home, ".todo", "config.json")

					data, err := os.ReadFile(cfgPath)
					if err != nil {
						if os.IsNotExist(err) {
							fmt.Println("No config found. No todo lists have been initialized yet.")
							return nil
						}
						return cli.Exit(fmt.Sprintf("error reading config: %v", err), 2)
					}

					if c.Bool("raw") {
						fmt.Println(string(data))
						return nil
					}

					// parse into struct
					var cfg struct {
						Paths []string `json:"paths"`
					}
					if err := json.Unmarshal(data, &cfg); err != nil {
						// fallback to raw
						fmt.Println(string(data))
						return nil
					}

					if len(cfg.Paths) == 0 {
						fmt.Println("No configured todo storage files found.")
						return nil
					}

					fmt.Println("Configured todo storage files:")
					for i, p := range cfg.Paths {
						fmt.Printf("%d. %s\n", i+1, p)
					}

					return nil
				},
			},
		},
	}
}
