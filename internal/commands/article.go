package commands

import (
	"fmt"
	"os"
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


func HandleArticleCreate(){
}
