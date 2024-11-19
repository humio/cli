package api

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/humio/cli/internal/api/humiographql"
	"golang.org/x/sync/errgroup"
)

type Files struct {
	client *Client
}

type File struct {
	Name        string
	ContentHash string
}

func (c *Client) Files() *Files { return &Files{client: c} }

func (f *Files) List(searchDomainName string) ([]File, error) {
	resp, err := humiographql.ListFiles(context.Background(), f.client, searchDomainName)
	if err != nil {
		return nil, err
	}

	respFiles := resp.GetSearchDomain().GetFiles()
	files := make([]File, len(respFiles))
	for i, file := range respFiles {
		nameAndPath := file.GetNameAndPath()
		files[i] = File{
			Name:        nameAndPath.GetName(),
			ContentHash: file.GetContentHash(),
		}
	}
	return files, nil
}

func (f *Files) Delete(searchDomainName string, fileName string) error {
	_, err := humiographql.RemoveFile(context.Background(), f.client, searchDomainName, fileName)
	return err
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

	if resp == nil {
		return fmt.Errorf("failed to get response")
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

	if resp == nil {
		return nil, fmt.Errorf("failed to get response")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server responded with %s: %s", resp.Status, string(body))
	}

	return resp.Body, nil
}
