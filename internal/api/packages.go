package api

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/humio/cli/internal/api/humiographql"
)

// Packages is a API client for working with Humio packages.
type Packages struct {
	client *Client
}

// Packages constructs a Packages API client.
func (c *Client) Packages() *Packages { return &Packages{client: c} }

// ValidationResponse contain the results of a package validation.
type ValidationResponse struct {
	InstallationErrors []string `json:"installationErrors"`
	ParseErrors        []string `json:"parseErrors"`
}

// IsValid returns true if there are no errors in the package
func (resp *ValidationResponse) IsValid() bool {
	return (len(resp.InstallationErrors) == 0) && (len(resp.ParseErrors) == 0)
}

// InstalledPackage contain the details of an installed package
type InstalledPackage struct {
	ID          string
	InstalledBy struct {
		Username  string
		Timestamp time.Time
	}
	UpdatedBy struct {
		Username  string
		Timestamp time.Time
	}
	Source          string
	AvailableUpdate *string
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

// Validate checks a package declaration validity against a Humio
// server.
func (p *Packages) Validate(viewName string, absPath string) (*ValidationResponse, error) {
	var zipFilePath string
	var err error

	isDir, err := isDirectory(absPath)
	if err != nil {
		return nil, err
	}

	if isDir {
		zipFilePath, err = createTempZipFromFolder(absPath)
	} else {
		zipFilePath = absPath
	}
	if err != nil {
		return nil, err
	}

	urlPath := "api/v1/packages/analyze?view=" + url.QueryEscape(viewName)

	// #nosec G304
	fileReader, err := os.Open(zipFilePath)
	if err != nil {
		return nil, err
	}
	// #nosec G307
	defer fileReader.Close()

	response, err := p.client.HTTPRequestContext(context.Background(), "POST", urlPath, fileReader, ZIPContentType)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("failed to get response")
	}
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("bad response. %s", response.Status)
	}

	var report ValidationResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

// ListInstalled returns a list of installed packages
func (p *Packages) ListInstalled(searchDomainName string) ([]InstalledPackage, error) {
	resp, err := humiographql.ListInstalledPackages(context.Background(), p.client, searchDomainName)
	if err != nil {
		return nil, err
	}
	respInstalledPackages := resp.GetSearchDomain().GetInstalledPackages()
	packages := make([]InstalledPackage, len(respInstalledPackages))
	for idx, installedPackage := range respInstalledPackages {
		installedBy := installedPackage.GetInstalledBy()
		updatedBy := installedPackage.GetUpdatedBy()
		packages[idx] = InstalledPackage{
			ID: installedPackage.GetId(),
			InstalledBy: struct {
				Username  string
				Timestamp time.Time
			}{
				Username:  installedBy.GetUsername(),
				Timestamp: installedBy.GetTimestamp(),
			},
			UpdatedBy: struct {
				Username  string
				Timestamp time.Time
			}{
				Username:  updatedBy.GetUsername(),
				Timestamp: updatedBy.GetTimestamp(),
			},
			Source:          string(installedPackage.GetSource()),
			AvailableUpdate: installedPackage.GetAvailableUpdate(),
		}
	}

	return packages, nil
}

// InstallArchive installs a local package (zip file).
func (p *Packages) InstallArchive(viewName string, pathToZip string, queryOwnership string) (*ValidationResponse, error) {
	// #nosec G304
	fileReader, err := os.Open(pathToZip)
	if err != nil {
		return nil, err
	}
	// #nosec G307
	defer fileReader.Close()

	urlPath := "api/v1/packages/install?view=" + url.QueryEscape(viewName) + "&overwrite=true" + "&queryOwnershipType=" + queryOwnership
	response, err := p.client.HTTPRequestContext(context.Background(), "POST", urlPath, fileReader, ZIPContentType)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("failed to get response")
	}
	if response.StatusCode >= 400 {
		return nil, detailedInstallationError(response)
	}

	var report ValidationResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func detailedInstallationError(response *http.Response) error {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	body, err := io.ReadAll(response.Body) // response body is []byte
	if err != nil {
		return fmt.Errorf("the package could not be installed, err=%q", err)
	}

	var result InstallationErrors
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("the package could not be installed and the reason returned could not be parsed, err=%q response=%q", err, body)
	}

	allErrors := append(result.ParseErrors, result.InstallationErrors...)

	formattedErrors := make([]string, len(allErrors))
	for i, v := range allErrors {
		formattedErrors[i] = "\n    - " + v
	}

	return fmt.Errorf("%s", strings.Join(formattedErrors, ""))
}

