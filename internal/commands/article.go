package commands

import (
	"fmt"
	"net/http"
	"os"
	"ssb/internal/commands/utils"
	"ssb/internal/schemas"
)

func HandleArticle(args []string) {
	switch args[0] {
	case "create":
		fmt.Println("hello from create article")
	case "delete":
		fmt.Println("hello from delete article")
	default:
		fmt.Println("expected one of 'create' or 'delete'")
		os.Exit(1)
	}
}

/*
title := "title"
body := "body"
data := schemas.ArticleCreateSchema {
	Title: title,
	Body: body,
}
*/

func HandleArticleCreate(data schemas.ArticleCreateSchema, client HTTPClient) error {
	ep, err := utils.BuildEndpoint("articles")
	if err != nil {
		return fmt.Errorf("could not build endpoint: %w", err)
	}

	req, err := utils.NewRequestBuilder(http.MethodPost, ep).
		WithAuth().
		WithJSON(data).
		Build()

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not make request: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("could not create article: %w", err)
	}

	return nil
}
