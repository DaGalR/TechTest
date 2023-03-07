build:
	cd lambdas/create_order/ && env GOOS=linux GOARCH=amd64 go build -o ../../bin/create_order
	cd lambdas/create_payment/ && env GOOS=linux GOARCH=amd64 go build -o ../../bin/create_payment
	cd lambdas/sqs_handler && env GOOS=linux GOARCH=amd64 go build -o ../../bin/sqs_handler
	cd lambdas/update_order/ && env GOOS=linux GOARCH=amd64 go build -o ../../bin/update_order
setenv:
	export AWS_ACCESS_KEY_ID=AKIAT724GV2MVPFYKLPZ && export AWS_SECRET_ACCESS_KEY=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
init:
	cd infra/ && terraform init

plan:
	cd infra/ && terraform plan

apply:
	cd infra/ && terraform apply --auto-approve

destroy:
	cd infra/ && terraform destroy --auto-approve