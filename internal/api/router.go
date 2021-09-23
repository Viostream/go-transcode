package api

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	"github.com/m1k1o/go-transcode/internal/config"
)

var conf *config.YamlConf

func init() {
	var err error
	conf, err = config.LoadConf("streams.yaml")
	if err != nil {
		panic(err)
	}
}

type ApiManagerCtx struct {
	Conf *config.YamlConf
}

func New() *ApiManagerCtx {
	return &ApiManagerCtx{Conf: conf}
}

func (a *ApiManagerCtx) Mount(r *chi.Mux) {
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		//nolint
		w.Write([]byte("pong"))
	})

	r.Group(a.HLS)
	r.Group(a.Http)
}

func transcodeStart(folder string, profile string, input string) (*exec.Cmd, error) {
	url, ok := conf.Streams[input]
	if !ok {
		return nil, fmt.Errorf("stream not found")
	}

	re := regexp.MustCompile(`^[0-9A-Za-z_-]+$`)
	if !re.MatchString(profile) {
		return nil, fmt.Errorf("invalid profile path")
	}

	profilePath := fmt.Sprintf("%s/%s/%s.sh", conf.BaseDir, folder, profile)
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return nil, err
	}

	log.Info().Str("profilePath", profilePath).Str("url", url).Msg("command startred")
	return exec.Command(profilePath, url), nil
}