type InstallationErrors struct {
	InstallationErrors []string `json:"installationErrors"`
	ParseErrors        []string `json:"parseErrors"`
	ResponseType       string   `json:"responseType"`
}

type (
	// VersionedPackageSpecifier is the ID and version of a package, e.g foo/bar@2.0.1
	VersionedPackageSpecifier string
	// UnversionedPackageSpecifier is the ID of a package, e.g foo/bar
	UnversionedPackageSpecifier string
)

// UninstallPackage uninstalls a package by name.
func (p *Packages) UninstallPackage(searchDomainName string, packageID string) error {
	_, err := humiographql.UninstallPackage(context.Background(), p.client, searchDomainName, packageID)
	return err
}

// CreateArchive creates a archive by bundling the files in packageDirPath in a zip file.
func (p *Packages) CreateArchive(packageDirPath string, targetFileName string) error {
	// #nosec G304
	outFile, err := os.Create(targetFileName)
	if err != nil {
		return err
	}
	// #nosec G307
	defer outFile.Close()

	return createZipFromFolder(packageDirPath, outFile)
}

// InstallFromDirectory installs a package from a directory containing the package files.
func (p *Packages) InstallFromDirectory(packageDirPath string, targetRepoOrView string, queryOwnership string) (*ValidationResponse, error) {
	zipFilePath, err := createTempZipFromFolder(packageDirPath)
	if err != nil {
		return nil, err
	}

	// #nosec G304
	zipFile, err := os.Open(zipFilePath)
	if err != nil {
		return nil, err
	}

	// #nosec G307
	defer zipFile.Close()
	defer os.Remove(zipFile.Name())

	return p.InstallArchive(targetRepoOrView, zipFilePath, queryOwnership)
}

func createTempZipFromFolder(baseFolder string) (string, error) {
	tempDir := os.TempDir()
	zipFile, err := os.CreateTemp(tempDir, "humio-package.*.zip")
	if err != nil {
		return "", err
	}
	// #nosec G307
	defer zipFile.Close()

	err = createZipFromFolder(baseFolder, zipFile)
	if err != nil {
		return "", err
	}

	return zipFile.Name(), nil
}

func createZipFromFolder(baseFolder string, outFile *os.File) error {
	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	err := addFiles(w, baseFolder, "")
	if err != nil {
		return err
	}

	return w.Close()
}

func isValidFolderOrFile(name string) bool {
	return !strings.HasPrefix(name, "_") && !strings.HasPrefix(name, ".")
}

func addFiles(w *zip.Writer, basePath string, baseInZip string) error {
	// Open the Directory
	files, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !isValidFolderOrFile(file.Name()) {
			continue
		}

		if !file.IsDir() {
			// #nosec G304
			src, err := os.Open(filepath.Join(basePath, file.Name()))
			if err != nil {
				return err
			}

			// Add some files to the archive.
			dst, err := w.Create(path.Join(baseInZip, file.Name()))
			if err != nil {
				_ = src.Close()
				return err
			}
			_, err = io.Copy(dst, src)
			_ = src.Close()
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			// Drill down
			newBase := path.Join(basePath, file.Name())
			err := addFiles(w, newBase, path.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
