package lamdba

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"os"
)

func DisplayInfo() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := lambda.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	result, err := svc.ListFunctions(nil)
	if err != nil {
		fmt.Println("Cannot list functions")
		os.Exit(0)
	}

	fmt.Println("Functions:")

	for _, f := range result.Functions {
		fmt.Println("Name:        " + aws.StringValue(f.FunctionName))
		fmt.Println("Description: " + aws.StringValue(f.Description))
		fmt.Println("")
	}
}
