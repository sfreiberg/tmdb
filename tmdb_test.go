package tmdb

import (
	"context"
	"os"
	"testing"
)

func TestMovie(t *testing.T) {
	apiKey := os.Getenv("TMDB_KEY")
	if apiKey == "" {
		t.Fatal("TMDB_KEY is not set")
	}
	client := New(apiKey)
	opts := &SearchMoviesOptions{Query: "john wick", Year: 2014}
	resp, err := client.SearchMovies(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}
	if resp.TotalResults != 1 {
		t.Fatalf("Expected 1 TotalResults; got %v\n", resp.TotalResults)
	}
	movie, err := client.Movie(context.Background(), resp.Results[0].ID)
	if err != nil {
		t.Fatalf("Unable to get movie: %s\n", err)
	}
	if movie.IMDBID != "tt2911666" {
		t.Fatalf("IMDBID == %s; expected tt2911666\n", movie.IMDBID)
	}
}
