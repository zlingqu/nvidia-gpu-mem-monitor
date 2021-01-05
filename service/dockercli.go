package service

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ContainersStruct struct {
	ContainerID string
	Image       string
	ConName     string
}

// ListContains 这个方法没有用到
func ListContains(cli *client.Client) []ContainersStruct {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	containerList := []ContainersStruct{}
	for _, container := range containers {
		containerList = append(containerList, ContainersStruct{container.ID[:10], container.Image, container.Names[0]})
	}
	return containerList
}

func GetContainsPodInfo(cli *client.Client, containID string) (string, string) {
	containers, _ := cli.ContainerInspect(context.Background(), containID)
	podName := containers.Config.Labels["io.kubernetes.pod.name"]
	podNamespace := containers.Config.Labels["io.kubernetes.pod.namespace"]
	// return containers.State.Pid
	return podName, podNamespace
}
