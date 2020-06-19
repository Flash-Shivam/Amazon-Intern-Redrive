import json
import boto3
#import jsonify

client = boto3.client('lambda')
sqs = boto3.client('sqs')


def lambda_handler(event, context):
    name = event['Name']
    account_id = event['AccountID']
    
    response1 = sqs.get_queue_url(
    QueueName= name ,
    QueueOwnerAWSAccountId=account_id
    )
    queue_url = response1['QueueUrl']
    #print(response1)
    response2 = sqs.receive_message(
    QueueUrl=queue_url,
    AttributeNames=[
        'SentTimestamp'
    ],
    MaxNumberOfMessages=10,
    MessageAttributeNames=[
        'All'
    ],
    VisibilityTimeout=20,
    WaitTimeSeconds=0
    )
    
    #if "Messages" in response2:
    x = response2['Messages']
    
    
    while True :
        for i in range(0, len(x)):
            inputforinvoker = {
        'MessageId':  x[i]['MessageId'],
        'Type': 'Notification',
        'TopicARN' : 'TrialBasis123abd'
        
    }
            
            response3 = client.invoke(
        FunctionName = 'LambdaToInvoke',
        InvocationType = 'RequestResponse',
        Payload=json.dumps(inputforinvoker)
        )

            response = sqs.delete_message(
    QueueUrl=queue_url,
    ReceiptHandle=(x[i]['ReceiptHandle'])
    )
        #j.append(response)
        
        response2 = sqs.receive_message(
    QueueUrl=queue_url,
    AttributeNames=[
        'SentTimestamp'
    ],
    MaxNumberOfMessages=10,
    MessageAttributeNames=[
        'All'
    ],
    VisibilityTimeout=20,
    WaitTimeSeconds=0
    )
        if "Messages" in response2:
            x = response2['Messages']
        else:
            break

    

    #print(message, len(response2['Messages']))
    
    
    #return responseJson, x
    
    return "Done"
    #return inputforinvoker