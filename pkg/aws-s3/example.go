package awss3

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func Example() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")))

	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	client := s3.NewFromConfig(cfg)

	type result struct {
		Output *s3.PutObjectOutput
		Err    error
	}

	results := make(chan result, 2)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("bucket-serv-micro"),
			Key:    aws.String("PruebaFinalll2.0s"),
			Body:   strings.NewReader("Pruba de contenidos"),
		})
		results <- result{Output: output, Err: err}
	}()

	// go func() {
	// 	defer wg.Done()
	// 	output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
	// 		Bucket: aws.String("bucket-serv-micro"),
	// 		Key:    aws.String("pruebasdfasd"),
	// 		Body:   strings.NewReader("bar body contentde prueba"),
	// 	})
	// 	results <- result{Output: output, Err: err}
	// }()

	wg.Wait()

	close(results)

	for result := range results {
		if result.Err != nil {
			log.Printf("error: %v", result.Err)
			continue
		}
		fmt.Printf("etag: %v", aws.ToString(result.Output.ETag))
	}
}
