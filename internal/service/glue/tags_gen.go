// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package glue

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/aws/aws-sdk-go/service/glue/glueiface"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/logging"
	tftags "terraform-provider-awsgps/internal/tags"
	"terraform-provider-awsgps/internal/types/option"
	"terraform-provider-awsgps/names"
)

// listTags lists glue service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func listTags(ctx context.Context, conn glueiface.GlueAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &glue.GetTagsInput{
		ResourceArn: aws.String(identifier),
	}

	output, err := conn.GetTagsWithContext(ctx, input)

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return KeyValueTags(ctx, output.Tags), nil
}

// ListTags lists glue service tags and set them in Context.
// It is called from outside this package.
func (p *servicePackage) ListTags(ctx context.Context, meta any, identifier string) error {
	tags, err := listTags(ctx, meta.(*conns.AWSClient).GlueConn(ctx), identifier)

	if err != nil {
		return err
	}

	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(tags)
	}

	return nil
}

// map[string]*string handling

// Tags returns glue service tags.
func Tags(tags tftags.KeyValueTags) map[string]*string {
	return aws.StringMap(tags.Map())
}

// KeyValueTags creates tftags.KeyValueTags from glue service tags.
func KeyValueTags(ctx context.Context, tags map[string]*string) tftags.KeyValueTags {
	return tftags.New(ctx, tags)
}

// getTagsIn returns glue service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) map[string]*string {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets glue service tags in Context.
func setTagsOut(ctx context.Context, tags map[string]*string) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(KeyValueTags(ctx, tags))
	}
}

// updateTags updates glue service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTags(ctx context.Context, conn glueiface.GlueAPI, identifier string, oldTagsMap, newTagsMap any) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.Glue)
	if len(removedTags) > 0 {
		input := &glue.UntagResourceInput{
			ResourceArn:  aws.String(identifier),
			TagsToRemove: aws.StringSlice(removedTags.Keys()),
		}

		_, err := conn.UntagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.Glue)
	if len(updatedTags) > 0 {
		input := &glue.TagResourceInput{
			ResourceArn: aws.String(identifier),
			TagsToAdd:   Tags(updatedTags),
		}

		_, err := conn.TagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

// UpdateTags updates glue service tags.
// It is called from outside this package.
func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).GlueConn(ctx), identifier, oldTags, newTags)
}
