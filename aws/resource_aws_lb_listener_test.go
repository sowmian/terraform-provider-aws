package aws

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSLBListener_basic(t *testing.T) {
	var conf elbv2.Listener
	lbName := fmt.Sprintf("testlistener-basic-%s", acctest.RandStringFromCharSet(13, acctest.CharSetAlphaNum))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_lb_listener.front_end",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_basic(lbName, targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists("aws_lb_listener.front_end", &conf),
					resource.TestCheckResourceAttrSet("aws_lb_listener.front_end", "load_balancer_arn"),
					resource.TestCheckResourceAttrSet("aws_lb_listener.front_end", "arn"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "port", "80"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.#", "1"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.type", "forward"),
					resource.TestCheckResourceAttrSet("aws_lb_listener.front_end", "default_action.0.target_group_arn"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.fixed_response.#", "0"),
				),
			},
		},
	})
}

func TestAccAWSLBListenerBackwardsCompatibility(t *testing.T) {
	var conf elbv2.Listener
	lbName := fmt.Sprintf("testlistener-basic-%s", acctest.RandStringFromCharSet(13, acctest.CharSetAlphaNum))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_alb_listener.front_end",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfigBackwardsCompatibility(lbName, targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists("aws_alb_listener.front_end", &conf),
					resource.TestCheckResourceAttrSet("aws_alb_listener.front_end", "load_balancer_arn"),
					resource.TestCheckResourceAttrSet("aws_alb_listener.front_end", "arn"),
					resource.TestCheckResourceAttr("aws_alb_listener.front_end", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("aws_alb_listener.front_end", "port", "80"),
					resource.TestCheckResourceAttr("aws_alb_listener.front_end", "default_action.#", "1"),
					resource.TestCheckResourceAttr("aws_alb_listener.front_end", "default_action.0.type", "forward"),
					resource.TestCheckResourceAttrSet("aws_alb_listener.front_end", "default_action.0.target_group_arn"),
					resource.TestCheckResourceAttr("aws_alb_listener.front_end", "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr("aws_alb_listener.front_end", "default_action.0.fixed_response.#", "0"),
				),
			},
		},
	})
}

func TestAccAWSLBListener_https(t *testing.T) {
	var conf elbv2.Listener
	lbName := fmt.Sprintf("testlistener-https-%s", acctest.RandStringFromCharSet(13, acctest.CharSetAlphaNum))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_lb_listener.front_end",
		Providers:     testAccProvidersWithTLS,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_https(lbName, targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists("aws_lb_listener.front_end", &conf),
					resource.TestCheckResourceAttrSet("aws_lb_listener.front_end", "load_balancer_arn"),
					resource.TestCheckResourceAttrSet("aws_lb_listener.front_end", "arn"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "port", "443"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.#", "1"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.type", "forward"),
					resource.TestCheckResourceAttrSet("aws_lb_listener.front_end", "default_action.0.target_group_arn"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.fixed_response.#", "0"),
					resource.TestCheckResourceAttrSet("aws_lb_listener.front_end", "certificate_arn"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "ssl_policy", "ELBSecurityPolicy-2015-05"),
				),
			},
		},
	})
}

func TestAccAWSLBListener_redirect(t *testing.T) {
	var conf elbv2.Listener
	lbName := fmt.Sprintf("testlistener-redirect-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_lb_listener.front_end",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_redirect(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists("aws_lb_listener.front_end", &conf),
					resource.TestCheckResourceAttrSet("aws_lb_listener.front_end", "load_balancer_arn"),
					resource.TestCheckResourceAttrSet("aws_lb_listener.front_end", "arn"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "port", "80"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.#", "1"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.type", "redirect"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.target_group_arn", ""),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.redirect.#", "1"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.redirect.0.host", "#{host}"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.redirect.0.path", "/#{path}"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.redirect.0.port", "443"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.redirect.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.redirect.0.query", "#{query}"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.redirect.0.status_code", "HTTP_301"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.fixed_response.#", "0"),
				),
			},
		},
	})
}

func TestAccAWSLBListener_fixedResponse(t *testing.T) {
	var conf elbv2.Listener
	lbName := fmt.Sprintf("testlistener-fixedresponse-%s", acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_lb_listener.front_end",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_fixedResponse(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists("aws_lb_listener.front_end", &conf),
					resource.TestCheckResourceAttrSet("aws_lb_listener.front_end", "load_balancer_arn"),
					resource.TestCheckResourceAttrSet("aws_lb_listener.front_end", "arn"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "port", "80"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.#", "1"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.type", "fixed-response"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.target_group_arn", ""),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.fixed_response.#", "1"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.fixed_response.0.content_type", "text/plain"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.fixed_response.0.message_body", "Fixed response content"),
					resource.TestCheckResourceAttr("aws_lb_listener.front_end", "default_action.0.fixed_response.0.status_code", "200"),
				),
			},
		},
	})
}

