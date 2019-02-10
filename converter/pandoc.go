package converter

import (
	"context"
	"io/ioutil"
	"log"

	docker "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"github.com/pkg/errors"
)

type pandoc struct {
	cli *docker.Client
}

const (
	dockerPandocImage = "jagregory/pandoc"
)

// NewPandocConverter a converter based on pandoc application
func NewPandocConverter() (Converter, error) {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "error trying to access to local docker daemon")
	}

	return &pandoc{cli: cli}, nil
}

func (p *pandoc) Convert(opts ...[]string) error {

	cmd := []string{"-s", "-r", "html", "https://go.googlesource.com/proposal/+/master/design/go2draft-contracts.md", "-o", "/source/test.epub"}

	c, err := p.cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:      dockerPandocImage,
			Cmd:        cmd,
			WorkingDir: "/source",
			Tty:        true,
		},
		nil,
		nil,
		"",
	)

	if err != nil {
		return errors.Wrapf(err, "error trying to create container %s for convert files to epub", c.ID)
	}

	if err := p.exec(c.ID); err != nil {
		return err
	}

	if err := p.saveOutputFile(c.ID); err != nil {
		return err
	}

	defer p.removeContainer(c.ID)

	return nil
}

func (p *pandoc) saveOutputFile(containerID string) error {
	archive, _, err := p.cli.CopyFromContainer(context.Background(), containerID, "/source/test.epub")
	if err != nil {
		return errors.Wrapf(err, "final file %s is not created", "test.epub")
	}

	content, err := ioutil.ReadAll(archive)
	if err != nil {
		return errors.Wrap(err, "the output file can not be created")
	}

	ioutil.WriteFile("test.tar.gz", content, 0644)
	return nil
}

func (p *pandoc) exec(containerID string) error {
	err := p.cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return errors.Wrap(err, "error trying to convert files into a epub")
	}

	c, errC := p.cli.ContainerWait(context.Background(), containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errC:
		return errors.Wrapf(err, "error trying to convert files into a epub")
	case result := <-c:
		if result.StatusCode != 0 {
			return errors.Wrapf(err, "expected a status code equal to '0', got %d", result.StatusCode)
		}
	}

	return nil
}

func (p *pandoc) removeContainer(containerID string) error {
	if err := p.cli.ContainerRemove(
		context.Background(),
		containerID, types.ContainerRemoveOptions{}); err != nil {
		err := errors.Wrapf(err, "error trying to remove the container %s", containerID)
		log.Println(err)
	}

	return nil
}
