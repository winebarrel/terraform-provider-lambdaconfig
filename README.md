# terraform-provider-lambdaconfig

[![CI](https://github.com/winebarrel/terraform-provider-lambdaconfig/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/terraform-provider-lambdaconfig/actions/workflows/ci.yml)

Terraform provider for AWS Lambda function configuration.

## Usage

```tf
terraform {
  required_providers {
    lambdaconfig = {
      source  = "winebarrel/lambdaconfig"
      version = ">= 0.2.0"
    }
  }
}

provider "lambdaconfig" {
  # region = "ap-northeast-1"
}

# import {
#   to = lambdaconfig_concurrency.my_func
#   id = "my_func"
# }

resource "lambdaconfig_concurrency" "my_func" {
  function_name                  = "my_func"
  reserved_concurrent_executions = 1
}
```

## Run locally for development

```sh
cp lambdaconfig.tf.sample lambdaconfig.tf
make
make tf-plan
make tf-apply
```
