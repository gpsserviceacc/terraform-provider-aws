// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package globalaccelerator

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/globalaccelerator"
	"github.com/aws/aws-sdk-go/service/globalaccelerator/globalacceleratoriface"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/logging"
	tftags "terraform-provider-awsgps/internal/tags"
	"terraform-provider-awsgps/internal/types/option"
	"terraform-provider-awsgps/names"
)

// listTags lists globalaccelerator service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func listTags(ctx context.Context, conn globalacceleratoriface.GlobalAcceleratorAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &globalaccelerator.ListTagsForResourceInput{
		ResourceArn: aws.String(identifier),
	}

	output, err := conn.ListTagsForResourceWithContext(ctx, input)

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return KeyValueTags(ctx, output.Tags), nil
}

// ListTags lists globalaccelerator service tags and set them in Context.
// It is called from outside this package.
func (p *servicePackage) ListTags(ctx context.Context, meta any, identifier string) error {
	tags, err := listTags(ctx, meta.(*conns.AWSClient).GlobalAcceleratorConn(ctx), identifier)

	if err != nil {
		return err
	}

	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(tags)
	}

	return nil
}

// []*SERVICE.Tag handling

// Tags returns globalaccelerator service tags.
func Tags(tags tftags.KeyValueTags) []*globalaccelerator.Tag {
	result := make([]*globalaccelerator.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &globalaccelerator.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from globalaccelerator service tags.
func KeyValueTags(ctx context.Context, tags []*globalaccelerator.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsIn returns globalaccelerator service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) []*globalaccelerator.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets globalaccelerator service tags in Context.
func setTagsOut(ctx context.Context, tags []*globalaccelerator.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(KeyValueTags(ctx, tags))
	}
}

// updateTags updates globalaccelerator service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTags(ctx context.Context, conn globalacceleratoriface.GlobalAcceleratorAPI, identifier string, oldTagsMap, newTagsMap any) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.GlobalAccelerator)
	if len(removedTags) > 0 {
		input := &globalaccelerator.UntagResourceInput{
			ResourceArn: aws.String(identifier),
			TagKeys:     aws.StringSlice(removedTags.Keys()),
		}

		_, err := conn.UntagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.GlobalAccelerator)
	if len(updatedTags) > 0 {
		input := &globalaccelerator.TagResourceInput{
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

// UpdateTags updates globalaccelerator service tags.
// It is called from outside this package.
func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).GlobalAcceleratorConn(ctx), identifier, oldTags, newTags)
}
