package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

type OMDbService struct {
	client  *resty.Client
	apiKey  string
	baseURL string
}

type MovieInfo struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	Response   string `json:"Response"`
	Error      string `json:"Error"`
}

func NewOMDbService() *OMDbService {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(2 * time.Second)

	return &OMDbService{
		client:  client,
		apiKey:  os.Getenv("OMDB_API_KEY"),
		baseURL: "http://www.omdbapi.com",
	}
}

func (s *OMDbService) GetMovieByTitle(title string) (*MovieInfo, error) {
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"t":      title,
			"apikey": s.apiKey,
			"plot":   "full",
		}).
		Get(s.baseURL)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var movieInfo MovieInfo
	if err := json.Unmarshal(resp.Body(), &movieInfo); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if movieInfo.Response == "False" {
		return nil, fmt.Errorf("movie not found: %s", movieInfo.Error)
	}

	return &movieInfo, nil
}

func (s *OMDbService) SearchMovies(query string, page int) ([]MovieInfo, error) {
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"s":      query,
			"page":   fmt.Sprintf("%d", page),
			"apikey": s.apiKey,
		}).
		Get(s.baseURL)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var result struct {
		Search   []MovieInfo `json:"Search"`
		Response string      `json:"Response"`
		Error    string      `json:"Error"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Response == "False" {
		return nil, fmt.Errorf("search failed: %s", result.Error)
	}

	return result.Search, nil
}
