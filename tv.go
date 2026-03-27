package tmdb

import (
	"context"
	"fmt"
	"net/url"
)

type SearchTVOptions struct {
	Lang             string // optional
	Query            string // required
	Page             int    // optional
	IncludeAdult     bool   // optional
	FirstAirDateYear int    // optional
}

func (o *SearchTVOptions) values() url.Values {
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
	if o.FirstAirDateYear != 0 {
		v.Set("first_air_date_year", fmt.Sprintf("%v", o.FirstAirDateYear))
	}
	return v
}

type SearchTVResponse struct {
	Page    int `json:"page"`
	Results []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		GenreIds         []int   `json:"genre_ids"`
		ID               int     `json:"id"`
		OriginalLanguage string  `json:"original_language"`
		OriginalName     string  `json:"original_name"`
		Overview         string  `json:"overview"`
		Popularity       float64 `json:"popularity"`
		PosterPath       string  `json:"poster_path"`
		FirstAirDate     string  `json:"first_air_date,omitempty"`
		Name             string  `json:"name"`
		OriginCountry    []string `json:"origin_country"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

func (c *Client) SearchTV(ctx context.Context, opts *SearchTVOptions) (*SearchTVResponse, error) {
	v := opts.values()
	v.Set("api_key", c.key)
	u := fmt.Sprintf("%s/search/tv?%s", tmdbURL, v.Encode())
	r := new(SearchTVResponse)

	if err := c.get(ctx, u, r); err != nil {
		return nil, err
	}
	return r, nil
}

type TVEpisode struct {
	AirDate       string  `json:"air_date"`
	EpisodeNumber int     `json:"episode_number"`
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Overview      string  `json:"overview"`
	SeasonNumber  int     `json:"season_number"`
	ShowID        int     `json:"show_id"`
	StillPath     string  `json:"still_path"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
}

type TVSeason struct {
	AirDate      string  `json:"air_date"`
	EpisodeCount int     `json:"episode_count"`
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	SeasonNumber int     `json:"season_number"`
	VoteAverage  float64 `json:"vote_average"`
}

type TVExternalIDs struct {
	IMDBID      string `json:"imdb_id"`
	TVDBID      int    `json:"tvdb_id"`
	FacebookID  string `json:"facebook_id"`
	InstagramID string `json:"instagram_id"`
	TwitterID   string `json:"twitter_id"`
	WikidataID  string `json:"wikidata_id"`
}

type TVResponse struct {
	Adult        bool   `json:"adult"`
	BackdropPath string `json:"backdrop_path"`
	CreatedBy    []struct {
		ID          int    `json:"id"`
		CreditID    string `json:"credit_id"`
		Name        string `json:"name"`
		Gender      int    `json:"gender"`
		ProfilePath string `json:"profile_path"`
	} `json:"created_by"`
	EpisodeRunTime []int  `json:"episode_run_time"`
	FirstAirDate   string `json:"first_air_date"`
	Genres         []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage         string      `json:"homepage"`
	ID               int         `json:"id"`
	InProduction     bool        `json:"in_production"`
	Languages        []string    `json:"languages"`
	LastAirDate      string      `json:"last_air_date"`
	LastEpisodeToAir *TVEpisode  `json:"last_episode_to_air"`
	NextEpisodeToAir *TVEpisode  `json:"next_episode_to_air"`
	Name             string      `json:"name"`
	Networks         []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		LogoPath      string `json:"logo_path"`
		OriginCountry string `json:"origin_country"`
	} `json:"networks"`
	NumberOfEpisodes int      `json:"number_of_episodes"`
	NumberOfSeasons  int      `json:"number_of_seasons"`
	OriginCountry    []string `json:"origin_country"`
	OriginalLanguage string   `json:"original_language"`
	OriginalName     string   `json:"original_name"`
	Overview         string   `json:"overview"`
	Popularity       float64  `json:"popularity"`
	PosterPath       string   `json:"poster_path"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	Seasons     []TVSeason     `json:"seasons"`
	Status      string         `json:"status"`
	Tagline     string         `json:"tagline"`
	Type        string         `json:"type"`
	VoteAverage float64        `json:"vote_average"`
	VoteCount   int            `json:"vote_count"`
	ExternalIDs *TVExternalIDs `json:"external_ids,omitempty"`
}

// TV fetches TV show details by ID. It uses append_to_response=external_ids
// to include external IDs (TVDB, IMDB, etc.) in a single request.
func (c *Client) TV(ctx context.Context, id int) (*TVResponse, error) {
	v := url.Values{}
	v.Set("api_key", c.key)
	v.Set("append_to_response", "external_ids")
	u := fmt.Sprintf("%s/tv/%v?%s", tmdbURL, id, v.Encode())
	r := new(TVResponse)

	if err := c.get(ctx, u, r); err != nil {
		return nil, err
	}
	return r, nil
}
