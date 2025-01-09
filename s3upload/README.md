### Summary of Learning: AWS S3, SQS, and Golang Integration

#### **1. Setting Up Permissions**

To allow S3 to send events to SQS, you need to:

1. **Create an SQS Queue**:

   ```bash
   aws sqs create-queue --queue-name okharch-bucket-events
   ```

   Output:

   ```json
   {
       "QueueUrl": "https://sqs.us-east-1.amazonaws.com/897729111371/okharch-bucket-events"
   }
   ```

2. **Retrieve the SQS Queue ARN**:

   ```bash
   aws sqs get-queue-attributes --queue-url https://sqs.us-east-1.amazonaws.com/897729111371/okharch-bucket-events --attribute-names QueueArn
   ```

   Output:

   ```json
   {
       "Attributes": {
           "QueueArn": "arn:aws:sqs:us-east-1:897729111371:okharch-bucket-events"
       }
   }
   ```

3. **Set Permissions on the SQS Queue**:

   ```bash
   aws sqs set-queue-attributes --queue-url https://sqs.us-east-1.amazonaws.com/897729111371/okharch-bucket-events --attributes '{
       "Policy": "{\"Version\":\"2012-10-17\",\"Id\":\"S3ToSQSPolicy\",\"Statement\":[{\"Sid\":\"AllowS3ToSendMessage\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"s3.amazonaws.com\"},\"Action\":\"sqs:SendMessage\",\"Resource\":\"arn:aws:sqs:us-east-1:897729111371:okharch-bucket-events\",\"Condition\":{\"ArnEquals\":{\"aws:SourceArn\":\"arn:aws:s3:::okharch-bucket\"}}}]}"
   }'
   ```

4. **Configure S3 Bucket Notification**:

   ```bash
   aws s3api put-bucket-notification-configuration --bucket okharch-bucket --notification-configuration '{
       "QueueConfigurations": [
           {
               "QueueArn": "arn:aws:sqs:us-east-1:897729111371:okharch-bucket-events",
               "Events": ["s3:ObjectCreated:*"]
           }
       ]
   }'
   ```

---

#### **2. Creating and Managing Buckets**

1. **Create a Bucket**:

   ```bash
   aws s3api create-bucket --bucket okharch-bucket --region us-east-1
   ```

2. **List All Buckets**:

   ```bash
   aws s3api list-buckets
   ```

3. **List All Objects in a Bucket**:

   ```bash
   aws s3 ls s3://okharch-bucket --recursive
   ```

4. **Delete All Objects in a Bucket**:

   ```bash
   aws s3 rm s3://okharch-bucket --recursive
   ```

5. **Delete a Bucket**:

   ```bash
   aws s3api delete-bucket --bucket okharch-bucket --region us-east-1
   ```

---

#### **3. Building a Golang Application to Watch Bucket Events**

**Features**:

- The application listens for new S3 object creation events through an SQS queue.
- Logs received events to both the console and a file.
- Deletes processed messages from the SQS queue.

**Code**:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Create SQS client
	sqsClient := sqs.NewFromConfig(cfg)

	// Replace with your specific SQS Queue URL
	queueURL := "https://sqs.us-east-1.amazonaws.com/897729111371/okharch-bucket-events"

	// Open a log file
	logFile, err := os.OpenFile("s3_events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Configure logger to write to both console and file
	logger := log.New(logFile, "", log.LstdFlags)
	consoleLogger := log.New(os.Stdout, "", log.LstdFlags)

	fmt.Println("Listening for messages from SQS queue...")
	for {
		// Receive messages from the queue
		output, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 10,   // Fetch up to 10 messages
			WaitTimeSeconds:     20,   // Enable long polling
			VisibilityTimeout:   30,   // Messages will be invisible for 30 seconds if not deleted
		})
		if err != nil {
			logger.Printf("Failed to receive messages: %v", err)
			consoleLogger.Printf("Failed to receive messages: %v", err)
			time.Sleep(5 * time.Second) // Avoid excessive retries on errors
			continue
		}

		// Process each message
		for _, message := range output.Messages {
			// Log the message body
			logMessage := fmt.Sprintf("Received message: %s", aws.ToString(message.Body))
			logger.Println(logMessage)
			consoleLogger.Println(logMessage)

			// Delete the message after processing
			_, err := sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				logger.Printf("Failed to delete message: %v", err)
				consoleLogger.Printf("Failed to delete message: %v", err)
			} else {
				deletedMessageLog := fmt.Sprintf("Message deleted: %s", aws.ToString(message.MessageId))
				logger.Println(deletedMessageLog)
				consoleLogger.Println(deletedMessageLog)
			}
		}
	}
}
```

**Run Instructions**:

1. Initialize Go module and install dependencies:

   ```bash
   go mod init s3-sqs-listener
   go get github.com/aws/aws-sdk-go-v2
   ```

2. Run the application:

   ```bash
   go run main.go
   ```

3. Upload a file to the bucket to trigger events:

   ```bash
   aws s3 cp testfile.txt s3://okharch-bucket/
   ```

4. Check logs in the terminal and `s3_events.log` file.

---

This summary covers permissions, bucket management, and building a Golang application to monitor bucket events.



