package api

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"golang.org/x/sync/errgroup"
)

type Files struct {
	client *Client
}

type File struct {
	ID          string
	Name        string
	ContentHash string
}

func (c *Client) Files() *Files { return &Files{client: c} }

func (f *Files) List(viewName string) ([]File, error) {
	var query struct {
		SearchDomain struct {
			Files []File
		} `graphql:"searchDomain(name:$viewName)"`
	}

	variables := map[string]interface{}{
		"viewName": viewName,
	}

	err := f.client.Query(&query, variables)
	return query.SearchDomain.Files, err
}

func (f *Files) Delete(viewName string, fileName string) error {
	var query struct {
		RemoveFile struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"removeFile(name:$viewName, fileName: $fileName)"`
	}

	variables := map[string]interface{}{
		"viewName": viewName,
		"fileName": fileName,
	}

	return f.client.Mutate(&query, variables)
}

func (f *Files) Upload(viewName string, fileName string, reader io.Reader) error {
	pr, pw := io.Pipe()

	multipartWriter := multipart.NewWriter(pw)

	var resp *http.Response

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		var err error
		resp, err = f.client.HTTPRequestContext(ctx, http.MethodPost, fmt.Sprintf("api/v1/dataspaces/%s/files", url.PathEscape(viewName)), pr, multipartWriter.FormDataContentType())
		return err
	})

	eg.Go(func() error {
		defer pw.Close()

		file, err := multipartWriter.CreateFormFile("file", fileName)
		if err != nil {
			return err
		}
		_, err = io.Copy(file, reader)
		if err != nil {
			return err
		}

		return multipartWriter.Close()
	})

	err := eg.Wait()
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server responded with %s: %s", resp.Status, string(body))
	}

	return nil
}

func (f *Files) Download(viewName string, fileName string) (io.Reader, error) {
	resp, err := f.client.HTTPRequest(http.MethodGet, fmt.Sprintf("api/v1/dataspaces/%s/files/%s", url.PathEscape(viewName), url.PathEscape(fileName)), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server responded with %s: %s", resp.Status, string(body))
	}

	return resp.Body, nil
}
