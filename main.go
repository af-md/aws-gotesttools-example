package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"aws-gotesttools-example/handler"
)

func main() {

	// create config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	//create handler
	h := handler.NewHandler(cfg)

	// register handler to the http server
	http.HandleFunc("/", h.HandleRequest)

	// the httpadapter package will adapt aws event to a http request
	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}
