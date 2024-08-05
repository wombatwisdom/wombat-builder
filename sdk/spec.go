package sdk

import "github.com/wombatwisdom/wombat-builder/library"

// Spec contains the specification of the build we want to create
type Spec struct {
	Goos           string                `json:"goos"`
	Goarch         string                `json:"goarch"`
	GoVersion      string                `json:"goversion"`
	BenthosVersion string                `json:"benthosversion"`
	Packages       []library.PackageSpec `json:"packages"`
}
