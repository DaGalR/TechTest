#BASIC CONFIGURATION FOR TERRAFORM
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}
#DYNAMO CONFIG
resource "aws_dynamodb_table" "transactions_table" {
    name = "transactions_table"
    billing_mode = "PROVISIONED"
    read_capacity = "30"
    write_capacity = "30"
    hash_key = "PK"
    range_key = "SK"

    attribute{
        name = "PK"
        type = "S"
    }

    attribute {
        name = "SK"
        type = "S"
    }

}

#LAMBDAS PERMISSIONS AND DECLARATIONS
resource "aws_iam_role" "iam_for_lambda" {
 name = "iam_for_lambda"

 assume_role_policy = jsonencode({
   "Version" : "2012-10-17",
   "Statement" : [
     {
       "Effect" : "Allow",
       "Principal" : {
         "Service" : "lambda.amazonaws.com"
       },
       "Action" : "sts:AssumeRole"
     }
   ]
  })
}
          
resource "aws_iam_role_policy_attachment" "lambda_policy" {
   role = aws_iam_role.iam_for_lambda.name
   policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_sqs_role_policy" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
}
          
resource "aws_iam_role_policy" "dynamodb-lambda-policy" {
   name = "dynamodb_lambda_policy"
   role = aws_iam_role.iam_for_lambda.id
   policy = jsonencode({
      "Version" : "2012-10-17",
      "Statement" : [
        {
           "Effect" : "Allow",
           "Action" : ["dynamodb:*"],
           "Resource" : "${aws_dynamodb_table.transactions_table.arn}"
        }
      ]
   })
}

resource "aws_iam_role_policy" "sqs-policy" {
  name        = "sqs-policy"
  role = aws_iam_role.iam_for_lambda.id
  policy      = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "sqs:SendMessage"
        ]
        Effect   = "Allow"
        Resource = "${aws_sqs_queue.terraform_queue.arn}"
      }
    ]
  })
}

data "archive_file" "create-order-archive"{
  source_file = "../bin/create_order"
  output_path = "../bin/create_order.zip"
  type = "zip"
}

data "archive_file" "create-payment-archive"{
  source_file = "../bin/create_payment"
  output_path = "../bin/create_payment.zip"
  type = "zip"
}

data "archive_file" "sqs-handler-archive"{
  source_file = "../bin/sqs_handler"
  output_path = "../bin/sqs_handler.zip"
  type = "zip"
}

resource "aws_lambda_function" "create-order" {
  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.transactions_table.name,
      SQS_QUEUE = aws_sqs_queue.terraform_queue.id
    }
  }
  memory_size = "128"
  timeout = 10
  runtime = "go1.x"
  architectures = ["x86_64"]
  function_name = "create_order"
  handler = "create_order"
  role = aws_iam_role.iam_for_lambda.arn
  filename = data.archive_file.create-order-archive.output_path
}

resource "aws_lambda_function" "create-payment" {
  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.transactions_table.name,
      SQS_QUEUE = aws_sqs_queue.terraform_queue.id
    }
  }
  memory_size = "128"
  timeout = 10
  runtime = "go1.x"
  architectures = ["x86_64"]
  function_name = "create_payment"
  handler = "create_payment"
  role = aws_iam_role.iam_for_lambda.arn
  filename = data.archive_file.create-payment-archive.output_path
}

resource "aws_lambda_function" "lambda_sqs_handler" {
  environment {
    variables = {
      ORDERS_URL = "${aws_api_gateway_deployment.transactions_deploy.invoke_url}${aws_api_gateway_resource.order-resource.path}"
    }
  }
  memory_size = "128"
  timeout = 10
  runtime = "go1.x"
  architectures = ["x86_64"]
  function_name = "sqs_handler"
  handler = "sqs_handler"
  role = aws_iam_role.iam_for_lambda.arn
  filename = data.archive_file.sqs-handler-archive.output_path
}

#API GATEWAY
resource "aws_api_gateway_rest_api" "api" {
  name = "transactions_api"
}
#GATEWAY RESOURCES
resource "aws_api_gateway_resource" "order-resource" {
  path_part   = "order"
  parent_id   = "${aws_api_gateway_rest_api.api.root_resource_id}"
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
}
resource "aws_api_gateway_resource" "payment-resource" {
  path_part   = "payment"
  parent_id   = "${aws_api_gateway_rest_api.api.root_resource_id}"
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
}

#MODEL AND VALIDATOR FOR CUSTOM REQUEST
resource "aws_api_gateway_model" "create-order-req" {
  rest_api_id  = aws_api_gateway_rest_api.api.id
  name         = "POSTRequestModelCreateOrder"
  description  = "A JSON schema"
  content_type = "application/json"
  schema       = file("${path.module}/create_order.json")
}

resource "aws_api_gateway_request_validator" "create-order-validator" {
  name                        = "POSTRequestModelCreateOrder"
  rest_api_id                 = aws_api_gateway_rest_api.api.id
  validate_request_body       = true
  validate_request_parameters = false
}

