package main

// movie represents data about a record album.
type movie struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	ReleasedYear int      `json:"released_year"`
	Rating       float64  `json:"rating"`
	Genres       []string `json:"genres"`
}

func main() {
}