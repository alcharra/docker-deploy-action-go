package validator

type ComposeFile struct {
	Version  string                        `yaml:"version"`
	Services map[string]ServiceDefinition  `yaml:"services"`
	Networks map[string]*NetworkDefinition `yaml:"networks,omitempty"`
	Volumes  map[string]*VolumeDefinition  `yaml:"volumes,omitempty"`
	Configs  map[string]*ConfigDefinition  `yaml:"configs,omitempty"`
	Secrets  map[string]*SecretDefinition  `yaml:"secrets,omitempty"`
}

type StringMap map[string]string

type ServiceDefinition struct {
	Image      string            `yaml:"image"`
	Build      interface{}       `yaml:"build,omitempty"`
	Env        StringMap         `yaml:"environment,omitempty"`
	Labels     StringMap         `yaml:"labels,omitempty"`
	Configs    []ConfigReference `yaml:"configs,omitempty"`
	Secrets    []SecretReference `yaml:"secrets,omitempty"`
	Deploy     *DeployDefinition `yaml:"deploy,omitempty"`
	Ports      []interface{}     `yaml:"ports,omitempty"`
	Command    interface{}       `yaml:"command,omitempty"`
	Entrypoint interface{}       `yaml:"entrypoint,omitempty"`
	Volumes    []interface{}     `yaml:"volumes,omitempty"`
	Networks   []interface{}     `yaml:"networks,omitempty"`
}

type DeployDefinition struct {
	Replicas  *int               `yaml:"replicas,omitempty"`
	Placement *PlacementStrategy `yaml:"placement,omitempty"`
}

type PlacementStrategy struct {
	Constraints []string `yaml:"constraints,omitempty"`
}

type ConfigReference struct {
	Source string `yaml:"source"`
	Target string `yaml:"target,omitempty"`
}

type SecretReference struct {
	Source string `yaml:"source"`
	Target string `yaml:"target,omitempty"`
}

type NetworkDefinition struct {
	Driver   string      `yaml:"driver,omitempty"`
	External interface{} `yaml:"external,omitempty"`
}

type VolumeDefinition struct {
	Driver   string      `yaml:"driver,omitempty"`
	External interface{} `yaml:"external,omitempty"`
}

type ConfigDefinition struct {
	File     string      `yaml:"file,omitempty"`
	External interface{} `yaml:"external,omitempty"`
	Name     string      `yaml:"name,omitempty"`
}

type SecretDefinition struct {
	File     string      `yaml:"file,omitempty"`
	External interface{} `yaml:"external,omitempty"`
	Name     string      `yaml:"name,omitempty"`
}
