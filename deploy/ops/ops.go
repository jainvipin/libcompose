
package ops

import (
	"os"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"github.com/docker/go-connections/nat"
	log "github.com/Sirupsen/logrus"
)

const (
	cfgFile = "src/github.com/docker/libcompose/deploy/example/ops.json"
)

type UserPolicyInfo struct {
	User string
	Networks string
	NetworkPolicies string
	DefaultNetworkPolicy string
}

type NetworkPolicyInfo struct {
	Name string
	Rules []string
}

type opsPolicy struct {
	UserPolicy []UserPolicyInfo
	NetworkPolicy []NetworkPolicyInfo
}

var ops opsPolicy

func LoadOps() error {
	composeFile := os.Getenv("GOPATH") + "/" + cfgFile
	composeFile = "./ops.json"
	return loadOpsWithFile(composeFile)
}

func loadOpsWithFile(fileName string) error {

	composeBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("error reading the config file: %s", err)
	}

	if err := json.Unmarshal(composeBytes, &ops); err != nil {
		log.Errorf("error unmarshaling json %#v \n", err)
		return err
	}

	return nil
}

func UserOpsCheckNetwork(userName, network string) error {
	for _, policy := range ops.UserPolicy {
		if policy.User != userName {
			continue
		}
		allowedNetworks := strings.Split(policy.Networks, ",")
		for _, allowedNetwork := range allowedNetworks {
			if allowedNetwork == network || allowedNetwork == "all" {
				return nil
			}
		}
	}

	return errors.New("Deny unspecified user")
}

func UserOpsGetDefaultNetworkPolicy(userName string) (string, error) {
	for _, policy := range ops.UserPolicy {
		if policy.User != userName {
			continue
		}
		if policy.DefaultNetworkPolicy != "" {
			return policy.DefaultNetworkPolicy, nil
		}
	}

	return "", errors.New("Default Policy Not Found")
}

func UserOpsCheckNetworkPolicy(userName, netPolicy string) error {
	for _, policy := range ops.UserPolicy {
		if policy.User != userName {
			continue
		}
		allowedNetPolicies := strings.Split(policy.NetworkPolicies, ",")
		for _, allowedNetPolicy := range allowedNetPolicies {
			if allowedNetPolicy == netPolicy || allowedNetPolicy == "all" {
				return nil
			}
		}
	}

	return errors.New("Deny unspecified user")
}

func GetRules(policyName string) ([]nat.Port, error) {
	portList := []nat.Port{}

	for _, policy := range ops.NetworkPolicy {
		if policy.Name != policyName {
			continue
		}
		for _, rule := range policy.Rules {
			var natPort nat.Port
			var err error

			clauses := strings.Split(rule, " ")
			if len(clauses) <= 0 {
				return portList, errors.New("none found")
			}
			switch clauses[0] {
				case "permit":
					if len(clauses) <= 1 {
						return portList, errors.New("Incomplete permit clause")
					}
					protoPort := strings.Split(clauses[1], "/")
					if len(protoPort) == 0 {
						return portList, errors.New("Empty proto/port in permit clause")
					}
					switch protoPort[0] {
						case "tcp", "udp":
							if len(protoPort) <= 1 {
								return portList, errors.New("Invalid permit clause: port or protocol missing")
							}
							pNum, _ := strconv.Atoi(protoPort[1])
							if pNum < 0 || pNum > 65535 {
								return portList, errors.New("Invalid port in permit clause")
							}
							natPort, err = nat.NewPort(protoPort[0], protoPort[1]);
							if err != nil {
								return portList, err
							}
						case "icmp":
							natPort, err = nat.NewPort(protoPort[0], "0");
							if err != nil {
								return portList, err
							}
						case "app":
							natPort, err = nat.NewPort("app", "0")
							if err != nil {
								return portList, err
							}
						case "all":
							natPort, err = nat.NewPort("all", "0")
							if err != nil {
								return portList, err
							}
						default:
							return portList, errors.New("Invalid proto in permit clause")
					}

				case "deny":
					return portList, errors.New("Not supported")

				default:
					return portList, errors.New("Invalid clause")
			}
			portList = append(portList, natPort)
		}
	}

	if len(portList) == 0 {
		return portList, errors.New("Unrecognized policy")
	}

	return portList, nil
}
