package tmdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type Movie struct {
	Adult        bool     `json:"adult"`
	BackdropPath string   `json:"backdrop_path"`
	GenreIds     []int    `json:"genre_ids"`
	ID           int      `json:"id"`
	OriginalLang string   `json:"original_language"`
	OriginalTitle string  `json:"original_title"`
	Overview     string   `json:"overview"`
	Popularity   float64  `json:"popularity"`
	PosterPath   string   `json:"poster_path"`
	ReleaseDate  string   `json:"release_date"`
	Title        string   `json:"title"`
	Video        bool     `json:"video"`
	VoteAverage  float64  `json:"vote_average"`
	VoteCount    int      `json:"vote_count"`
}

type ApiResponse struct {
	Dates struct {
		Maximum string `json:"maximum"`
		Minimum string `json:"minimum"`
	} `json:"dates"`
	Page    int     `json:"page"`
	Results []Movie `json:"results"`
}

var TMDBUrl = map[string]string{
	"validate": "https://api.themoviedb.org/3/authentication?api_key=%s",
	"popular": "https://api.themoviedb.org/3/movie/popular?language=en-US&page=1&api_key=%s",
	"top": "https://api.themoviedb.org/3/movie/top_rated?language=en-US&page=1&api_key=%s",
	"upcoming": "https://api.themoviedb.org/3/movie/upcoming?language=en-US&page=1&api_key=%s",
	"playing": "https://api.themoviedb.org/3/movie/now_playing?language=en-US&page=1&api_key=%s",
}


var Genres = map[int]string{
	28:    "Action",
	12:    "Adventure",
	16:    "Animation",
	35:    "Comedy",
	80:    "Crime",
	99:    "Documentary",
	18:    "Drama",
	10751: "Family",
	14:    "Fantasy",
	36:    "History",
	27:    "Horror",
	10402: "Music",
	9648:  "Mystery",
	10749: "Romance",
	878:   "Science Fiction",
	10770: "TV Movie",
	53:    "Thriller",
	10752: "War",
	37:    "Western",
}


var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func ValidateKey(key string) error {
	if key == "" {
		return errors.New("API key cannot be empty")
	}

	url, exists := TMDBUrl["validate"]
	if !exists {
		return errors.New("validation endpoint not found")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(url, key), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid API key, status code: %d", res.StatusCode)
	}

	return nil
}


func GetMovies(Mtype string, key string) ([]Movie, error) {
	url, exists := TMDBUrl[Mtype]
	if !exists {
		return nil, errors.New("invalid type (options are: latest, popular, top, upcoming, playing)")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(url, key), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid API key, status code: %d", res.StatusCode)
	}

	var movieRes ApiResponse
	err = json.NewDecoder(res.Body).Decode(&movieRes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}


	return movieRes.Results, nil
}


func DisplayMovies(movies []Movie) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Title", "Year", "Genres", "Popularity", "Vote Average"})
	table.SetBorder(true)

	for _, movie := range movies {
		genres := getGenres(movie.GenreIds, ",")
		table.Append([]string{movie.Title, movie.ReleaseDate, genres, fmt.Sprintf("%.2f", movie.Popularity), fmt.Sprintf("%.2f", movie.VoteAverage)})
	}


	table.Render()
}


func getGenres(genres []int, separator string) string {
	if len(genres) == 0 {
		return "N/A"
	}

	if(separator == "") {
		separator = ","
	}

	genresString := make([]string, len(genres))
	for i, genre := range genres {
		genresString[i] = Genres[genre]
	}

	if len(genresString) == 0 {
		return "N/A"
	}

	if len(genresString) == 1 {
		return genresString[0]
	}

	return strings.Join(genresString, fmt.Sprintf("%s ", separator))
}
