package converter

import (
	docker "docker.io/go-docker"
	"github.com/pkg/errors"
)

type pandoc struct {
	cli *docker.Client
}

func NewPandocConverter() (Converter, error) {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "Error trying to access to local docker daemon")
	}

	return &pandoc{cli: cli}, nil
}

func (p *pandoc) Convert(opts ...[]string) {

}
