package main

import (
	"fmt"
	"log"
	"os"

	"github.com/carolynvs/emote/config"
	"github.com/carolynvs/emote/emoticons"
	"github.com/spf13/cobra"
)

func main() {
	app, err := emoticons.New()
	if err != nil {
		log.Fatal(err)
	}
	cmd, err := buildEmoteCommand(app)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildEmoteCommand(app *emoticons.App) (emote *cobra.Command, err error) {
	var dest string

	emote = &cobra.Command{
		Use: "emote",
		PreRun: func(cmd *cobra.Command, args []string) {
			app.Out = cmd.OutOrStdout()
		},
		Run: func(cmd *cobra.Command, args []string) {
			emoticonName := args[0]
			err = app.Emote(emoticonName, dest)
		},
		Args: cobra.ExactArgs(1),
	}
	if err != nil {
		return
	}

	emote.Flags().StringVar(&dest, "dest", "clipboard", "Where to send your emoticon")

	emote.AddCommand(buildConfigCommand(app))

	listCommand, err := buildListCommand(app)
	if err != nil {
		return
	}
	emote.AddCommand(listCommand)

	return
}

func buildConfigCommand(app *emoticons.App) *cobra.Command {
	config := &cobra.Command{
		Use: "config",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Config file is", config.File(app.Viper))
		},
		Args: cobra.ExactArgs(0),
	}

	return config
}

func buildListCommand(app *emoticons.App) (list *cobra.Command, err error) {
	list = &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			err = app.PrintEmotesList()
		},
		Args: cobra.ExactArgs(0),
	}

	return
}
