package commands

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"ssb/internal/commands/utils"
	"ssb/internal/schemas"
)

func HandleArticle(args []string) {
	switch args[0] {
	case "init":
		if len(args) < 2 {
			fmt.Println("provide a name for the article directory")
			os.Exit(1)
		}
		HandleArticleInit(args[1])

	case "create":
		fmt.Println("not implemented")
	case "update":
		fmt.Println("not implemented")
	case "delete":
		fmt.Println("not implemented")
	default:
		fmt.Println("expected one of 'init, create, update' or 'delete'")
		os.Exit(1)
	}
}

func HandleArticleInit(name string) error {
	if err := os.Mkdir(name, 0755); err != nil {
		return fmt.Errorf("could not create dir: %w", err)
	}

	imgDir := filepath.Join(name, "images")
	if err := os.Mkdir(imgDir, 0755); err != nil {
		return fmt.Errorf("could not create 'images' directory: %w", err)
	}

	doc := filepath.Join(name, "article.md")
	if err := createEmptyFile(doc); err != nil {
		return fmt.Errorf("could not create article %s: %w", doc, err)
	}

	metadata := filepath.Join(name, "metadata.json")
	if err := createEmptyFile(metadata); err != nil {
		return fmt.Errorf("could not create metadata file %s: %w", metadata, err)
	}

	return nil
}

func createEmptyFile(name string) error {
	d := []byte("")
	if err := os.WriteFile(name, d, 0644); err != nil {
		return err
	}
	return nil
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
