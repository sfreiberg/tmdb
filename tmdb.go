package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const tmdbURL = "https://api.themoviedb.org/3"

type Client struct {
	key string
}

func New(key string) *Client {
	return &Client{key: key}
}

type SearchMoviesOptions struct {
	Lang               string // optional
	Query              string // required
	Page               int    // optional
	IncludeAdult       bool   // optional
	Region             string // optional
	Year               int    // optional
	PrimaryReleaseYear int    // optional
}

func (o *SearchMoviesOptions) values() url.Values {
	v := url.Values{}
	if o.Lang != "" {
		v.Set("language", o.Lang)
	}
	if o.Query != "" {
		v.Set("query", o.Query)
	}
	if o.Page != 0 {
		v.Set("page", fmt.Sprintf("%v", o.Page))
	}
	if !o.IncludeAdult {
		v.Set("include_adult", "false")
	}
	if o.Region != "" {
		v.Set("region", o.Region)
	}
	if o.Year != 0 {
		v.Set("year", fmt.Sprintf("%v", o.Year))
	}
	if o.PrimaryReleaseYear != 0 {
		v.Set("primary_release_year", fmt.Sprintf("%v", o.PrimaryReleaseYear))
	}
	return v
}

type SearchMoviesResponse struct {
	Page    int `json:"page"`
	Results []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		GenreIds         []int   `json:"genre_ids"`
		ID               int     `json:"id"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		Popularity       float64 `json:"popularity"`
		PosterPath       string  `json:"poster_path"`
		ReleaseDate      string  `json:"release_date,omitempty"`
		Title            string  `json:"title"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

func (c *Client) SearchMovies(ctx context.Context, opts *SearchMoviesOptions) (*SearchMoviesResponse, error) {
	v := opts.values()
	v.Set("api_key", c.key)
	u := fmt.Sprintf("%s/search/movie?%s", tmdbURL, v.Encode())
	r := new(SearchMoviesResponse)

	if err := c.get(ctx, u, r); err != nil {
		return nil, err
	}
	return r, nil
}

type MovieResponse struct {
	Adult               bool   `json:"adult"`
	BackdropPath        string `json:"backdrop_path"`
	BelongsToCollection struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		PosterPath   string `json:"poster_path"`
		BackdropPath string `json:"backdrop_path"`
	} `json:"belongs_to_collection"`
	Budget int `json:"budget"`
	Genres []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage            string  `json:"homepage"`
	ID                  int     `json:"id"`
	IMDBID              string  `json:"imdb_id"`
	OriginalLanguage    string  `json:"original_language"`
	OriginalTitle       string  `json:"original_title"`
	Overview            string  `json:"overview"`
	Popularity          float64 `json:"popularity"`
	PosterPath          string  `json:"poster_path"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate     string `json:"release_date"`
	Revenue         int    `json:"revenue"`
	Runtime         int    `json:"runtime"`
	SpokenLanguages []struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

func (c *Client) Movie(ctx context.Context, id int) (*MovieResponse, error) {
	v := url.Values{}
	v.Set("api_key", c.key)
	u := fmt.Sprintf("%s/movie/%v?%s", tmdbURL, id, v.Encode())
	r := new(MovieResponse)

	if err := c.get(ctx, u, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) get(ctx context.Context, url string, respObj interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(respObj)
}
