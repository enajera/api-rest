variable "AWS_ACCESS_KEY_ID" {}
variable "AWS_SECRET_ACCESS_KEY" {}

terraform {

  cloud {
    organization = "vinn-org"
    workspaces {
      name = "aws-workspace"
    }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region     = "us-east-1"
 }

resource "aws_s3_bucket" "example" {
  bucket = "enajera-bucket-projects"
}


resource "aws_elastic_beanstalk_application" "example" {
  name        = "api-rest"
  description = "Zincsearch - go api rest"
}

resource "aws_elastic_beanstalk_environment" "example" {
  name                = "api-rest-dev"
  application         = aws_elastic_beanstalk_application.example.name
  solution_stack_name = "64bit Amazon Linux 2 running Go 1"
}
