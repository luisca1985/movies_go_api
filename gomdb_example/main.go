package main

import (
	"fmt"
	"strings"

	"github.com/eefret/gomdb"
)

func main() {
	api := gomdb.Init("4733c4c1")
	query := &gomdb.QueryData{Title: "Macbeth", SearchType: gomdb.MovieSearch}
	res, err := api.Search(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.Search)

	// query = &gomdb.QueryData{Title: "Macbeth", Year: "2015"}
	query = &gomdb.QueryData{Title: "Gladiator"}
	res2, err := api.MovieByTitle(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Title")
	fmt.Println(res2.Title)
	fmt.Println("Year")
	fmt.Println(res2.Year)
	fmt.Println("Rated")
	fmt.Println(res2.Rated)
	fmt.Println("Genre")
	fmt.Println(res2.Genre)
	fmt.Println("Genre slice")
	fmt.Println(strings.Split(strings.ReplaceAll(res2.Genre, " ", ""), ","))

	res3, err := api.MovieByImdbID("tt2884018")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res3)
}
