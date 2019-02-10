package converter

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	docker "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"github.com/pkg/errors"
)

type pandoc struct {
	cli *docker.Client

	outputDir string
	webFiles  []string
}

const (
	dockerPandocImage = "jagregory/pandoc"
)

var (
	currentTime       = time.Now()
	outputPandocEpub  = "/source/fogo_pandoc_" + currentTime.Format("20060102150405") + ".epub"
	outputTarFileName = "fogo_generated_epub_" + currentTime.Format("20060102150405") + ".tar.gz"
)

// NewPandocConverter a converter based on pandoc application
func NewPandocConverter() (Converter, error) {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "error trying to access to local docker daemon")
	}

	return &pandoc{cli: cli}, nil
}

func (p *pandoc) WithOutputDir(outputDir string) error {
	p.outputDir = outputDir
	return nil
}

func (p *pandoc) WithWebFiles(webFiles []string) error {
	if len(webFiles) == 0 {
		return errors.New("you must specify almost one file")
	}
	p.webFiles = webFiles
	return nil
}

func (p *pandoc) Convert() error {

	cmd := []string{"-o", outputPandocEpub, "-s", "-r", "html"}
	cmd = append(cmd, p.webFiles...)

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
	archive, _, err := p.cli.CopyFromContainer(context.Background(), containerID, outputPandocEpub)
	if err != nil {
		return errors.Wrapf(err, "final file %s is not created", outputPandocEpub)
	}

	content, err := ioutil.ReadAll(archive)
	if err != nil {
		return errors.Wrap(err, "the output file can not be created")
	}

	ioutil.WriteFile(p.outputDir+outputTarFileName, content, 0644)
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
