// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !generate
// +build !generate

package inspector

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/aws/aws-sdk-go/service/inspector/inspectoriface"
	"terraform-provider-awsgps/internal/conns"
	tftags "terraform-provider-awsgps/internal/tags"
	"terraform-provider-awsgps/names"
)

// Custom Inspector Classic tag service update functions using the same format as generated code.

// updateTags updates Inspector Classic resource tags.
// The identifier is the resource ARN.
func updateTags(ctx context.Context, conn inspectoriface.InspectorAPI, identifier string, oldTagsMap, newTagsMap any) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap).IgnoreSystem(names.Inspector)

	if len(newTags) > 0 {
		input := &inspector.SetTagsForResourceInput{
			ResourceArn: aws.String(identifier),
			Tags:        Tags(newTags),
		}

		_, err := conn.SetTagsForResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	} else if len(oldTags) > 0 {
		input := &inspector.SetTagsForResourceInput{
			ResourceArn: aws.String(identifier),
		}

		_, err := conn.SetTagsForResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

func createTags(ctx context.Context, conn inspectoriface.InspectorAPI, identifier string, tags []*inspector.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	return updateTags(ctx, conn, identifier, nil, KeyValueTags(ctx, tags))
}

// UpdateTags updates Inspector Classic service tags.
// It is called from outside this package.
func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).InspectorConn(ctx), identifier, oldTags, newTags)
}
