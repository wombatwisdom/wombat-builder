package library

// Spec represents a git repository containing benthos components.
type Spec struct {
	Name   string `json:"name"`
	Module string `json:"module"`
}

// VersionSpec represents a version of a library. It can refer to a specific tag or branch within a git repository
type VersionSpec struct {
	Name    string       `json:"name"`
	Bundles []BundleSpec `json:"bundles"`
}

// BundleSpec is a grouping of packages within a library.
// This is useful for when a library contains multiple packages that are often used together.
// Common examples are libraries that contain packages with different licenses, or packages that
// want to make a distinction between what is considered stable and what is considered experimental.
type BundleSpec struct {
	Name     string   `json:"name"`
	Packages []string `json:"packages"`
}

// PackageSpec represents a package within a library. It contains the components which can be used when
// including the package.
type PackageSpec struct {
	Name              string   `json:"name"`
	Fqn               string   `json:"fqn"`
	BloblangFunctions []string `json:"bloblang_functions,omitempty"`
	BloblangMethods   []string `json:"bloblang_methods,omitempty"`
	Buffers           []string `json:"buffers,omitempty"`
	Caches            []string `json:"caches,omitempty"`
	Inputs            []string `json:"inputs,omitempty"`
	Metrics           []string `json:"metrics,omitempty"`
	Outputs           []string `json:"outputs,omitempty"`
	Processors        []string `json:"processors,omitempty"`
	RateLimits        []string `json:"rate_limits,omitempty"`
	Scanners          []string `json:"scanners,omitempty"`
}
