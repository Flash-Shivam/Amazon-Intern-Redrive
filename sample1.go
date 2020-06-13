package main

import (
  "fmt"
  "flag"
  "reflect"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/sqs"
  "github.com/aws/aws-sdk-go/service/sns"
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


result1, err1 := svc1.CreateTopic(&sns.CreateTopicInput{
        Name: aws.String(*input1),
	})

	if err1 != nil {
        fmt.Println(err.Error())
        
    }
fmt.Println(*result1)


}
