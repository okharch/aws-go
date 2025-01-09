# AWS Go Learning Projects

This repository is dedicated to learning and experimenting with the [AWS SDK for Go](https://aws.github.io/aws-sdk-go-v2/). It serves as a collection of small, focused projects that demonstrate the capabilities of AWS services in conjunction with the Go programming language.

## Purpose
The primary objective of this repository is to:
- Provide hands-on examples of integrating AWS services using the Go SDK.
- Create a reference for common patterns, best practices, and use cases while working with AWS in Go.
- Serve as a learning resource for developers exploring AWS services in a Go-based environment.

## Projects

### 1. **S3 Upload**
#### Description:
The first project in this repository demonstrates handling AWS S3 and SQS event-driven workflows. It receives event notifications triggered by uploads to an S3 bucket, processes them, and performs the following:
- Fetches details about the uploaded object from the S3 event.
- Downloads the S3 object and determines if it is textual.
- Logs the first line of textual files or skips others.
- Logs all operations for debugging and monitoring purposes.

#### Key AWS Services:
- **Amazon S3 (Simple Storage Service):** Used for storage of files and to host event notifications for uploads.
- **Amazon SQS (Simple Queue Service):** Acts as a trigger for receiving and processing S3 bucket event notifications.

### Future Additions
This repository will expand to include more learning projects for other AWS services, such as DynamoDB, Lambda, ECS, and more.

---

## Setting Up

### Prerequisites
Before running any examples in this repository, ensure that you:
1. **Create your own AWS identity and access key pair**:
   You can configure this by running:
   ```bash
   aws configure
   ```
   This command will prompt you to enter:
    - AWS Access Key ID
    - AWS Secret Access Key
    - The default AWS Region (e.g., `us-east-1`)
    - Output format (optional)

2. **Create the necessary AWS resources**:
   Examples in this repository depend on AWS resources such as:
    - **S3 buckets**
    - **SQS queues**
    - Other AWS services based on the specific example.

   Ensure you create these resources in your AWS account using either:
    - The AWS Management Console
    - AWS CLI commands
    - Infrastructure as Code (IaC) tools like AWS CloudFormation or Terraform.

### Steps to Run
1. Clone this repository.
   ```bash
   git clone git@github.com:okharch/aws-go.git
   cd aws-go
   ```
2. Set up your AWS credentials and region as mentioned in the prerequisites.
3. Navigate to the specific project directory and follow any additional setup instructions (if provided).
4. Run the Go programs to interact with the configured AWS services.

---

## Contributing
Contributions are welcome! If you have any suggestions or new project ideas to add, feel free to:
1. Submit a pull request.
2. Open an issue for discussion.

---

Happy coding with AWS and Go!