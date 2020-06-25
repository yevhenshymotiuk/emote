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
	cmd := buildEmoteCommand(app)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildEmoteCommand(app *emoticons.App) *cobra.Command {
	var dest string

	emote := &cobra.Command{
		Use: "emote",
		PreRun: func(cmd *cobra.Command, args []string) {
			app.Out = cmd.OutOrStdout()
		},
		Run: func(cmd *cobra.Command, args []string) {
			emoticonName := args[0]
			app.Emote(emoticonName, dest)
		},
		Args: cobra.ExactArgs(1),
	}

	emote.Flags().StringVar(&dest, "dest", "clipboard", "Where to send your emoticon")

	emote.AddCommand(buildConfigCommand(app))
	emote.AddCommand(buildListCommand(app))

	return emote
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

func buildListCommand(app *emoticons.App) *cobra.Command {
	list := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			app.PrintEmotesList()
		},
		Args: cobra.ExactArgs(0),
	}

	return list
}