resource "aws_api_gateway_model" "create-payment-req" {
  rest_api_id  = aws_api_gateway_rest_api.api.id
  name         = "POSTRequestModelCreatePayment"
  description  = "A JSON schema"
  content_type = "application/json"
  schema       = file("${path.module}/create_payment.json")
}

resource "aws_api_gateway_request_validator" "create-payment-validator" {
  name                        = "POSTRequestModelCreatePayment"
  rest_api_id                 = aws_api_gateway_rest_api.api.id
  validate_request_body       = true
  validate_request_parameters = false
}

#GATEWAY METHODS
resource "aws_api_gateway_method" "create-order-method" {
  rest_api_id   = "${aws_api_gateway_rest_api.api.id}"
  resource_id   = "${aws_api_gateway_resource.order-resource.id}"
  http_method   = "POST"
  authorization = "NONE"
  request_validator_id = aws_api_gateway_request_validator.create-order-validator.id
  request_models = {
    "application/json" = aws_api_gateway_model.create-order-req.name
  }
}
resource "aws_api_gateway_method" "create-payment-method" {
  rest_api_id   = "${aws_api_gateway_rest_api.api.id}"
  resource_id   = "${aws_api_gateway_resource.payment-resource.id}"
  http_method   = "POST"
  authorization = "NONE"
  request_validator_id = aws_api_gateway_request_validator.create-payment-validator.id
  request_models = {
    "application/json" = aws_api_gateway_model.create-payment-req.name
  }
}
resource "aws_api_gateway_method_response" "create-order-response" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.order-resource.id
  http_method = aws_api_gateway_method.create-order-method.http_method
  status_code = "200"
}
resource "aws_api_gateway_method_response" "create-payment-response" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.payment-resource.id
  http_method = aws_api_gateway_method.create-payment-method.http_method
  status_code = "200"
}
#INTEGRATIONS WITH LAMBDAS
resource "aws_api_gateway_integration" "create-order-integration" {
  rest_api_id             = "${aws_api_gateway_rest_api.api.id}"
  resource_id             = "${aws_api_gateway_resource.order-resource.id}"
  http_method             = "${aws_api_gateway_method.create-order-method.http_method}"
  integration_http_method = "POST"
  type                    = "AWS"
  uri                     = "${aws_lambda_function.create-order.invoke_arn}"
}
resource "aws_api_gateway_integration_response" "create-order-integration-res" {
  depends_on  = [aws_api_gateway_integration.create-order-integration]
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.order-resource.id
  http_method = aws_api_gateway_method.create-order-method.http_method
  status_code = aws_api_gateway_method_response.create-order-response.status_code
}
resource "aws_api_gateway_integration" "create-payment-integration" {
  rest_api_id             = "${aws_api_gateway_rest_api.api.id}"
  resource_id             = "${aws_api_gateway_resource.payment-resource.id}"
  http_method             = "${aws_api_gateway_method.create-payment-method.http_method}"
  integration_http_method = "POST"
  type                    = "AWS"
  uri                     = "${aws_lambda_function.create-payment.invoke_arn}"
}
resource "aws_api_gateway_integration_response" "create-payment-integration-res" {
  depends_on  = [aws_api_gateway_integration.create-payment-integration]
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.payment-resource.id
  http_method = aws_api_gateway_method.create-payment-method.http_method
  status_code = aws_api_gateway_method_response.create-payment-response.status_code
}
#LAMBDA PERMISSIONS FOR GATEWAY
resource "aws_lambda_permission" "apigw-create-order-lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.create-order.function_name}"
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.api.execution_arn}/*/*/*"
}

resource "aws_lambda_permission" "apigw-create-payment-lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.create-payment.function_name}"
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.api.execution_arn}/*/*/*"
}

#API GATEWAY DEPLOY
resource "aws_api_gateway_deployment" "transactions_deploy" {
  depends_on = [aws_api_gateway_integration.create-order-integration, aws_api_gateway_integration.create-payment-integration]

  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  stage_name  = "v1"

}

#SQS
resource "aws_sqs_queue" "terraform_queue" {
  name                      = "api-events-queue"
  delay_seconds             = 90
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
  #redrive_policy            = "{\"deadLetterTargetArn\":\"${aws_sqs_queue.terraform_queue_deadletter.arn}\",\"maxReceiveCount\":3}"
}

resource "aws_lambda_event_source_mapping" "event_source_mapping" {
  event_source_arn = aws_sqs_queue.terraform_queue.arn
  function_name    = aws_lambda_function.create-order.arn
  enabled          = true
}

resource "aws_lambda_event_source_mapping" "sqs_event_handler_mapping" {
  event_source_arn = aws_sqs_queue.terraform_queue.arn
  function_name    = aws_lambda_function.lambda_sqs_handler.arn
  enabled          = true
}

output "sqs-url"{
  value = "${aws_sqs_queue.terraform_queue.id}"
}

output "url-order" {
  value = "${aws_api_gateway_deployment.transactions_deploy.invoke_url}${aws_api_gateway_resource.order-resource.path}"
}
output "url-payment" {
  value = "${aws_api_gateway_deployment.transactions_deploy.invoke_url}${aws_api_gateway_resource.payment-resource.path}"
}

