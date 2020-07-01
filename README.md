<h1>Redriver</h1>

**Redriver** is a tool for moving AWS SQS messages from one queue (DLQ) to invoke(trigger) a lambda function through a SNS topic 
<h1>Preview</h1>

![](https://github.com/Flash-Shivam/Amazon-Intern-Redrive/blob/master/Screenshot%20(190).png)

<h1>Features</h1>

- Queue name resolution. For ease of use, you only need to provide a queue name and not the full arn address.

- Progress indicator.

- User friendly info and error messages.

- Reliable delivery. Redriver will only delete messages from the source queue after they were enqueued to the destination.

<h1>Prerequisites</h1>

Specifying AWS credentials in either a shared credentials file or environment variables.
Creating the Credentials File

If you don’t have a shared credentials file (~/.aws/credentials), you can use any text editor to create one in your home directory. Add the following content to your credentials file, replacing <YOUR_ACCESS_KEY_ID> and <YOUR_SECRET_ACCESS_KEY> with your credentials.

```
[default]
aws_access_key_id = <YOUR_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_SECRET_ACCESS_KEY>
```


The [default] heading defines credentials for the default profile, which the Redriver will use unless you configure it to use another profile.

Optionally you can configure default region in ~/.aws/config

[default]
region=us-east-2

<h1>Setting Up</h1>

You will need to have Golang installed.

First create a SNS Topic named 'TestTopic'. While configuring the lambda that is to be invoked , add the SNS Topic as one of the trigger .See the picture for better understanding .

![](https://github.com/Flash-Shivam/Amazon-Intern-Redrive/blob/master/Screenshot%20(222).png)

```
git clone https://github.com/Flash-Shivam/Amazon-Intern-Redrive
```


<h1>Usage</h1>

```
usage: go run sample1.go -s SOURCE -f LAMBDAFUNCTIONNAME
```

Region will be set default to us-east-2, you can also override it with --region flag
