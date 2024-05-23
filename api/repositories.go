package api

import (
	"fmt"
	"strings"

	graphql "github.com/cli/shurcooL-graphql"
	"github.com/humio/cli/api/internal/humiographql"
)

type Repositories struct {
	client *Client
}

type Repository struct {
	ID                       string
	Name                     string
	Description              string
	RetentionDays            float64                      `graphql:"timeBasedRetention"`
	IngestRetentionSizeGB    float64                      `graphql:"ingestSizeBasedRetention"`
	StorageRetentionSizeGB   float64                      `graphql:"storageSizeBasedRetention"`
	SpaceUsed                int64                        `graphql:"compressedByteSize"`
	S3ArchivingConfiguration humiographql.S3Configuration `graphql:"s3ArchivingConfiguration"`
}

func (c *Client) Repositories() *Repositories { return &Repositories{client: c} }

func (r *Repositories) Get(name string) (Repository, error) {
	var query struct {
		Repository Repository `graphql:"repository(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	err := r.client.Query(&query, variables)

	if err != nil {
		// The graphql error message is vague if the repo already exists, so add a hint.
		return query.Repository, fmt.Errorf("%w. Does the repo already exist?", err)
	}

	return query.Repository, nil
}

type RepoListItem struct {
	ID        string
	Name      string
	SpaceUsed int64 `graphql:"compressedByteSize"`
}

func (r *Repositories) List() ([]RepoListItem, error) {
	var query struct {
		Repositories []RepoListItem `graphql:"repositories"`
	}

	err := r.client.Query(&query, nil)
	return query.Repositories, err
}

func (r *Repositories) Create(name string) error {
	var mutation struct {
		CreateRepository struct {
			Repository Repository
		} `graphql:"createRepository(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	err := r.client.Mutate(&mutation, variables)
	if err != nil {
		// The graphql error message is vague if the repo already exists, so add a hint.
		return fmt.Errorf("%w. Does the repo already exist?", err)
	}

	return nil
}

func (r *Repositories) Delete(name, reason string, allowDataDeletion bool) error {
	if !allowDataDeletion {
		return fmt.Errorf("repository may contain data and data deletion not enabled")
	}

	var mutation struct {
		DeleteSearchDomain struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"deleteSearchDomain(name: $name, deleteMessage: $reason)"`
	}
	variables := map[string]interface{}{
		"name":   graphql.String(name),
		"reason": graphql.String(reason),
	}

	return r.client.Mutate(&mutation, variables)
}

type DefaultGroupEnum string

const (
	DefaultGroupEnumMember     DefaultGroupEnum = "Member"
	DefaultGroupEnumAdmin      DefaultGroupEnum = "Admin"
	DefaultGroupEnumEliminator DefaultGroupEnum = "Eliminator"
)

func (e DefaultGroupEnum) String() string {
	return string(e)
}

func (e *DefaultGroupEnum) ParseString(s string) bool {
	switch strings.ToLower(s) {
	case "member":
		*e = DefaultGroupEnumMember
		return true
	case "admin":
		*e = DefaultGroupEnumAdmin
		return true
	case "eliminator":
		*e = DefaultGroupEnumEliminator
		return true
	default:
		return false
	}
}

func (r *Repositories) UpdateUserGroup(name, username string, groups ...DefaultGroupEnum) error {
	if len(groups) == 0 {
		return fmt.Errorf("at least one group must be defined")
	}

	var mutation struct {
		UpdateDefaultGroupMembershipsMutation struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateDefaultGroupMemberships(input: {viewName: $name, userName: $username, groups: $groups})"`
	}
	variables := map[string]interface{}{
		"name":     graphql.String(name),
		"username": graphql.String(username),
		"groups":   groups,
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) UpdateTimeBasedRetention(name string, retentionInDays float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}

	var mutation struct {
		UpdateRetention struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateRetention(repositoryName: $name, timeBasedRetention: $retentionInDays)"`
	}
	variables := map[string]interface{}{
		"name":            graphql.String(name),
		"retentionInDays": (*graphql.Float)(nil),
	}
	if retentionInDays > 0 {
		if retentionInDays < existingRepo.RetentionDays || existingRepo.RetentionDays == 0 {
			if !allowDataDeletion {
				return fmt.Errorf("repository may contain data and data deletion not enabled")
			}
		}
		variables["retentionInDays"] = graphql.Float(retentionInDays)
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) UpdateStorageBasedRetention(name string, storageInGB float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}

	var mutation struct {
		UpdateRetention struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateRetention(repositoryName: $name, storageSizeBasedRetention: $storageInGB)"`
	}
	variables := map[string]interface{}{
		"name":        graphql.String(name),
		"storageInGB": (*graphql.Float)(nil),
	}
	if storageInGB > 0 {
		if storageInGB < existingRepo.StorageRetentionSizeGB || existingRepo.StorageRetentionSizeGB == 0 {
			if !allowDataDeletion {
				return fmt.Errorf("repository may contain data and data deletion not enabled")
			}
		}
		variables["storageInGB"] = graphql.Float(storageInGB)
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) UpdateIngestBasedRetention(name string, ingestInGB float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}

	var mutation struct {
		UpdateRetention struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateRetention(repositoryName: $name, ingestSizeBasedRetention: $ingestInGB)"`
	}
	variables := map[string]interface{}{
		"name":       graphql.String(name),
		"ingestInGB": (*graphql.Float)(nil),
	}
	if ingestInGB > 0 {
		if ingestInGB < existingRepo.IngestRetentionSizeGB || existingRepo.IngestRetentionSizeGB == 0 {
			if !allowDataDeletion {
				return fmt.Errorf("repository may contain data and data deletion not enabled")
			}
		}
		variables["ingestInGB"] = graphql.Float(ingestInGB)
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) UpdateDescription(name, description string) error {
	var mutation struct {
		UpdateDescription struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateDescriptionForSearchDomain(name: $name, newDescription: $description)"`
	}

	variables := map[string]interface{}{
		"name":        graphql.String(name),
		"description": graphql.String(description),
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) EnableS3Archiving(name string) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}

	if existingRepo.S3ArchivingConfiguration.IsConfigured() == false {
		return fmt.Errorf("repository has no configuration for S3 archiving")
	}

	var mutation struct {
		S3EnableArchiving struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"s3EnableArchiving(repositoryName: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) DisableS3Archiving(name string) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}

	if existingRepo.S3ArchivingConfiguration.IsConfigured() == false {
		return fmt.Errorf("repository has no configuration for S3 archiving")
	}

	var mutation struct {
		S3DisableArchiving struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"s3DisableArchiving(repositoryName: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	return r.client.Mutate(&mutation, variables)
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

	archivingFormat, ferr := humiographql.NewS3ArchivingFormat(format)
	if ferr != nil {
		return ferr
	}

	var mutation struct {
		S3ConfigureArchiving struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"s3ConfigureArchiving(repositoryName: $name, bucket: $bucket, region: $region, format: $format)"`
	}

	variables := map[string]interface{}{
		"name":   graphql.String(name),
		"bucket": graphql.String(bucket),
		"region": graphql.String(region),
		"format": archivingFormat,
	}

	return r.client.Mutate(&mutation, variables)
}
