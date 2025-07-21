# TMDB (The Movie Database) Adapter

This adapter provides integration with The Movie Database (TMDB) API, allowing access to millions of movies, TV shows, cast, crew, and genre information.

## Features

- **Genre Information**: Get official movie and TV show genres
- **Movie Data**: Access popular, now playing, top rated, and upcoming movies
- **TV Show Data**: Retrieve popular, on-air, top rated, and airing today TV shows
- **Content Discovery**: Advanced discovery features for movies and TV shows with filtering
- **Trending Content**: Get trending movies, TV shows, and people (daily/weekly)
- **Search Functionality**: Search across movies, TV shows, people, or multiple content types
- **Detailed Information**: Get comprehensive details about specific movies and TV shows

## Supported Operations

| Operation ID          | Name                  | Description                                | Category |
| --------------------- | --------------------- | ------------------------------------------ | -------- |
| `genre_movie_list`  | Get Movie Genres      | Get the list of official movie genres      | Movies   |
| `genre_tv_list`     | Get TV Genres         | Get the list of official TV genres         | TV       |
| `movie_popular`     | Popular Movies        | Get a list of movies ordered by popularity | Movies   |
| `movie_now_playing` | Now Playing Movies    | Get movies currently in theaters           | Movies   |
| `movie_top_rated`   | Top Rated Movies      | Get the top rated movies on TMDB           | Movies   |
| `movie_upcoming`    | Upcoming Movies       | Get movies that are being released soon    | Movies   |
| `movie_details`     | Movie Details         | Get primary information about a movie      | Movies   |
| `tv_popular`        | Popular TV Shows      | Get TV shows ordered by popularity         | TV       |
| `tv_on_the_air`     | TV Shows On The Air   | Get TV shows that air in the next 7 days   | TV       |
| `tv_top_rated`      | Top Rated TV Shows    | Get the top rated TV shows on TMDB         | TV       |
| `tv_airing_today`   | TV Shows Airing Today | Get TV shows airing today                  | TV       |
| `tv_details`        | TV Show Details       | Get primary information about a TV show    | TV       |
| `discover_movie`    | Discover Movies       | Discover movies with advanced filtering    | Movies   |
| `discover_tv`       | Discover TV Shows     | Discover TV shows with advanced filtering  | TV       |
| `trending`          | Trending              | Get trending movies, TV shows and people   | Trending |
| `search_movie`      | Search Movies         | Search for movies by title                 | Search   |
| `search_tv`         | Search TV Shows       | Search for TV shows by name                | Search   |
| `search_person`     | Search People         | Search for people by name                  | Search   |
| `search_multi`      | Multi Search          | Search across multiple content types       | Search   |

## Authentication

This adapter uses API Key authentication. The API key is passed as a query parameter (`api_key`) with all requests.

### Required Configuration

- `api_key`: Your TMDB API key

### Getting an API Key

1. Visit [TMDB website](https://www.themoviedb.org/)
2. Create an account
3. Go to Settings > API
4. Generate an API key

## Usage Examples

### Get Popular Movies

```json
{
  "operation_id": "movie_popular",
  "parameters": {
    "language": "en-US",
    "page": 1,
    "region": "US"
  }
}
```

### Search for Movies

```json
{
  "operation_id": "search_movie",
  "parameters": {
    "query": "The Shawshank Redemption",
    "language": "en-US",
    "year": 1994
  }
}
```

### Get Movie Details

```json
{
  "operation_id": "movie_details",
  "parameters": {
    "movie_id": 278,
    "language": "en-US",
    "append_to_response": "credits,videos"
  }
}
```

### Discover Movies with Filters

```json
{
  "operation_id": "discover_movie",
  "parameters": {
    "with_genres": "28,12",
    "primary_release_year": 2024,
    "vote_average.gte": 7.0,
    "sort_by": "popularity.desc"
  }
}
```

### Get Trending Content

```json
{
  "operation_id": "trending",
  "parameters": {
    "media_type": "movie",
    "time_window": "week",
    "language": "en-US"
  }
}
```

## Common Parameters

- **language**: ISO 639-1 language code (default: en-US)
- **page**: Page number for pagination (1-1000, default: 1)
- **region**: ISO 3166-1 country code for region-specific data
- **include_adult**: Include adult content in search results

## API Endpoints

All requests are made to the TMDB API v3 base URL: `https://api.themoviedb.org/3`

## Error Handling

The adapter handles common TMDB API errors:

- **401 Unauthorized**: Invalid API key
- **404 Not Found**: Resource not found
- **422 Unprocessable Entity**: Invalid parameters
- **429 Too Many Requests**: Rate limit exceeded

## Rate Limiting

TMDB enforces rate limiting on their API. The adapter includes retry logic with exponential backoff to handle rate limit responses.

## Links

- [TMDB API Documentation](https://developer.themoviedb.org/docs/getting-started)
- [TMDB Website](https://www.themoviedb.org/)
- [TMDB API Reference](https://developer.themoviedb.org/reference/intro/getting-started)
