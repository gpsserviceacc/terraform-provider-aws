// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package elasticache

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	tftags "terraform-provider-awsgps/internal/tags"
	"terraform-provider-awsgps/internal/types/option"
)

// []*SERVICE.Tag handling

// TagsV2 returns elasticache service tags.
func TagsV2(tags tftags.KeyValueTags) []awstypes.Tag {
	result := make([]awstypes.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := awstypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// keyValueTagsV2 creates tftags.KeyValueTags from elasticache service tags.
func keyValueTagsV2(ctx context.Context, tags []awstypes.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.ToString(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsInV2 returns elasticache service tags from Context.
// nil is returned if there are no input tags.
func getTagsInV2(ctx context.Context) []awstypes.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := TagsV2(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOutV2 sets elasticache service tags in Context.
func setTagsOutV2(ctx context.Context, tags []awstypes.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(keyValueTagsV2(ctx, tags))
	}
}
