# Tech Challenge

Withing this repo you will find he folders holfing the code for the tech challenge for backend position. The Terraform code can be found within the "infra" folder, the lambdas are in the "lambdas" folder. Don't delete "mocks"folder or any of its contents as it is used for the unit tests. The rest of the folders ("adapters", "service", "domain", and "entrypoints") hold code following hexagonal architecure pattern.

# Make

In the source folder is a Makefile with the following commands:

- build: which will compile the four lambda functions created to solve the challenge.

- init: which initializes Terraform for creating the serverless app.

- plan: which devices a plan to develop the architecture.

- apply: which builds the architecture (for this it is important to first run the "make build" command, so that Terraform can get access to the .zip files generated from each lambda). This command outputs the URLS of the four lambda functions created.

- destroy: which deletes all the elements from the architecture.

# Tests

The unit tests were written assisted by the "mockery" Go package wich helps building mocks for the interfaces in the code so it is easier to test.

You will also find a file named "coverage.html" which holds the results from the unit tests. To generate this file you should run the comand:

$ go test ./... -coverprofile=coverage.html

And to display the test coverage you should run

$ go tool cover -html=coverage.html

Both commands from the source folder.

# Instructions

To run the application, first set thee environment variables provided with these two commands:

$ export AWS_ACCESS_KEY_ID=<the_key_id_provided>
$ AWS_SECRET_ACCESS_KEY=<the_access_key_id_provided>

Then you should run the make commands described above in the same order as written:

$ make build

$ make init

$ make plan

$ make apply

You should use make destroy only when you are done reviewing the live architecture in AWS with your IAM role.

To send requests to the API Gateway, you should select the api gateway named "transactions-api" from the console, then you will see te entrypoints calld "ORDER", "PAYMENT" and "UPDATE-ORDER".

The app working flow should be: send an order, (which will appear in the Dynamo table called "transactions_table") and then send a payment (which will also appear in the same table. Also note that the order status will change from "Incomplete" to "Ready for shipping" after the payment was correctly received.). To send the body for the endpoints, you should use JSON objects like the following: (These objects are to be used as input in "Request Body" field when clicking the "TEST" button on the desired entrypoint or also you can use POSTMAN, create a POST method with the URL of each lambda which you get from the "make apply" command output and set the request body)

- For order entrypoint:

{
"order_id":"00001",
"user_id": "dani",
"item":"Papas",
"quantity":2,
"total_price": 4.5
}

- For payments entrypoint:

{
"order_id":"00001",
"status": "Complete"
}

- For update order entrypoint:

{
"order_id":"00001",
"new_status": "Ready_for_shipping"
}

The entrypoint "UPDATE-ORDER" entrypoint is called internally by de SQSHandler lambda function. This entrypoint changes the order status once the payment is correctly received. You should only use this entrypoint from the API Gateway diretly with the purpose of testing its functionality.

Also, to check about the SQS queue receiving the events, you will find it by the name "api-events-queue".

# IAM ROLE

The IAM role's console, username and password to see the architecture will be provided to you as well.

# IMPORTANT ADDITIONAL INFORMATION

Due to my working schedule and other inconveniences I could not manage to write ALL the unit tests for all the code (you will realize this by looking at the coverage.html file), still I know the ideal would be 100% coverage for all files.
