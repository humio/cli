package api

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/shurcooL/graphql"
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
	ID string
}

// Validate checks a package declaration validity against a Humio
// server.
func (p *Packages) Validate(repoOrViewName string, absDiretoryPath string) (*ValidationResponse, error) {
	zipFilePath, err := createTempZipFromFolder(absDiretoryPath)

	if err != nil {
		return nil, err
	}

	urlPath := "api/v1/packages/analyze?view=" + url.QueryEscape(repoOrViewName)

	fileReader, openErr := os.Open(zipFilePath)

	if openErr != nil {
		return nil, openErr
	}
	defer fileReader.Close()

	response, httpErr := p.client.HTTPRequestContext(context.Background(), "POST", urlPath, fileReader, "application/zip")

	if httpErr != nil {
		return nil, httpErr
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("Bad response. %s", response.Status)
	}

	var report ValidationResponse
	decoder := json.NewDecoder(response.Body)
	decodeErr := decoder.Decode(&report)

	if decodeErr != nil {
		return nil, decodeErr
	}

	return &report, nil
}

// ListInstalled returns a list of installed packages
func (p *Packages) ListInstalled(repoOrViewName string) ([]InstalledPackage, error) {
	var q struct {
		Repository struct {
			InstalledPackages []InstalledPackage
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(repoOrViewName),
	}

	graphqlErr := p.client.Query(&q, variables)

	var installedPackages []InstalledPackage
	if graphqlErr == nil {
		installedPackages = q.Repository.InstalledPackages
	}

	return installedPackages, graphqlErr
}

// InstallArchive installs a local package (zip file).
func (p *Packages) InstallArchive(repoOrViewName string, pathToZip string) (*ValidationResponse, error) {

	fileReader, openErr := os.Open(pathToZip)

	if openErr != nil {
		return nil, openErr
	}
	defer fileReader.Close()

	urlPath := "api/v1/packages/install?view=" + url.QueryEscape(repoOrViewName)

	response, httpErr := p.client.HTTPRequestContext(context.Background(), "POST", urlPath, fileReader, "application/zip")

	if httpErr != nil {
		return nil, httpErr
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("Bad response. %s", response.Status)
	}

	var report ValidationResponse
	decoder := json.NewDecoder(response.Body)
	decodeErr := decoder.Decode(&report)

	if decodeErr != nil {
		return nil, decodeErr
	}

	return &report, nil
}

type (
	// VersionedPackageSpecifier is the ID and version of a package, e.g foo/bar@2.0.1
	VersionedPackageSpecifier string
	// UnversionedPackageSpecifier is the ID of a package, e.g foo/bar
	UnversionedPackageSpecifier string
)

// UninstallPackage uninstalls a package by name.
func (p *Packages) UninstallPackage(repoOrViewName string, packageID string) error {

	var m struct {
		StartDataRedistribution struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"uninstallPackage(packageId: $packageId, viewName: $viewName)"`
	}

	variables := map[string]interface{}{
		"packageId": UnversionedPackageSpecifier(packageID),
		"viewName":  graphql.String(repoOrViewName),
	}

	graphqlErr := p.client.Mutate(&m, variables)

	return graphqlErr
}

// CreateArchive creates a archive by bundling the files in packageDirPath in a zip file.
func (p *Packages) CreateArchive(packageDirPath string, targetFileName string) error {

	outFile, err := os.Create(targetFileName)

	if err != nil {
		return err
	}
	defer outFile.Close()

	return createZipFromFolder(packageDirPath, outFile)
}

// InstallFromDirectory installs a package from a directory containing the package files.
func (p *Packages) InstallFromDirectory(packageDirPath string, targetRepoOrView string) (*ValidationResponse, error) {
	zipFilePath, err := createTempZipFromFolder(packageDirPath)

	if err != nil {
		return nil, err
	}

	zipFile, err := os.Open(zipFilePath)

	if err != nil {
		return nil, err
	}

	defer zipFile.Close()
	defer os.Remove(zipFile.Name())

	if err != nil {
		return nil, err
	}

	return p.InstallArchive(targetRepoOrView, zipFilePath)
}

func createTempZipFromFolder(baseFolder string) (string, error) {
	tempDir := os.TempDir()
	zipFile, err := ioutil.TempFile(tempDir, "humio-package.*.zip")

	if err != nil {
		return "", err
	}

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
	addFiles(w, baseFolder, "")

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		return err
	}
	return nil
}

func addFiles(w *zip.Writer, basePath string, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(path.Join(basePath, file.Name()))
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(path.Join(baseInZip, file.Name()))
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {
			// Drill down
			newBase := path.Join(basePath, file.Name())
			addFiles(w, newBase, path.Join(baseInZip, file.Name()))
		}
	}
}
