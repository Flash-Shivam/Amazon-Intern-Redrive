Redriver

Redriver is a tool for moving AWS SQS messages from one queue (DLQ) to invoke(trigger) a lambda function through a SNS topic .

Feautures

Queue name resolution. For ease of use, you only need to provide a queue name and not the full arn address.
Progress indicator.
User friendly info and error messages.
Reliable delivery. Redriver will only delete messages from the source queue after they were enqueued to the destination.

Prerequisites

Specifying AWS credentials in either a shared credentials file or environment variables.
Creating the Credentials File

If you don’t have a shared credentials file (~/.aws/credentials), you can use any text editor to create one in your home directory. Add the following content to your credentials file, replacing <YOUR_ACCESS_KEY_ID> and <YOUR_SECRET_ACCESS_KEY> with your credentials.

[default]
aws_access_key_id = <YOUR_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_SECRET_ACCESS_KEY>

The [default] heading defines credentials for the default profile, which the SQSMover will use unless you configure it to use another profile.

Optionally you can configure default region in ~/.aws/config

[default]
region=us-west-2

Setting Up

You will need to have Golang installed.

git clone https://github.com/Flash-Shivam/Amazon-Intern-Redrive

Usage

usage: go run sample1.go -s SOURCE -f LAMBDAFUNCTIONNAME

Region will default to us-east-2, you can also override it with --region flag
