package config

import (
	"fmt"
	"regexp"
)

type Config struct {
	// App general configuration
	App App `yaml:"app"`

	// Upgrader configuration
	Upgrader Upgrader `yaml:"upgrader"`

	// PvcMonitor configuration
	PvcMonitor PvcMonitor `yaml:"pvcMonitor"`

	// MaxRecentLogLines optional max tail log lines in messages,
	// if it's not provided it will get all log lines
	MaxRecentLogLines int64 `yaml:"maxRecentLogLines"`

	// IgnoreFailedGracefulShutdown if set to true, containers which are
	// forcefully killed during shutdown (as their graceful shutdown failed)
	// are not reported as error
	IgnoreFailedGracefulShutdown bool `yaml:"ignoreFailedGracefulShutdown"`

	// Namespaces is an optional list of namespaces that you want to watch or
	// forbid, if it's not provided it will watch all namespaces.
	// If you want to forbid a namespace, configure it with !<namespace name>
	// You can either set forbidden namespaces or allowed, not both
	Namespaces []string `yaml:"namespaces"`

	// Reasons is an  optional list of reasons that you want to watch or forbid,
	// if it's not provided it will watch all reasons.
	// If you want to forbid a reason, configure it with !<reason>
	// You can either set forbidden reasons or allowed, not both
	Reasons []string `yaml:"reasons"`

	// IgnoreContainerNames optional list of container names to ignore
	IgnoreContainerNames []string `yaml:"ignoreContainerNames"`

	// IgnorePodLabels is an optional list of labels to help exclude pods
	IgnorePodLabels []IgnorePodLabelRule `yaml:"ignorePodLabels"`

	// Alert is a map contains a map of each provider configuration
	// e.g. {"slack": {"webhook": "URL"}}
	Alert map[string]map[string]interface{} `yaml:"alert"`

	// AllowedNamespaces, ForbiddenNamespaces are calculated internally
	// after loading Namespaces configuration
	AllowedNamespaces   []string
	ForbiddenNamespaces []string

	// AllowedReasons, ForbiddenReasons are calculated internally after loading
	// Reasons configuration
	AllowedReasons   []string
	ForbiddenReasons []string
}

// App confing struct
type App struct {
	// ProxyURL to be used in outgoing http(s) requests except Kubernetes
	// requests to cluster
	ProxyURL string `yaml:"proxyURL"`

	// ClusterName to used in notifications to indicate which cluster has
	// issue
	ClusterName string `yaml:"clusterName"`

	// DisableUpdateCheck if set to true, welcome message will not be
	// sent to configured notification channels
	DisableStartupMessage bool `yaml:"disableStartupMessage"`
}

// Upgrader confing struct
type Upgrader struct {
	// DisableUpdateCheck if set to true, does not check for and
	// notify about kwatch updates
	DisableUpdateCheck bool `yaml:"disableUpdateCheck"`
}

// PvcMonitor confing struct
type PvcMonitor struct {
	// Enabled if set to true, it will check pvc usage periodically
	// By default, this value is true
	Enabled bool `yaml:"enabled"`

	// Interval is the frequency (in minutes) to check pvc usage in the cluster
	// By default, this value is 5
	Interval int `yaml:"interval"`

	// Threshold is the percentage of accepted pvc usage. if current usage
	// exceeds this value, it will send a notification.
	// By default, this value is 80
	Threshold float64 `yaml:"threshold"`
}

// IgnorePodLabelRule config struct
type IgnorePodLabelRule struct {
	// Label is the value of the label to inspect.
	Label string `yaml:"label"`

	// Value is an exact string to match to identify pods to ignore.
	Value string `yaml:"value"`

	// ValueRegex is a regular expression to use to identify pods to ignore. Takes precedence over Value if both are supplied.
	ValueRegex string `yaml:"valueRegex"`

	matcher *regexp.Regexp
}

// IsValid determines whether the rule appears to be well-formed.
func (r *IgnorePodLabelRule) IsValid() error {
	if len(r.Label) == 0 {
		return fmt.Errorf("No label supplied: %#v", r)
	}

	if len(r.Value) == 0 && len(r.ValueRegex) == 0 {
		return fmt.Errorf("No value or valueRegex supplied: %#v", r)
	}

	if len(r.ValueRegex) > 0 {
		if r.Matcher() == nil {
			return fmt.Errorf("Invalid regex: %q", r.ValueRegex)
		}
	}

	return nil
}

// Matcher returns a compiled regular expression to use for matching if the rule is configured as such
func (r *IgnorePodLabelRule) Matcher() *regexp.Regexp {
	if len(r.ValueRegex) > 0 && r.matcher == nil {
		r.matcher = regexp.MustCompile(r.ValueRegex)
	}

	return r.matcher
}
