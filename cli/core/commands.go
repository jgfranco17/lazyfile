package core

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"gtithub.com/jgfranco17/lazyfile/cli/files"
)

type ShellExecutor interface {
	Exec(ctx context.Context, name string, args string) (int, string, error)
}

func CommandListFiles() *cobra.Command {
	var printAsTree bool
	cmd := &cobra.Command{
		Use:   "list",
		Args:  cobra.MaximumNArgs(1),
		Short: "List the files in a directory",
		Long:  "View the file system in an interactive TUI render",
		RunE: func(cmd *cobra.Command, args []string) error {
			var path string
			if len(args) > 0 {
				path = args[0]
			}
			entries, err := files.GetDirectoryContents(path)
			if err != nil {
				return err
			}
			baseDir, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			fmt.Println("Path:", baseDir)
			files.ListDirectoryContents(entries, printAsTree)
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().BoolVarP(&printAsTree, "tree", "t", false, "Render as tree")
	return cmd
}