func testAccCheckAWSLBListenerExists(n string, res *elbv2.Listener) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Listener ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).elbv2conn

		describe, err := conn.DescribeListeners(&elbv2.DescribeListenersInput{
			ListenerArns: []*string{aws.String(rs.Primary.ID)},
		})

		if err != nil {
			return err
		}

		if len(describe.Listeners) != 1 ||
			*describe.Listeners[0].ListenerArn != rs.Primary.ID {
			return errors.New("Listener not found")
		}

		*res = *describe.Listeners[0]
		return nil
	}
}

func testAccCheckAWSLBListenerDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).elbv2conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_lb_listener" && rs.Type != "aws_alb_listener" {
			continue
		}

		describe, err := conn.DescribeListeners(&elbv2.DescribeListenersInput{
			ListenerArns: []*string{aws.String(rs.Primary.ID)},
		})

		if err == nil {
			if len(describe.Listeners) != 0 &&
				*describe.Listeners[0].ListenerArn == rs.Primary.ID {
				return fmt.Errorf("Listener %q still exists", rs.Primary.ID)
			}
		}

		// Verify the error
		if isAWSErr(err, elbv2.ErrCodeListenerNotFoundException, "") {
			return nil
		} else {
			return fmt.Errorf("Unexpected error checking LB Listener destroyed: %s", err)
		}
	}

	return nil
}

func testAccAWSLBListenerConfig_basic(lbName, targetGroupName string) string {
	return fmt.Sprintf(`resource "aws_lb_listener" "front_end" {
   load_balancer_arn = "${aws_lb.alb_test.id}"
   protocol = "HTTP"
   port = "80"

   default_action {
     target_group_arn = "${aws_lb_target_group.test.id}"
     type = "forward"
   }
}

resource "aws_lb" "alb_test" {
  name            = "%s"
  internal        = true
  security_groups = ["${aws_security_group.alb_test.id}"]
  subnets         = ["${aws_subnet.alb_test.*.id}"]

  idle_timeout = 30
  enable_deletion_protection = false

  tags {
    Name = "TestAccAWSALB_basic"
  }
}

resource "aws_lb_target_group" "test" {
  name = "%s"
  port = 8080
  protocol = "HTTP"
  vpc_id = "${aws_vpc.alb_test.id}"

  health_check {
    path = "/health"
    interval = 60
    port = 8081
    protocol = "HTTP"
    timeout = 3
    healthy_threshold = 3
    unhealthy_threshold = 3
    matcher = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = "list"
}

data "aws_availability_zones" "available" {}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-lb-listener-basic"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = "${aws_vpc.alb_test.id}"
  cidr_block              = "${element(var.subnets, count.index)}"
  map_public_ip_on_launch = true
  availability_zone       = "${element(data.aws_availability_zones.available.names, count.index)}"

  tags {
    Name = "tf-acc-lb-listener-basic-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = "${aws_vpc.alb_test.id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "TestAccAWSALB_basic"
  }
}`, lbName, targetGroupName)
}

func testAccAWSLBListenerConfigBackwardsCompatibility(lbName, targetGroupName string) string {
	return fmt.Sprintf(`resource "aws_alb_listener" "front_end" {
   load_balancer_arn = "${aws_alb.alb_test.id}"
   protocol = "HTTP"
   port = "80"

   default_action {
     target_group_arn = "${aws_alb_target_group.test.id}"
     type = "forward"
   }
}

resource "aws_alb" "alb_test" {
  name            = "%s"
  internal        = true
  security_groups = ["${aws_security_group.alb_test.id}"]
  subnets         = ["${aws_subnet.alb_test.*.id}"]

  idle_timeout = 30
  enable_deletion_protection = false

  tags {
    Name = "TestAccAWSALB_basic"
  }
}

resource "aws_alb_target_group" "test" {
  name = "%s"
  port = 8080
  protocol = "HTTP"
  vpc_id = "${aws_vpc.alb_test.id}"

  health_check {
    path = "/health"
    interval = 60
    port = 8081
    protocol = "HTTP"
    timeout = 3
    healthy_threshold = 3
    unhealthy_threshold = 3
    matcher = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = "list"
}

data "aws_availability_zones" "available" {}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-lb-listener-bc"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = "${aws_vpc.alb_test.id}"
  cidr_block              = "${element(var.subnets, count.index)}"
  map_public_ip_on_launch = true
  availability_zone       = "${element(data.aws_availability_zones.available.names, count.index)}"

  tags {
    Name = "tf-acc-lb-listener-bc-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = "${aws_vpc.alb_test.id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "TestAccAWSALB_basic"
  }
}`, lbName, targetGroupName)
}

