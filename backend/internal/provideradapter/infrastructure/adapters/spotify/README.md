# Spotify Adapter Setup Guide

## Service Overview
The Spotify adapter allows your application to connect to the Spotify Web API through OAuth 2.0, accessing user music data, playlists, playback controls, and search functionality.

## Prerequisites
- Valid Spotify account (free or premium)
- Spotify developer account
- Basic understanding of OAuth 2.0 flow
- Active Spotify app (mobile/desktop) for playback control

## Important Note: Premium Requirements
Some Spotify API operations require the end user to have a Premium account:
- **Search, Get Playlists, Create/Modify Playlists**: Work with free accounts
- **Playback Control (start_playback, pause, skip, etc.)**: Require Premium account
- When Premium is required but not available, API returns: `403 PREMIUM_REQUIRED`

## Application Creation Steps

### Step 1: Create Spotify App
1. Visit [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/)
2. Log in to your Spotify account
3. Click the **Create app** button
4. Fill in the app information:
   - **App name**: Enter your app name
   - **App description**: Describe the app purpose
   - **Website**: Enter your website URL
   - **Redirect URI**: Add callback URL

### Step 2: Configure Redirect URI
Add the following redirect URI:
```
https://your-domain.com/auth/spotify/callback
```

### Step 3: Accept Terms
Check the box to agree to Spotify's Developer Terms of Service, then click **Save**

### Step 4: Permission Scopes

**User Information**
- `user-read-private` - Access user's basic information
- `user-read-email` - Access user email address

**Playlists**
- `playlist-read-private` - Read private playlists
- `playlist-read-collaborative` - Read collaborative playlists
- `playlist-modify-public` - Modify public playlists
- `playlist-modify-private` - Modify private playlists

**Playback Control**
- `user-read-playback-state` - Read playback state
- `user-modify-playback-state` - Control playback (required for start_playback)
- `user-read-currently-playing` - Read currently playing

**Music Library**
- `user-library-read` - Read music library
- `user-library-modify` - Modify music library

## Available Operations

### Search
- **Operation ID**: `search`
- **Description**: Search for tracks, artists, albums, playlists, shows, episodes, or audiobooks
- **Required Scopes**: None
- **Parameters**:
  - `q` (required): Search query keywords
  - `type` (required): Types to search (e.g., "track", "artist", "album")
  - `limit`: Maximum results (1-50)
  - `market`: Country code (e.g., "US")

### Playback Control
- **Operation ID**: `start_playback`
- **Description**: Start or resume playback on user's active device
- **Required Scopes**: `user-modify-playback-state`
- **Parameters**:
  - `device_id`: Target device ID
  - `context_uri`: Spotify URI of context (album/playlist)
  - `uris`: Array of track URIs to play
  - `position_ms`: Start position in milliseconds

## Demo Use Case: AI-Powered Mood Music

This adapter enables creating an AI assistant that:
1. Reads your daily notes from Notion
2. Checks current weather conditions
3. Generates a mood summary based on your notes and weather
4. Searches Spotify for songs matching your mood
5. Automatically plays music that matches your emotional state

Example workflow:
```
User notes: "Completed big project, feeling accomplished"
Weather: "Sunny, 75°F"
AI summary: "Energetic and celebratory mood"
Spotify search: "upbeat celebration victory"
Result: Plays motivational, happy songs
```

