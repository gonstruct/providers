package amazon_s3

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gonstruct/providers/entities"
)

func storageVisibilityToS3ACL(visibility entities.Visibility) types.ObjectCannedACL {
	switch visibility {
	case entities.VisibilityPublic:
		return types.ObjectCannedACLPublicRead
	case entities.VisibilityPrivate:
		return types.ObjectCannedACLPrivate
	default:
		return ""
	}
}

func s3ACLToStorageVisibility(acl types.ObjectCannedACL) entities.Visibility {
	switch acl {
	case types.ObjectCannedACLPublicRead:
		return entities.VisibilityPublic
	case types.ObjectCannedACLPrivate:
		return entities.VisibilityPrivate
	default:
		return ""
	}
}