func testAccAWSLBListenerConfig_https(lbName, targetGroupName string) string {
	return fmt.Sprintf(`resource "aws_lb_listener" "front_end" {
   load_balancer_arn = "${aws_lb.alb_test.id}"
   protocol = "HTTPS"
   port = "443"
   ssl_policy = "ELBSecurityPolicy-2015-05"
   certificate_arn = "${aws_iam_server_certificate.test_cert.arn}"

   default_action {
     target_group_arn = "${aws_lb_target_group.test.id}"
     type = "forward"
   }
}

resource "aws_lb" "alb_test" {
  name            = "%s"
  internal        = false
  security_groups = ["${aws_security_group.alb_test.id}"]
  subnets         = ["${aws_subnet.alb_test.*.id}"]

  idle_timeout = 30
  enable_deletion_protection = false

  tags {
    Name = "TestAccAWSALB_basic"
  }

  depends_on = ["aws_internet_gateway.gw"]
}

resource "aws_lb_target_group" "test" {
  name = "%s"
  port = 8080
  protocol = "HTTP"
  vpc_id = "${aws_vpc.alb_test.id}"

  health_check {
    path = "/health"
    interval = 60
    port = 8081
    protocol = "HTTP"
    timeout = 3
    healthy_threshold = 3
    unhealthy_threshold = 3
    matcher = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = "list"
}

data "aws_availability_zones" "available" {}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-lb-listener-https"
  }
}

resource "aws_internet_gateway" "gw" {
    vpc_id = "${aws_vpc.alb_test.id}"

    tags {
        Name = "TestAccAWSALB_basic"
    }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = "${aws_vpc.alb_test.id}"
  cidr_block              = "${element(var.subnets, count.index)}"
  map_public_ip_on_launch = true
  availability_zone       = "${element(data.aws_availability_zones.available.names, count.index)}"

  tags {
    Name = "tf-acc-lb-listener-https-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = "${aws_vpc.alb_test.id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "TestAccAWSALB_basic"
  }
}

resource "aws_iam_server_certificate" "test_cert" {
  name = "terraform-test-cert-%d"
  certificate_body = "${tls_self_signed_cert.example.cert_pem}"
  private_key      = "${tls_private_key.example.private_key_pem}"
}

resource "tls_private_key" "example" {
  algorithm = "RSA"
}

resource "tls_self_signed_cert" "example" {
  key_algorithm   = "RSA"
  private_key_pem = "${tls_private_key.example.private_key_pem}"

  subject {
    common_name  = "example.com"
    organization = "ACME Examples, Inc"
  }

  validity_period_hours = 12

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}
`, lbName, targetGroupName, acctest.RandInt())
}

func testAccAWSLBListenerConfig_redirect(lbName string) string {
	return fmt.Sprintf(`resource "aws_lb_listener" "front_end" {
  load_balancer_arn = "${aws_lb.alb_test.id}"
  protocol = "HTTP"
  port = "80"

  default_action {
    type = "redirect"
    redirect {
      port = "443"
      protocol = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb" "alb_test" {
  name            = "%s"
  internal        = true
  security_groups = ["${aws_security_group.alb_test.id}"]
  subnets         = ["${aws_subnet.alb_test.*.id}"]

  idle_timeout = 30
  enable_deletion_protection = false

  tags {
    Name = "TestAccAWSALB_redirect"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = "list"
}

data "aws_availability_zones" "available" {}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-lb-listener-redirect"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = "${aws_vpc.alb_test.id}"
  cidr_block              = "${element(var.subnets, count.index)}"
  map_public_ip_on_launch = true
  availability_zone       = "${element(data.aws_availability_zones.available.names, count.index)}"

  tags {
    Name = "tf-acc-lb-listener-redirect-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = "${aws_vpc.alb_test.id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "TestAccAWSALB_redirect"
  }
}`, lbName)
}

func testAccAWSLBListenerConfig_fixedResponse(lbName string) string {
	return fmt.Sprintf(`resource "aws_lb_listener" "front_end" {
  load_balancer_arn = "${aws_lb.alb_test.id}"
  protocol = "HTTP"
  port = "80"

  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "Fixed response content"
      status_code = "200"
    }
  }
}

resource "aws_lb" "alb_test" {
  name            = "%s"
  internal        = true
  security_groups = ["${aws_security_group.alb_test.id}"]
  subnets         = ["${aws_subnet.alb_test.*.id}"]

  idle_timeout = 30
  enable_deletion_protection = false

  tags {
    Name = "TestAccAWSALB_fixedresponse"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = "list"
}

data "aws_availability_zones" "available" {}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-lb-listener-fixedresponse"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = "${aws_vpc.alb_test.id}"
  cidr_block              = "${element(var.subnets, count.index)}"
  map_public_ip_on_launch = true
  availability_zone       = "${element(data.aws_availability_zones.available.names, count.index)}"

  tags {
    Name = "tf-acc-lb-listener-fixedresponse-${count.index}"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = "${aws_vpc.alb_test.id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "TestAccAWSALB_fixedresponse"
  }
}`, lbName)
}
