package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// Define command line flags
	endpoint := flag.String("endpoint", "", "S3-compatible API endpoint (required)")
	region := flag.String("region", "us-east-1", "AWS region")
	bucket := flag.String("bucket", "", "S3 bucket name (required)")
	accessKey := flag.String("access-key", "", "Access key (required)")
	secretKey := flag.String("secret-key", "", "Secret key (required)")
	filePath := flag.String("file", "", "Path to the file you want to upload (required)")
	objectKey := flag.String("key", "", "Object key (name) in the bucket. If not specified, the filename will be used")
	useSSL := flag.Bool("ssl", true, "Use SSL/TLS for connection")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()

	// Show help if requested or missing required parameters
	if *help || *endpoint == "" || *bucket == "" || *accessKey == "" || *secretKey == "" || *filePath == "" {
		printUsage()
		return
	}

	// Initialize the S3 client
	client, err := createS3Client(*endpoint, *region, *accessKey, *secretKey, *useSSL)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	// Default to filename if object key is not specified
	key := *objectKey
	if key == "" {
		key = filepath.Base(*filePath)
	}

	// Upload the file
	err = uploadFile(client, *bucket, key, *filePath)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}

	fmt.Printf("Successfully uploaded file '%s' to bucket '%s' with key '%s'\n", *filePath, *bucket, key)
}

func printUsage() {
	fmt.Println("S3 Upload CLI - Upload files to S3-compatible storage")
	fmt.Println("\nUsage:")
	fmt.Println("  s3upload [options]")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExample:")
	fmt.Println("  s3upload -endpoint=https://s3.example.com -bucket=my-bucket \\")
	fmt.Println("           -access-key=myaccesskey -secret-key=mysecretkey \\")
	fmt.Println("           -file=./myfile.txt -key=uploads/myfile.txt")
}

// Create an S3 client with custom configuration
func createS3Client(endpoint, region, accessKey, secretKey string, useSSL bool) (*s3.Client, error) {
	// Update endpoint URL scheme if needed
	if !strings.HasPrefix(endpoint, "http") {
		if useSSL {
			endpoint = "https://" + endpoint
		} else {
			endpoint = "http://" + endpoint
		}
	}

	// Create AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config: %w", err)
	}

	// Create and return the S3 client with custom endpoint
	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true // Use path-style addressing for compatibility with most S3-compatible APIs
	}), nil
}

// Upload a file to S3-compatible storage
func uploadFile(client *s3.Client, bucket, key, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info for content length
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Create the upload input
	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: aws.Int64(fileInfo.Size()),
	}

	// Perform the upload
	_, err = client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}
