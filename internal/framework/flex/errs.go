// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package flex

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"terraform-provider-awsgps/internal/errs/fwdiag"
)

// must panics if the provided Diagnostics has errors.
func must(diags diag.Diagnostics) {
	fwdiag.Must[any](nil, diags)
}
