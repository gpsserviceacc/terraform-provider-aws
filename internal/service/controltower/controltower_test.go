// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package controltower_test

import (
	"testing"

	"terraform-provider-awsgps/internal/acctest"
)

func TestAccControlTower_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"LandingZone": {
			"basic":      testAccLandingZone_basic,
			"disappears": testAccLandingZone_disappears,
			"tags":       testAccLandingZone_tags,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}
