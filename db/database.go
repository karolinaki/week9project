package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// NewDynamoDB initializes a new DynamoDB client.
func NewDynamoDB() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"), // Replace with your desired region
		Credentials: credentials.NewStaticCredentials("ASIAW7OMQ3DYAVT7CCEQ", "PO1Y82NGpTGWs9wqmffpwIwT03MU3JnYYbikjMzE", "IQoJb3JpZ2luX2VjEDYaCXVzLWVhc3QtMSJFMEMCH2UNCg+Kyh1zpvC95d/noeyIZ+DaEUFJ6diIrRksuWICICpIDFKiXzFkcZKg1Y3MO7PoZ0e9caW6qeINaBE2U8+AKvUCCF8QABoMNDc5ODU0NjQ3NTM2Igw4dixdg6jALZu9veQq0gKqvLQyg9mt5L6iZ/PLN9FftjXeGBDzpR9D/riZnSOr3//SoRBifyjFTYgR/jzi4F9fTNkcwtumgolNDRgfKd/7+5H4NE6cw6v60LB+FjUpavK2D4nW0E2kfnaNmmdDQK3QKFV1Xc5HJhHf7+Rf7v23kZpZXPm9o7h0jcso0ewysr5qsQGsndCNxl/+DELkThSvhIQrihyZALB7iRQpefzJfC/dNTCVdpFQnqeR5NDWBsctmvAL2CZQ2AnyFPdaDcS+EMHoNU1T51AfCKgVC5NvVupeO7G8GvEVVLhwGopwSrlns9e7BdyaFeMUuBUogCpsgsmgUxB4PKfv6yeoCdX770XKQj/dh+A8irRUzoTZiyVOjMX+C5U0DfralpjWEBL2hRMBYx61C7q2cICdwGo9t2hdMybraVA6kJp+WNXOIZs4wVyIVyYdBXZNDcONN7XNFTDO0ZijBjqpAbU0O96vtqNoj5qH5ppTpnr9zO5m9fhqeZTA78W81pEFcWbSar6LeAJNVaqeiUYFBI3XEBoB8Kpn6v82gLfg9RCsEy2g7vQlKvJ8+nnE4L4t0PxTTxrQdG8o0l/Ami77uuOyYfBsFA6Nyamhk6GTyu/VD0SP/cVOiiRCu3KQK4R3O1mfxLXfw6mze21eJx0L6IU0eM2oKX7RclnOc9DXw3rwNLtXNJt3KaY="),
	})

	if err != nil {
		// Handle session creation error
		fmt.Println("Error creating session:", err)
	}

	return dynamodb.New(sess)
}
