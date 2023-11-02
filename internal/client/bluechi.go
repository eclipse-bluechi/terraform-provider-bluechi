package client

import (
	"strconv"
	"strings"
)

const (
	BlueChiControllerConfdDirectory = "/etc/bluechi/controller.conf.d/"
	BlueChiAgentConfdDirectory      = "/etc/bluechi/agent.conf.d/"
)

type BlueChiControllerConfig struct {
	AllowedNodeNames []string
	ManagerPort      *int64
	LogLevel         *string
	LogTarget        *string
	LogIsQuiet       *bool
}

func (cfg BlueChiControllerConfig) Serialize() string {
	res := "[bluechi-controller]\n"
	res += "AllowedNodeNames=" + strings.Join(cfg.AllowedNodeNames, ",\n\t")
	res += "\n"
	if cfg.ManagerPort != nil {
		res += "ManagerPort=" + strconv.FormatInt(*cfg.ManagerPort, 10) + "\n"
	}
	if cfg.LogLevel != nil {
		res += "LogLevel=" + *cfg.LogLevel + "\n"
	}
	if cfg.LogTarget != nil {
		res += "LogTarget=" + *cfg.LogTarget + "\n"
	}
	if cfg.LogIsQuiet != nil {
		res += "LogIsQuiet=" + strconv.FormatBool(*cfg.LogIsQuiet) + "\n"
	}

	return res
}

type BlueChiAgentConfig struct {
	NodeName          *string
	ManagerHost       *string
	ManagerPort       *int64
	ManagerAddress    *string
	HeartbeatInterval *int64
	LogLevel          *string
	LogTarget         *string
	LogIsQuiet        *bool
}

func (cfg BlueChiAgentConfig) Serialize() string {
	res := "[bluechi-agent]\n"
	res += "NodeName=" + *cfg.NodeName + "\n"
	res += "ManagerHost=" + *cfg.ManagerHost + "\n"
	res += "ManagerPort=" + strconv.FormatInt(*cfg.ManagerPort, 10) + "\n"
	if cfg.ManagerAddress != nil {
		res += "ManagerAddress=" + *cfg.ManagerAddress + "\n"
	}
	if cfg.HeartbeatInterval != nil {
		res += "HeartbeatInterval=" + strconv.FormatInt(*cfg.HeartbeatInterval, 10) + "\n"
	}
	if cfg.LogLevel != nil {
		res += "LogLevel=" + *cfg.LogLevel + "\n"
	}
	if cfg.LogTarget != nil {
		res += "LogTarget=" + *cfg.LogTarget + "\n"
	}
	if cfg.LogIsQuiet != nil {
		res += "LogIsQuiet=" + strconv.FormatBool(*cfg.LogIsQuiet) + "\n"
	}

	return res
}
