package amazon_s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/storage"
)

// GetVisibility returns the visibility of a file
func (adapter Adapter) GetVisibility(ctx context.Context, path string) (entities.Visibility, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return "", storage.Err("create S3 client", err)
	}

	result, err := client.GetObjectAcl(ctx, &s3.GetObjectAclInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return "", storage.PathErr("get visibility", path, err)
	}

	// Check for public-read grant
	for _, grant := range result.Grants {
		if grant.Grantee != nil && grant.Grantee.URI != nil {
			if *grant.Grantee.URI == "http://acs.amazonaws.com/groups/global/AllUsers" {
				return entities.VisibilityPublic, nil
			}
		}
	}

	return entities.VisibilityPrivate, nil
}

// SetVisibility changes the visibility of a file
func (adapter Adapter) SetVisibility(ctx context.Context, path string, visibility entities.Visibility) error {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return storage.Err("create S3 client", err)
	}

	var acl types.ObjectCannedACL
	switch visibility {
	case entities.VisibilityPublic:
		acl = types.ObjectCannedACLPublicRead
	case entities.VisibilityPrivate:
		acl = types.ObjectCannedACLPrivate
	default:
		acl = types.ObjectCannedACLPrivate
	}

	_, err = client.PutObjectAcl(ctx, &s3.PutObjectAclInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(path),
		ACL:    acl,
	})
	if err != nil {
		return storage.PathErr("set visibility", path, err)
	}

	return nil
}
