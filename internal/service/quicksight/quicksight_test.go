// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package quicksight_test

import (
	"testing"

	"terraform-provider-awsgps/internal/acctest"
)

func TestAccQuickSight_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"AccountSubscription": {
			"basic":      testAccAccountSubscription_basic,
			"disappears": testAccAccountSubscription_disappears,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}
