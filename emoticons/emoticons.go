package emoticons

import (
	"fmt"
	"io"
	"os"
	"path"

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
	emoticon := a.Config.Emoticon[name]

	switch dest {
	case "clipboard":
		clipboard.WriteAll(emoticon)
		fmt.Fprintf(a.Out, "'%s' was copied to the clipboard\n", emoticon)
	default:
		fmt.Fprintln(a.Out, emoticon)
	}
}
