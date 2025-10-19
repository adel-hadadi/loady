package utils

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
)

func GetContainerAddress(c container.InspectResponse) string {

	for _, net := range c.NetworkSettings.Networks {
		if net != nil && net.IPAddress != "" {
			for port, bindings := range c.NetworkSettings.Ports {
				if len(bindings) > 0 {
					return fmt.Sprintf("%s:%s", net.IPAddress, port.Port())
				}
			}

			return fmt.Sprintf("%s:80", net.IPAddress)
		}
	}

	for _, bindings := range c.NetworkSettings.Ports {
		if len(bindings) > 0 {
			host := bindings[0].HostIP
			if host == "" || host == "0.0.0.0" {
				host = "localhost"
			}
			return fmt.Sprintf("%s:%s", host, bindings[0].HostPort)
		}
	}

	return ""
}
