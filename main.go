package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
)

type Repository struct {
	Projects struct {
		Nodes []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			ForksCount  int    `json:"forksCount"`
		} `json:"nodes"`
	} `json:"projects"`
}

type Service struct {
	Name       string
	ForksCount int
}

// 1. repository layer
func getRepository() (resData Repository) {
	GITLAB_API_URL := goDotEnvVariable("GITLAB_API_URL")

	client := graphql.NewClient(GITLAB_API_URL)
	req := graphql.NewRequest(`
		query last_projects($n: Int = 5) {
			projects(last:$n) {
				nodes {
					name
					description
					forksCount
				}
			}
	  	}
	`)

	// set any variables
	req.Var("n", 10)

	// get gitlab access
	GIT_ACCESS := goDotEnvVariable("GITLAB_ACCESS")
	// set header fields
	req.Header.Set("Authorization", "Bearer "+GIT_ACCESS)
	req.Header.Set("content-type", "application/json")

	// define a Context for the request
	ctx := context.Background()

	if err := client.Run(ctx, req, &resData); err != nil {
		panic(err)
	}

	return
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// 2. service layer
func getService() Service {
	resData := getRepository()

	nameString := ""
	forksCount := 0
	nodes := resData.Projects.Nodes
	for _, value := range nodes {
		nameString += value.Name + ", "
		forksCount += value.ForksCount
	}

	result := Service{
		Name:       nameString,
		ForksCount: forksCount,
	}

	return result
}

func main() {
	service := getService()

	// 4. print results
	fmt.Println("Name: ", service.Name)
	fmt.Println("ForksCount: ", service.ForksCount)
}
