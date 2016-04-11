package nethooks

import (
	"os"
	"os/exec"
	"strings"

	"github.com/docker/go-connections/nat"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"golang.org/x/net/context"
)

type imageInfo struct {
	portID    int
	protoName string
}

func getImageInfo(imageName string) ([]nat.Port, error) {
	imageInfoList := []nat.Port{}

	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		dockerHost = client.DefaultDockerHost
	}

	defaultHeaders := map[string]string{"User-Agent": "deploy-client"}

	log.Infof("docker host = %s ", dockerHost)
	cl, err := client.NewClient(dockerHost, "v1.20", nil, defaultHeaders)
	if err != nil {
		log.Errorf("Unable to connect to docker. Error %v", err)
		return imageInfoList, err
	}

	imageInfo, _, err := cl.ImageInspectWithRaw(context.Background(), imageName, false)
	log.Debugf("Got the following container info %#v", imageInfo.ContainerConfig)

	if err != nil {
		log.Errorf("Unable to inspect image '%s'. Error %v", imageName, err)
		return imageInfoList, err
	}

	for port := range imageInfo.ContainerConfig.ExposedPorts {
		log.Infof("Extrated (protocol, port) = (%s, %s)", port.Proto(), port.Port())
		imageInfoList = append(imageInfoList, port)
	}

	return imageInfoList, nil
}

func getContainerIP(contName string) string {
	ipAddress := ""
	output, err := exec.Command("docker", "exec", contName, "/sbin/ip", "address", "show").CombinedOutput()
	if err != nil {
		log.Errorf("Unable to fetch container '%s' IP. Error %v", contName, err)
		return ipAddress
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "eth0") && strings.Contains(line, "inet ") {
			words := strings.Split(line, " ")
			for idx, word := range words {
				if word == "inet" {
					ipAddress = strings.Split(words[idx+1], "/")[0]
				}
			}
		}
	}

	return ipAddress
}

func populateEtcHosts(contName, dnsSvcName, ipAddress string) error {
	command := "echo " + ipAddress + "     " + dnsSvcName + " >> /etc/hosts"
	if _, err := exec.Command("docker", "exec", contName, "/bin/sh", "-c", command).CombinedOutput(); err != nil {
		log.Errorf("Unable to populate etc hosts. Error %v", err)
		return err
	}

	if output, err := exec.Command("docker", "exec", contName, "cat", "/etc/hosts").CombinedOutput(); err != nil {
		log.Infof("VJ ===> output = %s ", output)
	}
	return nil
}

func containerExists(contName string) bool {
  _, err := exec.Command("docker", "inspect", contName).CombinedOutput()
  return err == nil
}
