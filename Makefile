build:
	cd lambdas/create_order/ && env GOOS=linux GOARCH=amd64 go build -o ../../bin/create_order
	cd lambdas/create_payment/ && env GOOS=linux GOARCH=amd64 go build -o ../../bin/create_payment
	cd lambdas/sqs_handler && env GOOS=linux GOARCH=amd64 go build -o ../../bin/sqs_handler
	cd lambdas/update_order/ && env GOOS=linux GOARCH=amd64 go build -o ../../bin/update_order
init:
	cd infra/ && terraform init

plan:
	cd infra/ && terraform plan

apply:
	cd infra/ && terraform apply --auto-approve

destroy:
	cd infra/ && terraform destroy --auto-approve