package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
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
			MaxNumberOfMessages: 10, // Fetch up to 10 messages
			WaitTimeSeconds:     20, // Enable long polling
			VisibilityTimeout:   30, // Messages will be invisible for 30 seconds if not deleted
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
			// Determine if the file is textual and print first line if applicable
			err := processS3Event(aws.ToString(message.Body), logger, consoleLogger)
			if err != nil {
				logger.Printf("Error processing message: %v", err)
				consoleLogger.Printf("Error processing message: %v", err)
			}
			logger.Println(logMessage)
			consoleLogger.Println(logMessage)

			// Delete the message after processing
			_, err = sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
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

// processS3Event processes the S3 event, determines if the file is textual, and prints the first line
func processS3Event(messageBody string, logger *log.Logger, consoleLogger *log.Logger) error {
	// Parse the S3 event to extract bucket name and object key
	bucket, key, err := extractS3Details(messageBody)
	if err != nil {
		logger.Printf("Failed to extract S3 details: %v", err)
		consoleLogger.Printf("Failed to extract S3 details: %v", err)
		return err
	}

	// Simulate checking if the file is textual and printing the first line
	isTextual, firstLine, err := fetchAndCheckFile(bucket, key)
	if err != nil {
		logger.Printf("Failed to process file: %v", err)
		consoleLogger.Printf("Failed to process file: %v", err)
		return err
	}

	if isTextual {
		logger.Printf("First line of textual file: %s", firstLine)
		consoleLogger.Printf("First line of textual file: %s", firstLine)
	} else {
		logger.Printf("Non-textual file or file skipped: %s/%s", bucket, key)
		consoleLogger.Printf("Non-textual file or file skipped: %s/%s", bucket, key)
	}

	return nil
}

// Mock function to extract bucket and key from S3 event
func extractS3Details(messageBody string) (string, string, error) {
	// Parse messageBody to extract bucket name and key
	var event map[string]interface{}
	err := json.Unmarshal([]byte(messageBody), &event)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse S3 event message: %w", err)
	}

	records, ok := event["Records"].([]interface{})
	if !ok || len(records) == 0 {
		return "", "", fmt.Errorf("no records found in S3 event")
	}

	record, ok := records[0].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("invalid record structure")
	}

	s3Info, ok := record["s3"].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("missing s3 information in record")
	}

	bucketInfo, ok := s3Info["bucket"].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("missing bucket information in record")
	}
	bucketName, ok := bucketInfo["name"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid bucket name")
	}

	objectInfo, ok := s3Info["object"].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("missing object information in record")
	}
	objectKey, ok := objectInfo["key"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid object key")
	}

	return bucketName, objectKey, nil
}

// Mock function to fetch and check file type, returning first line if textual
func fetchAndCheckFile(bucket string, key string) (bool, string, error) {
	// Simulate downloading the file and checking if it is textual
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	// Create an S3 client using the loaded configuration
	s3.NewFromConfig(cfg)

	// Create a downloader passing it the S3 client
	downloader := s3manager.NewDownloader(s3.NewFromConfig(cfg))

	buffer := s3manager.NewWriteAtBuffer([]byte{})

	_, err = downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return false, "", fmt.Errorf("failed to download S3 object: %w", err)
	}

	content := string(buffer.Bytes())
	scanner := bufio.NewScanner(strings.NewReader(content))
	isTextual := true
	firstLine := ""

	if scanner.Scan() {
		firstLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		isTextual = false
		firstLine = ""
	}

	return isTextual, firstLine, nil
}
