package api

import (
	"context"
	"fmt"

	"github.com/humio/cli/internal/api/humiographql"
)

type Repositories struct {
	client *Client
}

type Repository struct {
	ID                       string
	Name                     string
	Description              *string
	RetentionDays            *float64
	IngestRetentionSizeGB    *float64
	StorageRetentionSizeGB   *float64
	SpaceUsed                int64
	S3ArchivingConfiguration S3Configuration
	AutomaticSearch          bool
}

func (c *Client) Repositories() *Repositories { return &Repositories{client: c} }

func (r *Repositories) Get(name string) (Repository, error) {
	getRepositoryResp, err := humiographql.GetRepository(context.Background(), r.client, name)
	if err != nil {
		return Repository{}, RepositoryNotFound(name)
	}

	repository := getRepositoryResp.GetRepository()
	configuration := repository.GetS3ArchivingConfiguration()
	s3ArchivingConfiguration := S3Configuration{}
	if configuration != nil {
		if configuration.GetFormat() == nil {
			return Repository{}, fmt.Errorf("archiving format not defined")
		}
		disabled := configuration.GetDisabled()
		s3ArchivingConfiguration = S3Configuration{
			Bucket:   configuration.GetBucket(),
			Region:   configuration.GetRegion(),
			Disabled: disabled != nil && *disabled,
			Format:   string(*configuration.GetFormat()),
		}
	}

	return Repository{
		ID:                       repository.GetId(),
		Name:                     repository.GetName(),
		Description:              repository.GetDescription(),
		RetentionDays:            repository.GetTimeBasedRetention(),
		IngestRetentionSizeGB:    repository.GetIngestSizeBasedRetention(),
		StorageRetentionSizeGB:   repository.GetStorageSizeBasedRetention(),
		SpaceUsed:                repository.GetCompressedByteSize(),
		S3ArchivingConfiguration: s3ArchivingConfiguration,
		AutomaticSearch:          repository.GetAutomaticSearch(),
	}, nil
}

type RepoListItem struct {
	ID        string
	Name      string
	SpaceUsed int64 `graphql:"compressedByteSize"`
}

func (r *Repositories) List() ([]RepoListItem, error) {
	listRepositories, err := humiographql.ListRepositories(context.Background(), r.client)
	if err != nil {
		return nil, err
	}
	repoList := make([]RepoListItem, len(listRepositories.GetRepositories()))
	for i, x := range listRepositories.GetRepositories() {
		repoList[i] = RepoListItem{
			ID:        x.GetId(),
			Name:      x.GetName(),
			SpaceUsed: x.GetCompressedByteSize(),
		}
	}
	return repoList, nil
}

func (r *Repositories) Create(name string) error {
	_, err := humiographql.CreateRepository(context.Background(), r.client, name)
	return err
}

func (r *Repositories) Delete(name, reason string, allowDataDeletion bool) error {
	_, err := r.Get(name)
	if err != nil {
		return err
	}

	if !allowDataDeletion {
		return fmt.Errorf("repository may contain data and data deletion not enabled")
	}

	_, err = humiographql.DeleteSearchDomain(context.Background(), r.client, name, reason)
	return err
}

func (r *Repositories) UpdateTimeBasedRetention(name string, retentionInDays *float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}

	if retentionInDays != nil && *retentionInDays > 0 {
		if existingRepo.RetentionDays == nil || *retentionInDays < *existingRepo.RetentionDays {
			if !allowDataDeletion {
				return fmt.Errorf("repository may contain data and data deletion not enabled")
			}
		}
	}

	_, err = humiographql.UpdateTimeBasedRetention(context.Background(), r.client, name, retentionInDays)
	return err
}

func (r *Repositories) UpdateStorageBasedRetention(name string, storageInGB *float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}

	if storageInGB != nil && *storageInGB > 0 {
		if existingRepo.StorageRetentionSizeGB == nil || *storageInGB < *existingRepo.StorageRetentionSizeGB {
			if !allowDataDeletion {
				return fmt.Errorf("repository may contain data and data deletion not enabled")
			}
		}
	}

	_, err = humiographql.UpdateStorageBasedRetention(context.Background(), r.client, name, storageInGB)
	return err
}

func (r *Repositories) UpdateIngestBasedRetention(name string, ingestInGB *float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}

	if ingestInGB != nil && *ingestInGB > 0 {
		if existingRepo.IngestRetentionSizeGB == nil || *ingestInGB < *existingRepo.IngestRetentionSizeGB {
			if !allowDataDeletion {
				return fmt.Errorf("repository may contain data and data deletion not enabled")
			}
		}
	}

	_, err = humiographql.UpdateIngestBasedRetention(context.Background(), r.client, name, ingestInGB)
	return err
}

func (r *Repositories) UpdateDescription(name, description string) error {
	_, err := r.Get(name)
	if err != nil {
		return err
	}

	_, err = humiographql.UpdateDescriptionForSearchDomain(context.Background(), r.client, name, description)
	return err
}

type S3Configuration struct {
	Bucket   string `graphql:"bucket"`
	Region   string `graphql:"region"`
	Disabled bool   `graphql:"disabled"`
	Format   string `graphql:"format"`
}

// IsEnabled - determine if S3Configuration is enabled based on values and the Disabled field
// to avoid a bool defaulting to false
func (s *S3Configuration) IsEnabled() bool {
	if !s.IsConfigured() {
		return false
	}
	return !s.Disabled
}

func (s *S3Configuration) IsConfigured() bool {
	if s.Bucket != "" && s.Region != "" && s.Format != "" {
		return true
	}
	return false
}

func (r *Repositories) EnableS3Archiving(name string) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}

	if !existingRepo.S3ArchivingConfiguration.IsConfigured() {
		return fmt.Errorf("repository has no configuration for S3 archiving")
	}

	_, err = humiographql.EnableS3Archiving(context.Background(), r.client, name)
	return err
}

func (r *Repositories) DisableS3Archiving(name string) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}

	if !existingRepo.S3ArchivingConfiguration.IsConfigured() {
		return fmt.Errorf("repository has no configuration for S3 archiving")
	}

	_, err = humiographql.DisableS3Archiving(context.Background(), r.client, name)
	return err
}

func (r *Repositories) UpdateS3ArchivingConfiguration(name string, bucket string, region string, format string) error {
	_, err := r.Get(name)
	if err != nil {
		return err
	}

	if bucket == "" {
		return fmt.Errorf("bucket name cannot have an empty value")
	}

	if region == "" {
		return fmt.Errorf("region cannot have an empty value")
	}

	_, err = humiographql.UpdateS3ArchivingConfiguration(context.Background(), r.client, name, bucket, region, humiographql.S3ArchivingFormat(format))
	return err
}

func (r *Repositories) UpdateAutomaticSearch(name string, automaticSearch bool) error {
	_, err := r.Get(name)
	if err != nil {
		return err
	}

	_, err = humiographql.SetAutomaticSearching(context.Background(), r.client, name, automaticSearch)
	return err
}
