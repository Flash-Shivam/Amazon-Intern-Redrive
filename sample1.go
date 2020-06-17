package main

import (
  "fmt"
  //"encoding/json"
  _ "strings"
  "strconv"
  "time"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/sqs"
  "github.com/aws/aws-sdk-go/service/sns"
  "github.com/aws/aws-sdk-go/service/lambda"
  "gopkg.in/alecthomas/kingpin.v2"
  "github.com/fatih/color"
  "github.com/apex/log/handlers/cli"
  "github.com/apex/log"
  "github.com/tj/go-progress"
  "github.com/tj/go/term"
 // "github.com/aws/aws-sdk-go/service/iam"
)


var (
	sourceQueue      = kingpin.Flag("source", "Source queue to retrieve message from").Short('s').Required().String()
	FuncName = kingpin.Flag("funcname", "Lambda function name to be redrived").Short('f').Required().String()
	region			= kingpin.Flag("region", "AWS Region for source and destination queues").Short('r').Default("us-east-2").String()
)


func logAwsError(message string, err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		log.Error(color.New(color.FgRed).Sprintf("%s. Error: %s", message, awsErr.Message()))
	} else {
		log.Error(color.New(color.FgRed).Sprintf("%s. Error: %s", message, err.Error()))
	}
}

func convertSuccessfulMessageToBatchRequestEntry(messages []*sqs.Message) []*sqs.DeleteMessageBatchRequestEntry {
	result := make([]*sqs.DeleteMessageBatchRequestEntry, len(messages))
	for i, message := range messages {
		result[i] = &sqs.DeleteMessageBatchRequestEntry{
			ReceiptHandle: message.ReceiptHandle,
			Id:            message.MessageId,
		}
	}

	return result
}

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

    log.SetHandler(cli.Default)

    fmt.Println()
	defer fmt.Println()
    //fmt.Println("Script starting ....")

    input1 := "TestTopic"

  


    kingpin.UsageTemplate(kingpin.CompactUsageTemplate)

	kingpin.Parse()

  
  svc := sqs.New(session.New(), aws.NewConfig().WithRegion(*region))

	err, sourceUrl := resolveQueueUrl(*sourceQueue, svc)

	if err != nil {
		return
	}

  //fmt.Println(sourceUrl)
  log.Info(color.New(color.FgCyan).Sprintf("Source queue url: %s", sourceUrl))


  	queueAttributes, err := svc.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(sourceUrl),
		AttributeNames: []*string{aws.String("All")},
	})

	if err != nil {
		logAwsError("Failed to resolve queue attributes", err)
		return
	}

	numberOfMessages, _ := strconv.Atoi(*queueAttributes.Attributes["ApproximateNumberOfMessages"])

    log.Info(color.New(color.FgCyan).Sprintf("Approximate number of messages in the source queue: %d", numberOfMessages))
    
    if numberOfMessages == 0 {
		log.Info("Looks like nothing to move. Done.")
		return
	}

  params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(sourceUrl), // Required
		VisibilityTimeout:   aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
		MaxNumberOfMessages: aws.Int64(10),
  }


//sess := session.New(), aws.NewConfig().WithRegion("us-east-2")


svc1 := sns.New(session.New(), aws.NewConfig().WithRegion(*region))

    result10, err := svc1.CreateTopic(&sns.CreateTopicInput{
        Name: aws.String(input1),
    })
	
	var topic_Arn string 
	var lambda_Arn string 
    topic_Arn = string(*result10.TopicArn)


    svc2 := lambda.New(session.New(), aws.NewConfig().WithRegion(*region))
    
    params2 := &lambda.GetFunctionInput{
		FunctionName: aws.String(*FuncName),
	}
	 res,err4   := svc2.GetFunction(params2)
	 if err4 != nil {
	fmt.Println("Cannot list functions")
    }
    
    lambda_Arn = string(*res.Configuration.FunctionArn)

	_ , err2 := svc1.Subscribe(&sns.SubscribeInput{
        Endpoint:              aws.String(lambda_Arn),
        Protocol:              aws.String("lambda"),
        ReturnSubscriptionArn: aws.Bool(true), // Return the ARN, even if user has yet to confirm
		TopicArn:              aws.String(topic_Arn),
		//principal:             aws.String("sns.us-east-2.amazonaws.com"),
    })
    if err2 != nil {
        fmt.Println(err.Error())
    }

    log.Info(color.New(color.FgCyan).Sprintf("Starting to move messages..."))
	fmt.Println()

	term.HideCursor()
	defer term.ShowCursor()
    totalMessages := numberOfMessages
	b := progress.NewInt(totalMessages)
	b.Width = 40
	b.StartDelimiter = color.New(color.FgCyan).Sprint("|")
	b.EndDelimiter = color.New(color.FgCyan).Sprint("|")
	b.Filled = color.New(color.FgCyan).Sprint("█")
	b.Empty = color.New(color.FgCyan).Sprint("░")
	b.Template(`		{{.Bar}} {{.Text}}{{.Percent | printf "%3.0f"}}%`)

	render := term.Renderer()

	messagesProcessed := 0
    start := time.Now()


for 	
{
		resp, err := svc.ReceiveMessage(params)

  if len(resp.Messages) == 0 {
			break
			
		}

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			
    }
    //fmt.Println("Number of Messages : " , len(resp.Messages))
    
    
    messagesToCopy := resp.Messages
	for i:= 0; i < len(resp.Messages) ; i++ {

		x := resp.Messages[i]
		_ , err := svc1.Publish(&sns.PublishInput{
        Message:  aws.String(*x.Body),
        TopicArn: aws.String(topic_Arn),
    })
    if err != nil {
        fmt.Println(err.Error())
    }
    
    
	//fmt.Println("i = " , i)
        messagesProcessed += 1

        if messagesProcessed > totalMessages {
			b.Total = float64(messagesProcessed)
		}

        b.ValueInt(messagesProcessed)
		render(b.String())
    }
    deleteMessageBatch := &sqs.DeleteMessageBatchInput{
				Entries:  convertSuccessfulMessageToBatchRequestEntry(messagesToCopy),
				QueueUrl: aws.String(sourceUrl),
			}

    _  , err22 := svc.DeleteMessageBatch(deleteMessageBatch)
    if err22 != nil {
        fmt.Println(err.Error())
    }

}
    t := time.Now()
    elapsed := t.Sub(start)
    log.Info(color.New(color.FgCyan).Sprintf("Done. Moved %s messages in %f seconds", strconv.Itoa(totalMessages), elapsed.Seconds()))
	//log.Info(color.New(color.FgCyan).Sprintf("Done. Moved %s messages", strconv.Itoa(totalMessages)))

	


}
