package tmdb

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// newTestClient returns a Client pointed at an httptest server that serves
// body for every request. Each request's URL is sent on the returned channel.
func newTestClient(t *testing.T, body string) (*Client, <-chan *url.URL) {
	t.Helper()
	reqs := make(chan *url.URL, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqs <- r.URL
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	c := New("test-key")
	c.baseURL = srv.URL
	return c, reqs
}

func checkQuery(t *testing.T, u *url.URL, want map[string]string) {
	t.Helper()
	q := u.Query()
	for k, v := range want {
		if got := q.Get(k); got != v {
			t.Errorf("query %q = %q; want %q", k, got, v)
		}
	}
}

func TestSearchMoviesOptionsValues(t *testing.T) {
	v := (&SearchMoviesOptions{Query: "alien"}).values()
	if got := v.Encode(); got != "include_adult=false&query=alien" {
		t.Errorf("minimal options encoded to %q", got)
	}

	v = (&SearchMoviesOptions{
		Lang:               "en-US",
		Query:              "alien",
		Page:               2,
		Region:             "US",
		Year:               1979,
		PrimaryReleaseYear: 1980,
	}).values()
	want := map[string]string{
		"language":             "en-US",
		"query":                "alien",
		"page":                 "2",
		"region":               "US",
		"year":                 "1979",
		"primary_release_year": "1980",
	}
	for k, w := range want {
		if got := v.Get(k); got != w {
			t.Errorf("%s = %q; want %q", k, got, w)
		}
	}
}

func TestSearchTVOptionsValues(t *testing.T) {
	v := (&SearchTVOptions{
		Lang:             "de",
		Query:            "dark",
		Page:             3,
		FirstAirDateYear: 2017,
	}).values()
	want := map[string]string{
		"language":            "de",
		"query":               "dark",
		"page":                "3",
		"first_air_date_year": "2017",
		"include_adult":       "false",
	}
	for k, w := range want {
		if got := v.Get(k); got != w {
			t.Errorf("%s = %q; want %q", k, got, w)
		}
	}
}

func TestSearchMovies(t *testing.T) {
	c, reqs := newTestClient(t, `{
		"page": 1,
		"results": [{"id": 245891, "title": "John Wick", "release_date": "2014-10-22", "genre_ids": [28, 53]}],
		"total_pages": 1,
		"total_results": 1
	}`)

	resp, err := c.SearchMovies(context.Background(), &SearchMoviesOptions{Query: "john wick", Year: 2014})
	if err != nil {
		t.Fatal(err)
	}
	u := <-reqs
	if u.Path != "/search/movie" {
		t.Errorf("path = %q", u.Path)
	}
	checkQuery(t, u, map[string]string{"api_key": "test-key", "query": "john wick", "year": "2014"})

	if resp.TotalResults != 1 || len(resp.Results) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	r := resp.Results[0]
	if r.ID != 245891 || r.Title != "John Wick" || len(r.GenreIds) != 2 {
		t.Errorf("unexpected result: %+v", r)
	}
}

func TestMovie(t *testing.T) {
	c, reqs := newTestClient(t, `{"id": 245891, "imdb_id": "tt2911666", "runtime": 101, "genres": [{"id": 28, "name": "Action"}]}`)

	m, err := c.Movie(context.Background(), 245891)
	if err != nil {
		t.Fatal(err)
	}
	u := <-reqs
	if u.Path != "/movie/245891" {
		t.Errorf("path = %q", u.Path)
	}
	checkQuery(t, u, map[string]string{"api_key": "test-key"})
	if m.IMDBID != "tt2911666" || m.Runtime != 101 || len(m.Genres) != 1 || m.Genres[0].Name != "Action" {
		t.Errorf("unexpected movie: %+v", m)
	}
}

func TestMovieWithCredits(t *testing.T) {
	c, reqs := newTestClient(t, `{
		"id": 245891,
		"title": "John Wick",
		"credits": {"cast": [{"id": 6384, "name": "Keanu Reeves", "character": "John Wick", "order": 0}]}
	}`)

	m, err := c.MovieWithCredits(context.Background(), 245891)
	if err != nil {
		t.Fatal(err)
	}
	u := <-reqs
	if u.Path != "/movie/245891" {
		t.Errorf("path = %q", u.Path)
	}
	checkQuery(t, u, map[string]string{"append_to_response": "credits"})
	if m.Title != "John Wick" {
		t.Errorf("embedded MovieResponse not decoded: %+v", m.MovieResponse)
	}
	if len(m.Credits.Cast) != 1 || m.Credits.Cast[0].Name != "Keanu Reeves" {
		t.Errorf("unexpected credits: %+v", m.Credits)
	}
}

func TestFindByIMDB(t *testing.T) {
	c, reqs := newTestClient(t, `{"movie_results": [{"id": 245891, "title": "John Wick"}]}`)

	r, err := c.FindByIMDB(context.Background(), "tt2911666")
	if err != nil {
		t.Fatal(err)
	}
	u := <-reqs
	if u.Path != "/find/tt2911666" {
		t.Errorf("path = %q", u.Path)
	}
	checkQuery(t, u, map[string]string{"external_source": "imdb_id", "api_key": "test-key"})
	if len(r.MovieResults) != 1 || r.MovieResults[0].ID != 245891 {
		t.Errorf("unexpected results: %+v", r)
	}
}

func TestSearchTV(t *testing.T) {
	c, reqs := newTestClient(t, `{"page": 1, "results": [{"id": 1399, "name": "Game of Thrones", "origin_country": ["US"]}], "total_results": 1}`)

	r, err := c.SearchTV(context.Background(), &SearchTVOptions{Query: "game of thrones"})
	if err != nil {
		t.Fatal(err)
	}
	u := <-reqs
	if u.Path != "/search/tv" {
		t.Errorf("path = %q", u.Path)
	}
	checkQuery(t, u, map[string]string{"query": "game of thrones"})
	if len(r.Results) != 1 || r.Results[0].Name != "Game of Thrones" || r.Results[0].OriginCountry[0] != "US" {
		t.Errorf("unexpected results: %+v", r)
	}
}

func TestTV(t *testing.T) {
	c, reqs := newTestClient(t, `{
		"id": 1399,
		"name": "Game of Thrones",
		"number_of_seasons": 8,
		"last_episode_to_air": {"episode_number": 6, "season_number": 8},
		"next_episode_to_air": null,
		"seasons": [{"season_number": 1, "episode_count": 10}],
		"external_ids": {"imdb_id": "tt0944947", "tvdb_id": 121361}
	}`)

	tv, err := c.TV(context.Background(), 1399)
	if err != nil {
		t.Fatal(err)
	}
	u := <-reqs
	if u.Path != "/tv/1399" {
		t.Errorf("path = %q", u.Path)
	}
	checkQuery(t, u, map[string]string{"append_to_response": "external_ids"})
	if tv.NumberOfSeasons != 8 || len(tv.Seasons) != 1 || tv.Seasons[0].EpisodeCount != 10 {
		t.Errorf("unexpected show: %+v", tv)
	}
	if tv.LastEpisodeToAir == nil || tv.LastEpisodeToAir.EpisodeNumber != 6 {
		t.Errorf("unexpected last episode: %+v", tv.LastEpisodeToAir)
	}
	if tv.NextEpisodeToAir != nil {
		t.Errorf("expected nil next episode, got %+v", tv.NextEpisodeToAir)
	}
	if tv.ExternalIDs == nil || tv.ExternalIDs.IMDBID != "tt0944947" || tv.ExternalIDs.TVDBID != 121361 {
		t.Errorf("unexpected external ids: %+v", tv.ExternalIDs)
	}
}

func TestIncludeAdultTrueIsSent(t *testing.T) {
	if got := (&SearchMoviesOptions{Query: "x", IncludeAdult: true}).values().Get("include_adult"); got != "true" {
		t.Errorf("movies include_adult = %q; want true", got)
	}
	if got := (&SearchTVOptions{Query: "x", IncludeAdult: true}).values().Get("include_adult"); got != "true" {
		t.Errorf("tv include_adult = %q; want true", got)
	}
}

func TestFindByIMDBEscapesID(t *testing.T) {
	c, reqs := newTestClient(t, `{"movie_results": []}`)
	if _, err := c.FindByIMDB(context.Background(), "tt1/../../movie/1"); err != nil {
		t.Fatal(err)
	}
	if u := <-reqs; u.Path != "/find/tt1/../../movie/1" || u.RawPath != "/find/tt1%2F..%2F..%2Fmovie%2F1" {
		t.Errorf("id was not path-escaped: path=%q raw=%q", u.Path, u.RawPath)
	}
}

func TestErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status_code": 7, "status_message": "Invalid API key: You must be granted a valid key."}`))
	}))
	defer srv.Close()
	c := New("bad")
	c.baseURL = srv.URL

	_, err := c.Movie(context.Background(), 1)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("err = %v (%T); want *APIError", err, err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized || apiErr.StatusMessage == "" {
		t.Errorf("unexpected APIError: %+v", apiErr)
	}
}

func TestContextIsHonored(t *testing.T) {
	c, _ := newTestClient(t, `{}`)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.Movie(ctx, 1); err == nil {
		t.Fatal("expected an error from a canceled context")
	}
}
