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

// App contains main dependencies of the application
type App struct {
	Out    io.Writer
	Config *config.Config
	Viper  *viper.Viper
}

// New constructs new instance of the application
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
// Emote sends chosen emoticon to the destination
func (a *App) Emote(name string, dest string) error {
	emoticon, prs := a.Config.Emoticon[name]
	if !prs {
		_, err := fmt.Fprintf(a.Out, "There is no %s emote in the list\n", name)
		if err != nil {
			return err
		}
		return nil
	}

	switch dest {
	case "clipboard":
		err := clipboard.WriteAll(emoticon)
		_, err = fmt.Fprintf(a.Out, "'%s' was copied to the clipboard\n", emoticon)
		if err != nil {
			return err
		}
	default:
		_, err := fmt.Fprintln(a.Out, emoticon)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrintEmotesList prints list of emotes
func (a *App) PrintEmotesList() error {
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
		_, err := fmt.Fprintln(a.Out, k, space, emoticons[k])
		if err != nil {
			return err
		}
	}

	return nil
}
