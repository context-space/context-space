# Spotify Adapter Setup Guide

## Service Overview
The Spotify adapter allows your application to connect to the Spotify Web API through OAuth 2.0, accessing user music data, playlists, and playback controls.

## Prerequisites
- Valid Spotify account (free or premium)
- Spotify developer account
- Basic understanding of OAuth 2.0 flow

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
- `user-modify-playback-state` - Control playback
- `user-read-currently-playing` - Read currently playing

**Music Library**
- `user-library-read` - Read music library
- `user-library-modify` - Modify music library

