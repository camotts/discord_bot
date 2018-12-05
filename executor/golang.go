package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
)

func ExecuteGoScript(goFile string) (io.Reader, io.Reader, error) {
	ctx := context.Background()
	fmt.Println("Creating docker client")
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to create docker client")
	}

	fmt.Println("Pulling go image")
	_, err = cli.ImagePull(ctx, "golang", types.ImagePullOptions{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to pull go image")
	}
	fmt.Println("Creating go container")
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "golang",
		Cmd:   []string{"go", "run", "/usr/src/problem/main.go"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: goFile,
				Target: "/usr/src/problem/main.go",
			},
		},
	}, nil, "")
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, nil, errors.Wrap(err, "Unable to start the container")
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, nil, errors.Wrap(err, "Error while waiting for container")
		}
	case st := <-statusCh:
		fmt.Println(st)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return nil, nil, err
	}
	stdout := make([]byte, 0)
	stderr := make([]byte, 0)
	outRet := bytes.NewBuffer(stdout)
	errRet := bytes.NewBuffer(stderr)
	stdcopy.StdCopy(outRet, errRet, out)
	//io.Copy(os.Stdout, out)
	out.Close()
	fmt.Println("Closed out channel")
	return outRet, errRet, nil
}
