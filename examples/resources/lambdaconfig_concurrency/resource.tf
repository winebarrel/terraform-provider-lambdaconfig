resource "lambdaconfig_concurrency" "my_func" {
  function_name                  = "my_func"
  reserved_concurrent_executions = 5
}
