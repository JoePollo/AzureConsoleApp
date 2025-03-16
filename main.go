package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func ErrorHandler(object string, err error) error {
	return fmt.Errorf("Failed to build %s due to error:\n%w", object, err)
}

func GetCredentials() (*azidentity.DefaultAzureCredential, error) {

	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, ErrorHandler("NewDefaultAzureCredential", err)
	}
	return creds, nil
}

func GetClient(environment string, serviceBusHostName string, credentials *azidentity.DefaultAzureCredential) (*azservicebus.Client, error) {
	serviceBusClient, err := azservicebus.NewClient(serviceBusHostName, credentials, nil)
	if err != nil {
		return nil, ErrorHandler("azservicebus.NewClient", err)
	}
	return serviceBusClient, nil
}

func SendMessage(message string, client *azservicebus.Client, queueName string) error {
	sender, err := client.NewSender(queueName, nil)
	if err != nil {
		return ErrorHandler("client.NewSender", err)
	}

	defer sender.Close(context.TODO())

	serviceBusMessage := &azservicebus.Message{
		Body: []byte(message),
	}

	err = sender.SendMessage(context.TODO(), serviceBusMessage, nil)
	if err != nil {
		return ErrorHandler("sender.SendMessage", err)
	}
	return nil
}

func main() {
	Env := os.Getenv("ENV")
	ServiceBusHostName := fmt.Sprintf("sbns-yym-%s-usce.servicebus.windows.net", Env)
	QueueName := fmt.Sprintf("sbq-yym-%s-usce", Env)
	creds, err := GetCredentials()
	if err != nil {
		log.Fatal(err)
	}
	serviceBusClient, err := GetClient(Env, ServiceBusHostName, creds)
	if err != nil {
		log.Fatal(err)
	}
	err = SendMessage("super cool message", serviceBusClient, QueueName)
	if err != nil {
		log.Fatal(err)
	}

}
