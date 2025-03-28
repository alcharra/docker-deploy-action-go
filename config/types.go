package config

type DeployConfig struct {
	SSHHost             string
	SSHPort             string
	SSHUser             string
	SSHKey              string
	SSHKeyPassphrase    string
	SSHKnownHosts       string
	Fingerprint         string
	Timeout             string
	ProjectPath         string
	DeployFile          string
	ExtraFiles          []string
	Mode                string
	StackName           string
	ComposePull         bool
	DockerNetwork       string
	DockerNetworkDriver string
	DockerNetworkAttach bool
	DockerPrune         string
	RegistryHost        string
	RegistryUser        string
	RegistryPass        string
	EnableRollback      bool
}
