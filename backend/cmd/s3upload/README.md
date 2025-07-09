# S3 Upload CLI

A simple command-line tool for uploading files to S3-compatible storage services.

## Features

- Upload files to any S3-compatible storage service (AWS S3, MinIO, Ceph, etc.)
- Configurable endpoint URL, region, and credentials
- Supports both HTTP and HTTPS
- Path-style addressing for compatibility with most S3-compatible APIs

## Installation

First, install the required dependencies:

```bash
go get github.com/aws/aws-sdk-go-v2/aws
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/credentials
go get github.com/aws/aws-sdk-go-v2/service/s3
```

Then build the tool:

```bash
go build -o s3upload cmd/s3upload/main.go
```

## Usage

```
S3 Upload CLI - Upload files to S3-compatible storage

Usage:
  s3upload [options]

Options:
  -access-key string
        Access key (required)
  -bucket string
        S3 bucket name (required)
  -endpoint string
        S3-compatible API endpoint (required)
  -file string
        Path to the file you want to upload (required)
  -help
        Show help
  -key string
        Object key (name) in the bucket. If not specified, the filename will be used
  -region string
        AWS region (default "us-east-1")
  -secret-key string
        Secret key (required)
  -ssl
        Use SSL/TLS for connection (default true)
```

## Examples

### Uploading to MinIO

```bash
./s3upload \
  -endpoint=play.min.io \
  -bucket=mybucket \
  -access-key=Q3AM3UQ867SPQQA43P2F \
  -secret-key=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG \
  -file=./myfile.txt
```

### Uploading to AWS S3

```bash
./s3upload \
  -endpoint=s3.amazonaws.com \
  -region=us-west-2 \
  -bucket=mybucket \
  -access-key=AKIAIOSFODNN7EXAMPLE \
  -secret-key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY \
  -file=./myfile.txt \
  -key=uploads/myfile.txt
```

### Uploading to a local MinIO instance without SSL

```bash
./s3upload \
  -endpoint=localhost:9000 \
  -bucket=mybucket \
  -access-key=minioadmin \
  -secret-key=minioadmin \
  -file=./myfile.txt \
  -ssl=false
```

## Development

To modify or extend this tool:

1. Add any additional features to `main.go`
2. Update the `printUsage()` function to document new options
3. Build using `go build` 