// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acmpca_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/acmpca"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-awsgps/internal/acctest"
	"terraform-provider-awsgps/internal/conns"
	tfacmpca "terraform-provider-awsgps/internal/service/acmpca"
	"terraform-provider-awsgps/internal/tfresource"
	"terraform-provider-awsgps/names"
)

func TestAccACMPCACertificateAuthorityCertificate_rootCA(t *testing.T) {
	ctx := acctest.Context(t)
	var v acmpca.GetCertificateAuthorityCertificateOutput
	resourceName := "aws_acmpca_certificate_authority_certificate.test"

	commonName := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ACMPCAServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil, // Certificate authority certificates cannot be deleted
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateAuthorityCertificateConfig_rootCA(commonName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCertificateAuthorityCertificateExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_authority_arn", "aws_acmpca_certificate_authority.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate", "aws_acmpca_certificate.test", "certificate"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_chain", "aws_acmpca_certificate.test", "certificate_chain"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccACMPCACertificateAuthorityCertificate_updateRootCA(t *testing.T) {
	ctx := acctest.Context(t)
	var v acmpca.GetCertificateAuthorityCertificateOutput
	resourceName := "aws_acmpca_certificate_authority_certificate.test"
	updatedResourceName := "aws_acmpca_certificate_authority_certificate.updated"

	commonName := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ACMPCAServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil, // Certificate authority certificates cannot be deleted
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateAuthorityCertificateConfig_rootCA(commonName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCertificateAuthorityCertificateExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_authority_arn", "aws_acmpca_certificate_authority.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate", "aws_acmpca_certificate.test", "certificate"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_chain", "aws_acmpca_certificate.test", "certificate_chain"),
				),
			},
			{
				Config: testAccCertificateAuthorityCertificateConfig_updateRootCA(commonName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCertificateAuthorityCertificateExists(ctx, updatedResourceName, &v),
					resource.TestCheckResourceAttrPair(updatedResourceName, "certificate_authority_arn", "aws_acmpca_certificate_authority.test", "arn"),
					resource.TestCheckResourceAttrPair(updatedResourceName, "certificate", "aws_acmpca_certificate.updated", "certificate"),
					resource.TestCheckResourceAttrPair(updatedResourceName, "certificate_chain", "aws_acmpca_certificate.updated", "certificate_chain"),
				),
			},
		},
	})
}

