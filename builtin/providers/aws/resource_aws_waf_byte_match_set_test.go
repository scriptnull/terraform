package aws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAWSWaf_all(t *testing.T) {
	cases := map[string]func(*testing.T){
		"ByteMatchSet_basic":                      TestAccAWSWafByteMatchSet_basic,
		"ByteMatchSet_changeNameForceNew":         TestAccAWSWafByteMatchSet_changeNameForceNew,
		"ByteMatchSet_disappears":                 TestAccAWSWafByteMatchSet_disappears,
		"IPSet_basic":                             TestAccAWSWafIPSet_basic,
		"IPSet_disappears":                        TestAccAWSWafIPSet_disappears,
		"IPSet_changeNameForceNew":                TestAccAWSWafIPSet_changeNameForceNew,
		"Rule_basic":                              TestAccAWSWafRule_basic,
		"Rule_changeNameForceNew":                 TestAccAWSWafRule_changeNameForceNew,
		"Rule_disappears":                         TestAccAWSWafRule_disappears,
		"SizeConstraintSet_basic":                 TestAccAWSWafSizeConstraintSet_basic,
		"SizeConstraintSet_changeNameForceNew":    TestAccAWSWafSizeConstraintSet_changeNameForceNew,
		"SizeConstraintSet_disappears":            TestAccAWSWafSizeConstraintSet_disappears,
		"SqlInjectionMatchSet_basic":              TestAccAWSWafSqlInjectionMatchSet_basic,
		"SqlInjectionMatchSet_changeNameForceNew": TestAccAWSWafSqlInjectionMatchSet_changeNameForceNew,
		"SqlInjectionMatchSet_disappears":         TestAccAWSWafSqlInjectionMatchSet_disappears,
		"WebAcl_basic":                            TestAccAWSWafWebAcl_basic,
		"WebAcl_changeNameForceNew":               TestAccAWSWafWebAcl_changeNameForceNew,
		"WebAcl_changeDefaultAction":              TestAccAWSWafWebAcl_changeDefaultAction,
		"WebAcl_disappears":                       TestAccAWSWafWebAcl_disappears,
		"XssMatchSet_basic":                       TestAccAWSWafXssMatchSet_basic,
		"XssMatchSet_changeNameForceNew":          TestAccAWSWafXssMatchSet_changeNameForceNew,
		"XssMatchSet_disappears":                  TestAccAWSWafXssMatchSet_disappears,
	}
	for name, tf := range cases {
		tf := tf
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tf(t)
		})
	}
}

func TestAccAWSWafByteMatchSet_basic(t *testing.T) {
	var v waf.ByteMatchSet
	byteMatchSet := fmt.Sprintf("byteMatchSet-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSWafByteMatchSetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSWafByteMatchSetConfig(byteMatchSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSWafByteMatchSetExists("aws_waf_byte_match_set.byte_set", &v),
					resource.TestCheckResourceAttr(
						"aws_waf_byte_match_set.byte_set", "name", byteMatchSet),
					resource.TestCheckResourceAttr(
						"aws_waf_byte_match_set.byte_set", "byte_match_tuples.#", "2"),
				),
			},
		},
	})
}

func TestAccAWSWafByteMatchSet_changeNameForceNew(t *testing.T) {
	var before, after waf.ByteMatchSet
	byteMatchSet := fmt.Sprintf("byteMatchSet-%s", acctest.RandString(5))
	byteMatchSetNewName := fmt.Sprintf("byteMatchSet-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSWafByteMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSWafByteMatchSetConfig(byteMatchSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSWafByteMatchSetExists("aws_waf_byte_match_set.byte_set", &before),
					resource.TestCheckResourceAttr(
						"aws_waf_byte_match_set.byte_set", "name", byteMatchSet),
					resource.TestCheckResourceAttr(
						"aws_waf_byte_match_set.byte_set", "byte_match_tuples.#", "2"),
				),
			},
			{
				Config: testAccAWSWafByteMatchSetConfigChangeName(byteMatchSetNewName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSWafByteMatchSetExists("aws_waf_byte_match_set.byte_set", &after),
					resource.TestCheckResourceAttr(
						"aws_waf_byte_match_set.byte_set", "name", byteMatchSetNewName),
					resource.TestCheckResourceAttr(
						"aws_waf_byte_match_set.byte_set", "byte_match_tuples.#", "2"),
				),
			},
		},
	})
}

