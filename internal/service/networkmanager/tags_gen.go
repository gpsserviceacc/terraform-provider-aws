// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package networkmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/networkmanager"
	"github.com/aws/aws-sdk-go/service/networkmanager/networkmanageriface"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/logging"
	tftags "terraform-provider-awsgps/internal/tags"
	"terraform-provider-awsgps/internal/types/option"
	"terraform-provider-awsgps/names"
)

// []*SERVICE.Tag handling

// Tags returns networkmanager service tags.
func Tags(tags tftags.KeyValueTags) []*networkmanager.Tag {
	result := make([]*networkmanager.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &networkmanager.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from networkmanager service tags.
func KeyValueTags(ctx context.Context, tags []*networkmanager.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsIn returns networkmanager service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) []*networkmanager.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets networkmanager service tags in Context.
func setTagsOut(ctx context.Context, tags []*networkmanager.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(KeyValueTags(ctx, tags))
	}
}

// updateTags updates networkmanager service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTags(ctx context.Context, conn networkmanageriface.NetworkManagerAPI, identifier string, oldTagsMap, newTagsMap any) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.NetworkManager)
	if len(removedTags) > 0 {
		input := &networkmanager.UntagResourceInput{
			ResourceArn: aws.String(identifier),
			TagKeys:     aws.StringSlice(removedTags.Keys()),
		}

		_, err := conn.UntagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.NetworkManager)
	if len(updatedTags) > 0 {
		input := &networkmanager.TagResourceInput{
			ResourceArn: aws.String(identifier),
			Tags:        Tags(updatedTags),
		}

		_, err := conn.TagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

// UpdateTags updates networkmanager service tags.
// It is called from outside this package.
func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).NetworkManagerConn(ctx), identifier, oldTags, newTags)
}
