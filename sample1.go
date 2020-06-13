package main

import (
  "fmt"
  "flag"
  "reflect"
  "strings"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/sqs"
  "github.com/aws/aws-sdk-go/service/sns"
  "github.com/aws/aws-sdk-go/service/lambda"
)


func resolveQueueUrl(queueName string, svc *sqs.SQS) (error, string) {
	params := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}
	resp, err := svc.GetQueueUrl(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return err, ""
	}

	return nil, *resp.QueueUrl
}



func main() {

  fmt.Println("Script starting ....")

  input := flag.String("arg1","","Argument1")
  input1 := flag.String("arg2","","Argument2")
  input2 := flag.String("arg3","","Argument3")

  flag.Parse()

  svc := sqs.New(session.New(), aws.NewConfig().WithRegion("us-east-2"))

	err, sourceUrl := resolveQueueUrl(*input, svc)

	if err != nil {
		return
	}

  fmt.Println(sourceUrl)

  params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(sourceUrl), // Required
		VisibilityTimeout:   aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
		MaxNumberOfMessages: aws.Int64(10),
  }
  
  resp, err := svc.ReceiveMessage(params)

  if len(resp.Messages) == 0 {
			fmt.Println("Batch doesn't have any messages, transfer complete")
			return
		}

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return
    }
    
    
fmt.Println(resp.Messages[0])

x := resp.Messages[0]

fmt.Println(reflect.TypeOf(x), *x.Body)

//sess := session.New(), aws.NewConfig().WithRegion("us-east-2")


svc1 := sns.New(session.New(), aws.NewConfig().WithRegion("us-east-2"))


result1, err1 := svc1.ListTopics(nil)
    if err1 != nil {
        fmt.Println(err.Error())
       
	}
	
	var topic_Arn string 
	var lambda_Arn string 


    for _, t := range result1.Topics {
		fmt.Println(*t.TopicArn)
		if strings.Contains(*t.TopicArn,*input1){
			topic_Arn = *t.TopicArn
			break
		}
		
	}

	svc2 := lambda.New(session.New(), aws.NewConfig().WithRegion("us-east-2"))

	result3, err3 := svc2.ListFunctions(nil)
if err3 != nil {
    fmt.Println("Cannot list functions")
}

for _, f := range result3.Functions {
    fmt.Println("Name:        " + aws.StringValue(f.FunctionName))
	fmt.Println("Description: " + aws.StringValue(f.Description))
	params2 := &lambda.GetFunctionInput{
		FunctionName: aws.String(aws.StringValue(f.FunctionName)),
	}
	 res,err4   := svc2.GetFunction(params2)
	 if err4 != nil {
	fmt.Println("Cannot list functions")
	}

	//fmt.Println(reflect.TypeOf(string(*res.Configuration.FunctionArn)))
	
	if strings.Contains(string(*res.Configuration.FunctionArn),*input2){
			lambda_Arn = string(*res.Configuration.FunctionArn)
			break
		}

}
fmt.Println(*input1, topic_Arn, lambda_Arn)

	result2, err2 := svc1.Subscribe(&sns.SubscribeInput{
        Endpoint:              aws.String(lambda_Arn),
        Protocol:              aws.String("lambda"),
        ReturnSubscriptionArn: aws.Bool(true), // Return the ARN, even if user has yet to confirm
        TopicArn:              aws.String(topic_Arn),
    })
    if err2 != nil {
        fmt.Println(err.Error())
    }

    fmt.Println(*result2.SubscriptionArn)

	result, err := svc1.Publish(&sns.PublishInput{
        Message:  aws.String("Work In Progress , almost there"),
        TopicArn: aws.String(topic_Arn),
    })
    if err != nil {
        fmt.Println(err.Error())
    }

    fmt.Println(*result.MessageId)
	
	fmt.Println(*input1, topic_Arn)


}
