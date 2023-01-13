terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
  region     = "us-east-1"
}

resource "aws_elastic_beanstalk_application" "app" {
  name        = "api-rest"
  description = "Zincsearch api-rest application built in Go"
}

# resource "aws_elastic_beanstalk_environment" "prod_environment" {
#   name                = "production"
#   application         = aws_elastic_beanstalk_application.app.name
#   solution_stack_name = "64bit Amazon Linux 2 v3.6.0 running Go 1"
# }

resource "aws_s3_bucket" "bucket" {
  bucket = "enajera-bucket-projects"
}