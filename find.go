package tmdb

import (
	"context"
	"fmt"
	"net/url"
)

type FindResponse struct {
	MovieResults []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		ID               int     `json:"id"`
		Title            string  `json:"title"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		PosterPath       string  `json:"poster_path"`
		GenreIds         []int   `json:"genre_ids"`
		Popularity       float64 `json:"popularity"`
		ReleaseDate      string  `json:"release_date"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
	} `json:"movie_results"`
}

// FindByIMDB looks up a movie by its IMDB ID and returns the TMDB ID
// and basic info. Returns nil if not found.
func (c *Client) FindByIMDB(ctx context.Context, imdbID string) (*FindResponse, error) {
	v := url.Values{}
	v.Set("api_key", c.key)
	v.Set("external_source", "imdb_id")
	u := fmt.Sprintf("%s/find/%s?%s", tmdbURL, imdbID, v.Encode())
	r := new(FindResponse)

	if err := c.get(ctx, u, r); err != nil {
		return nil, err
	}
	return r, nil
}

type CastMember struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Character   string `json:"character"`
	Order       int    `json:"order"`
	ProfilePath string `json:"profile_path"`
}

type Credits struct {
	Cast []CastMember `json:"cast"`
}

type MovieWithCreditsResponse struct {
	MovieResponse
	Credits Credits `json:"credits"`
}

// MovieWithCredits fetches movie details with cast credits in a single request.
func (c *Client) MovieWithCredits(ctx context.Context, id int) (*MovieWithCreditsResponse, error) {
	v := url.Values{}
	v.Set("api_key", c.key)
	v.Set("append_to_response", "credits")
	u := fmt.Sprintf("%s/movie/%v?%s", tmdbURL, id, v.Encode())
	r := new(MovieWithCreditsResponse)

	if err := c.get(ctx, u, r); err != nil {
		return nil, err
	}
	return r, nil
}
