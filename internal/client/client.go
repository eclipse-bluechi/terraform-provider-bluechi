package client

type Client interface {
	Connect() error
	Disconnect() error

	InstallBlueChi(bool, bool) error

	CreateControllerConfig(string, BlueChiControllerConfig) error
	RemoveControllerConfig(string) error
	RestartBlueChiController() error
	StopBlueChiController() error

	CreateAgentConfig(string, BlueChiAgentConfig) error
	RemoveAgentConfig(string) error
	RestartBlueChiAgent() error
	StopBlueChiAgent() error
}
