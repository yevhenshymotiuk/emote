package emoticons

import (
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/carolynvs/emote/config"
	"github.com/spf13/viper"
)

type App struct {
	Out    io.Writer
	Config *config.Config
	Viper  *viper.Viper
}

func New() (*App, error) {
	v := viper.New()
	v.SetConfigFile(path.Join(os.Getenv("HOME"), "emote.toml"))

	c, err := config.Load(v)
	if err != nil {
		return nil, err
	}
	a := &App{
		Out:    os.Stdout,
		Config: c,
		Viper:  v,
	}
	return a, nil
}

func (a *App) Emote(name string, dest string) {
	emoticon, prs := a.Config.Emoticon[name]
	if !prs {
		fmt.Fprintf(a.Out, "There is no %s emote in the list\n", name)
		return
	}

	switch dest {
	case "clipboard":
		clipboard.WriteAll(emoticon)
		fmt.Fprintf(a.Out, "'%s' was copied to the clipboard\n", emoticon)
	default:
		fmt.Fprintln(a.Out, emoticon)
	}
}

func (a *App) PrintEmotesList() {
	emoticons := a.Config.Emoticon

	var maxNameLength int

	var keys []string
	for k := range emoticons {
		keyLen := len(k)
		if keyLen > maxNameLength {
			maxNameLength = keyLen
		}

		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		space := strings.Repeat(" ", maxNameLength - len(k))
		fmt.Fprintln(a.Out, k, space, emoticons[k])
	}
}
