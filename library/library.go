package library

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
)

var nameRegex = regexp.MustCompile(`^[a-z0-9_\-]+$`)

type Client interface {
	Libraries() ([]string, error)
	Library(libId string) (*Spec, error)
	AddLibrary(library Spec) error
	Versions(libId string) ([]string, error)
	Version(libId string, version string) (*VersionSpec, error)
	AddVersion(libId string, version VersionSpec) error
	Packages(libId string, version string) ([]string, error)
	Package(libId string, version string, pkgId string) (*PackageSpec, error)
	AddPackage(libId string, version string, pkg PackageSpec) error
	UploadArtifact(libId string, verId string, pkgId string, goos string, goarch string, gover string, data io.ReadCloser) error
	Download(libId string, verId string, pkgId string, goos string, goarch string, gover string) (io.ReadCloser, error)
	UploadDocs(libId string, verId string, pkgId string, doc DocView) error
}

func NewFsClient(basePath string) Client {
	return &fsClient{basePath: basePath}
}

type fsClient struct {
	basePath string
}

func (c *fsClient) AddPackage(libId string, version string, pkg PackageSpec) error {
	// get the library
	_, err := c.Version(libId, version)
	if err != nil {
		return err
	}

	pkgFile := path.Join(c.basePath, libId, version, pkg.Name, "package.json")

	// -- create the lib dir
	if err := os.MkdirAll(path.Dir(pkgFile), 0755); err != nil {
		return fmt.Errorf("failed to create package directory %q: %w", path.Dir(pkgFile), err)
	}

	f, err := os.Create(pkgFile)
	if err != nil {
		return fmt.Errorf("failed to create package file %q: %w", pkgFile, err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(pkg); err != nil {
		return fmt.Errorf("failed to encode package file %q: %w", pkgFile, err)
	}

	return nil
}

func (c *fsClient) AddLibrary(library Spec) error {
	if !nameRegex.MatchString(library.Name) {
		return fmt.Errorf("invalid library name %q. must be lowercase and only contain - or _", library.Name)
	}

	libFile := path.Join(c.basePath, library.Name, "library.json")

	// -- create the lib dir
	if err := os.MkdirAll(path.Dir(libFile), 0755); err != nil {
		return fmt.Errorf("failed to create library directory %q: %w", path.Dir(libFile), err)
	}

	f, err := os.Create(libFile)
	if err != nil {
		return fmt.Errorf("failed to create library file %q: %w", libFile, err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(library); err != nil {
		return fmt.Errorf("failed to encode library file %q: %w", libFile, err)
	}

	return nil
}

func (c *fsClient) AddVersion(libId string, version VersionSpec) error {
	// get the library
	lib, err := c.Library(libId)
	if err != nil {
		return err
	}

	versionFile := path.Join(c.basePath, lib.Name, version.Name, "version.json")

	// -- create the lib dir
	if err := os.MkdirAll(path.Dir(versionFile), 0755); err != nil {
		return fmt.Errorf("failed to create version directory %q: %w", path.Dir(versionFile), err)
	}

	f, err := os.Create(versionFile)
	if err != nil {
		return fmt.Errorf("failed to create version file %q: %w", versionFile, err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(version); err != nil {
		return fmt.Errorf("failed to encode version file %q: %w", versionFile, err)
	}

	return nil
}

func (c *fsClient) Libraries() ([]string, error) {
	libDirs, err := os.ReadDir(c.basePath)
	if err != nil {
		return nil, err
	}

	libs := make([]string, 0, len(libDirs))
	for _, libDir := range libDirs {
		if libDir.IsDir() {
			libs = append(libs, libDir.Name())
		}
	}

	return libs, nil
}

func (c *fsClient) Library(libId string) (*Spec, error) {
	libDir := path.Join(c.basePath, libId)
	if _, err := os.Stat(libDir); err != nil {
		return nil, fmt.Errorf("library %q not found: %w", libId, err)
	}

	libFile := path.Join(libDir, "library.json")
	lf, err := os.Open(libFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open library file %q: %w", libFile, err)
	}

	var lib Spec
	if err := json.NewDecoder(lf).Decode(&lib); err != nil {
		return nil, fmt.Errorf("failed to decode library file %q: %w", libFile, err)
	}

	return &lib, nil
}

func (c *fsClient) Versions(libId string) ([]string, error) {
	libDir := path.Join(c.basePath, libId)
	if _, err := os.Stat(libDir); err != nil {
		return nil, fmt.Errorf("library %q not found: %w", libId, err)
	}

	verDirs, err := os.ReadDir(libDir)
	if err != nil {
		return nil, err
	}

	vers := make([]string, 0, len(verDirs))
	for _, verDir := range verDirs {
		if verDir.IsDir() {
			vers = append(vers, verDir.Name())
		}
	}

	return vers, nil
}

func (c *fsClient) Version(libId string, version string) (*VersionSpec, error) {
	verDir := path.Join(c.basePath, libId, version)
	if _, err := os.Stat(verDir); err != nil {
		return nil, fmt.Errorf("version %q not found: %w", version, err)
	}

	verFile := path.Join(verDir, "version.json")
	vf, err := os.Open(verFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open version file %q: %w", verFile, err)
	}

	var ver VersionSpec
	if err := json.NewDecoder(vf).Decode(&ver); err != nil {
		return nil, fmt.Errorf("failed to decode version file %q: %w", verFile, err)
	}

	return &ver, nil
}

func (c *fsClient) Packages(libId string, version string) ([]string, error) {
	verDir := path.Join(c.basePath, libId, version)
	if _, err := os.Stat(verDir); err != nil {
		return nil, fmt.Errorf("version %q not found: %w", version, err)
	}

	pkgDirs, err := os.ReadDir(verDir)
	if err != nil {
		return nil, err
	}

	pkgs := make([]string, 0, len(pkgDirs))
	for _, pkgDir := range pkgDirs {
		if pkgDir.IsDir() {
			pkgs = append(pkgs, pkgDir.Name())
		}
	}

	return pkgs, nil
}

func (c *fsClient) Package(libId string, version string, pkgId string) (*PackageSpec, error) {
	pkgFile := path.Join(c.basePath, libId, version, pkgId, "package.json")
	pf, err := os.Open(pkgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open package file %q: %w", pkgFile, err)
	}

	var pkg PackageSpec
	if err := json.NewDecoder(pf).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("failed to decode package file %q: %w", pkgFile, err)
	}

	return &pkg, nil
}

func (c *fsClient) UploadArtifact(libId string, verId string, pkgId string, goos string, goarch string, gover string, data io.ReadCloser) error {
	pkgDir := c.PackagePath(libId, verId, pkgId)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return fmt.Errorf("failed to create package directory %q: %w", pkgDir, err)
	}

	artifactFile := path.Join(pkgDir, fmt.Sprintf("%s_%s_%s_%s_%s_go%s.so", libId, pkgId, verId, goos, goarch, gover))

	af, err := os.Create(artifactFile)
	if err != nil {
		return fmt.Errorf("failed to create artifact file %q: %w", artifactFile, err)
	}
	defer af.Close()

	if _, err := io.Copy(af, data); err != nil {
		return fmt.Errorf("failed to write artifact file %q: %w", artifactFile, err)
	}

	return nil
}

func (c *fsClient) Download(libId string, verId string, pkgId string, goos string, goarch string, gover string) (io.ReadCloser, error) {
	pkgDir := c.PackagePath(libId, verId, pkgId)
	artifactFile := path.Join(pkgDir, fmt.Sprintf("%s_%s_%s_%s_%s_go%s.so", libId, pkgId, verId, goos, goarch, gover))

	af, err := os.Open(artifactFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open artifact file %q: %w", artifactFile, err)
	}

	return af, nil
}

func (c *fsClient) UploadDocs(libId string, verId string, pkgId string, doc DocView) error {
	pkgDir := c.PackagePath(libId, verId, pkgId)
	docFile := path.Join(pkgDir, doc.Kind, fmt.Sprintf("%s.json", doc.Name))

	if err := os.MkdirAll(path.Dir(docFile), 0755); err != nil {
		return fmt.Errorf("failed to create docfile directory %q: %w", pkgDir, err)
	}

	df, err := os.Create(docFile)
	if err != nil {
		return err
	}
	defer df.Close()

	if err := json.NewEncoder(df).Encode(doc); err != nil {
		return err
	}

	return nil
}

func (c *fsClient) PackagePath(libId string, verId string, pkgId string) string {
	return path.Join(c.basePath, libId, verId, pkgId)
}
