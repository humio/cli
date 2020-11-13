package api

import (
	"context"
	"fmt"
	"github.com/shurcooL/graphql"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
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
	var q struct {
		Repository struct {
			Files []File
		} `graphql:"repository(name:$viewName)"`
	}

	variables := map[string]interface{}{
		"viewName": graphql.String(viewName),
	}

	err := f.client.Query(&q, variables)

	return q.Repository.Files, err
}

func (f *Files) Delete(viewName string, fileName string) error {
	var q struct {
		RemoveFile struct {
			TypeName string `graphql:"__typename"`
		} `graphql:"removeFile(name:$viewName, fileName: $fileName)"`
	}

	variables := map[string]interface{}{
		"viewName": graphql.String(viewName),
		"fileName": graphql.String(fileName),
	}

	err := f.client.Mutate(&q, variables)

	return err
}

func (f *Files) Upload(viewName string, fileName string, reader io.Reader) error {
	pr, pw := io.Pipe()

	multipartWriter := multipart.NewWriter(pw)

	var resp *http.Response

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		var err error
		resp, err = f.client.HTTPRequestContext(ctx, http.MethodPost, fmt.Sprintf("/api/v1/dataspaces/%s/files", url.QueryEscape(viewName)), pr, multipartWriter.FormDataContentType())
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
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("server responded with %s: %s", resp.Status, string(body))
	}

	return nil
}

func (f *Files) Download(viewName string, fileName string) (io.Reader, error) {
	resp, err := f.client.HTTPRequest(http.MethodGet, fmt.Sprintf("/api/v1/dataspaces/%s/files/%s", url.QueryEscape(viewName), url.QueryEscape(fileName)), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("server responded with %s: %s", resp.Status, string(body))
	}

	return resp.Body, nil
}