func TestAccACMPCACertificateAuthorityCertificate_subordinateCA(t *testing.T) {
	ctx := acctest.Context(t)
	var v acmpca.GetCertificateAuthorityCertificateOutput
	resourceName := "aws_acmpca_certificate_authority_certificate.test"

	commonName := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ACMPCAServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil, // Certificate authority certificates cannot be deleted
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateAuthorityCertificateConfig_subordinateCA(commonName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCertificateAuthorityCertificateExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_authority_arn", "aws_acmpca_certificate_authority.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate", "aws_acmpca_certificate.test", "certificate"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_chain", "aws_acmpca_certificate.test", "certificate_chain"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCertificateAuthorityCertificateExists(ctx context.Context, resourceName string, certificate *acmpca.GetCertificateAuthorityCertificateOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ACMPCAConn(ctx)

		output, err := tfacmpca.FindCertificateAuthorityCertificateByARN(ctx, conn, rs.Primary.ID)
		if err != nil {
			return err
		}
		if tfresource.NotFound(err) {
			return fmt.Errorf("ACM PCA Certificate (%s) does not exist", rs.Primary.ID)
		}

		*certificate = *output

		return nil
	}
}

func testAccCertificateAuthorityCertificateConfig_rootCA(commonName string) string {
	return fmt.Sprintf(`
resource "aws_acmpca_certificate_authority_certificate" "test" {
  certificate_authority_arn = aws_acmpca_certificate_authority.test.arn

  certificate       = aws_acmpca_certificate.test.certificate
  certificate_chain = aws_acmpca_certificate.test.certificate_chain
}

resource "aws_acmpca_certificate" "test" {
  certificate_authority_arn   = aws_acmpca_certificate_authority.test.arn
  certificate_signing_request = aws_acmpca_certificate_authority.test.certificate_signing_request
  signing_algorithm           = "SHA512WITHRSA"

  template_arn = "arn:${data.aws_partition.current.partition}:acm-pca:::template/RootCACertificate/V1"

  validity {
    type  = "YEARS"
    value = 1
  }
}

resource "aws_acmpca_certificate_authority" "test" {
  permanent_deletion_time_in_days = 7
  type                            = "ROOT"

  certificate_authority_configuration {
    key_algorithm     = "RSA_4096"
    signing_algorithm = "SHA512WITHRSA"

    subject {
      common_name = %[1]q
    }
  }
}

data "aws_partition" "current" {}
`, commonName)
}

func testAccCertificateAuthorityCertificateConfig_updateRootCA(commonName string) string {
	return fmt.Sprintf(`
resource "aws_acmpca_certificate_authority_certificate" "updated" {
  certificate_authority_arn = aws_acmpca_certificate_authority.test.arn

  certificate       = aws_acmpca_certificate.updated.certificate
  certificate_chain = aws_acmpca_certificate.updated.certificate_chain
}

resource "aws_acmpca_certificate" "updated" {
  certificate_authority_arn   = aws_acmpca_certificate_authority.test.arn
  certificate_signing_request = aws_acmpca_certificate_authority.test.certificate_signing_request
  signing_algorithm           = "SHA512WITHRSA"

  template_arn = "arn:${data.aws_partition.current.partition}:acm-pca:::template/RootCACertificate/V1"

  validity {
    type  = "YEARS"
    value = 1
  }
}

resource "aws_acmpca_certificate_authority" "test" {
  permanent_deletion_time_in_days = 7
  type                            = "ROOT"

  certificate_authority_configuration {
    key_algorithm     = "RSA_4096"
    signing_algorithm = "SHA512WITHRSA"

    subject {
      common_name = %[1]q
    }
  }
}

data "aws_partition" "current" {}
`, commonName)
}

func testAccCertificateAuthorityCertificateConfig_subordinateCA(commonName string) string {
	return fmt.Sprintf(`
resource "aws_acmpca_certificate_authority_certificate" "test" {
  certificate_authority_arn = aws_acmpca_certificate_authority.test.arn

  certificate       = aws_acmpca_certificate.test.certificate
  certificate_chain = aws_acmpca_certificate.test.certificate_chain
}

resource "aws_acmpca_certificate" "test" {
  certificate_authority_arn   = aws_acmpca_certificate_authority.root.arn
  certificate_signing_request = aws_acmpca_certificate_authority.test.certificate_signing_request
  signing_algorithm           = "SHA512WITHRSA"

  template_arn = "arn:${data.aws_partition.current.partition}:acm-pca:::template/SubordinateCACertificate_PathLen0/V1"

  validity {
    type  = "YEARS"
    value = 1
  }
}

resource "aws_acmpca_certificate_authority" "test" {
  permanent_deletion_time_in_days = 7
  type                            = "SUBORDINATE"

  certificate_authority_configuration {
    key_algorithm     = "RSA_2048"
    signing_algorithm = "SHA512WITHRSA"

    subject {
      common_name = "sub.%[1]s"
    }
  }
}

resource "aws_acmpca_certificate_authority" "root" {
  permanent_deletion_time_in_days = 7
  type                            = "ROOT"

  certificate_authority_configuration {
    key_algorithm     = "RSA_4096"
    signing_algorithm = "SHA512WITHRSA"

    subject {
      common_name = %[1]q
    }
  }
}

resource "aws_acmpca_certificate_authority_certificate" "root" {
  certificate_authority_arn = aws_acmpca_certificate_authority.root.arn

  certificate       = aws_acmpca_certificate.root.certificate
  certificate_chain = aws_acmpca_certificate.root.certificate_chain
}

resource "aws_acmpca_certificate" "root" {
  certificate_authority_arn   = aws_acmpca_certificate_authority.root.arn
  certificate_signing_request = aws_acmpca_certificate_authority.root.certificate_signing_request
  signing_algorithm           = "SHA512WITHRSA"

  template_arn = "arn:${data.aws_partition.current.partition}:acm-pca:::template/RootCACertificate/V1"

  validity {
    type  = "YEARS"
    value = 2
  }
}

data "aws_partition" "current" {}
`, commonName)
}