func TestAccAWSWafByteMatchSet_disappears(t *testing.T) {
	var v waf.ByteMatchSet
	byteMatchSet := fmt.Sprintf("byteMatchSet-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSWafByteMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSWafByteMatchSetConfig(byteMatchSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSWafByteMatchSetExists("aws_waf_byte_match_set.byte_set", &v),
					testAccCheckAWSWafByteMatchSetDisappears(&v),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAWSWafByteMatchSetDisappears(v *waf.ByteMatchSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).wafconn

		wt := newWAFToken(conn, "global")
		token, err := wt.Acquire()
		if err != nil {
			return fmt.Errorf("Error getting change token: %s", err)
		}

		req := &waf.UpdateByteMatchSetInput{
			ChangeToken:    token,
			ByteMatchSetId: v.ByteMatchSetId,
		}

		for _, ByteMatchTuple := range v.ByteMatchTuples {
			ByteMatchUpdate := &waf.ByteMatchSetUpdate{
				Action: aws.String("DELETE"),
				ByteMatchTuple: &waf.ByteMatchTuple{
					FieldToMatch:         ByteMatchTuple.FieldToMatch,
					PositionalConstraint: ByteMatchTuple.PositionalConstraint,
					TargetString:         ByteMatchTuple.TargetString,
					TextTransformation:   ByteMatchTuple.TextTransformation,
				},
			}
			req.Updates = append(req.Updates, ByteMatchUpdate)
		}

		_, err = conn.UpdateByteMatchSet(req)
		wtErr := wt.Release()
		if wtErr != nil {
			return wtErr
		}
		if err != nil {
			return errwrap.Wrapf("[ERROR] Error updating ByteMatchSet: {{err}}", err)
		}

		token, err = wt.Acquire()
		if err != nil {
			return errwrap.Wrapf("[ERROR] Error getting change token: {{err}}", err)
		}

		opts := &waf.DeleteByteMatchSetInput{
			ChangeToken:    token,
			ByteMatchSetId: v.ByteMatchSetId,
		}
		_, err = conn.DeleteByteMatchSet(opts)
		wtErr = wt.Release()
		if wtErr != nil {
			return wtErr
		}
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckAWSWafByteMatchSetExists(n string, v *waf.ByteMatchSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No WAF ByteMatchSet ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).wafconn
		resp, err := conn.GetByteMatchSet(&waf.GetByteMatchSetInput{
			ByteMatchSetId: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		if *resp.ByteMatchSet.ByteMatchSetId == rs.Primary.ID {
			*v = *resp.ByteMatchSet
			return nil
		}

		return fmt.Errorf("WAF ByteMatchSet (%s) not found", rs.Primary.ID)
	}
}

func testAccCheckAWSWafByteMatchSetDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_waf_byte_match_set" {
			continue
		}

		conn := testAccProvider.Meta().(*AWSClient).wafconn
		resp, err := conn.GetByteMatchSet(
			&waf.GetByteMatchSetInput{
				ByteMatchSetId: aws.String(rs.Primary.ID),
			})

		if err == nil {
			if *resp.ByteMatchSet.ByteMatchSetId == rs.Primary.ID {
				return fmt.Errorf("WAF ByteMatchSet %s still exists", rs.Primary.ID)
			}
		}

		// Return nil if the ByteMatchSet is already destroyed
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "WAFNonexistentItemException" {
				return nil
			}
		}

		return err
	}

	return nil
}

func testAccAWSWafByteMatchSetConfig(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_byte_match_set" "byte_set" {
  name = "%s"
  byte_match_tuples {
    text_transformation = "NONE"
    target_string = "badrefer1"
    positional_constraint = "CONTAINS"
    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }

  byte_match_tuples {
    text_transformation = "NONE"
    target_string = "badrefer2"
    positional_constraint = "CONTAINS"
    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }
}`, name)
}

func testAccAWSWafByteMatchSetConfigChangeName(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_byte_match_set" "byte_set" {
  name = "%s"
  byte_match_tuples {
    text_transformation = "NONE"
    target_string = "badrefer1"
    positional_constraint = "CONTAINS"
    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }

  byte_match_tuples {
    text_transformation = "NONE"
    target_string = "badrefer2"
    positional_constraint = "CONTAINS"
    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }
}`, name)
}
