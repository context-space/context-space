# Slack Adapter Setup Guide

## Service Overview
The Slack adapter allows your application to connect to the Slack API through OAuth 2.0, enabling integration with Slack workspaces, including sending messages, managing channels, and accessing user information.

## Prerequisites
- Slack account and workspace access permissions
- App installation permissions for the target workspace
- Basic understanding of OAuth 2.0 and Slack API

## Application Creation Steps

### Step 1: Create Slack App
1. Visit [Slack API Official Website](https://api.slack.com/)
2. Click **Your Apps** and log in to your Slack account
3. Click the **Create New App** button
4. Select **From scratch**
5. Enter app name and select development workspace
6. Click **Create App**

### Step 2: Configure OAuth Settings
1. In the app management panel, select **OAuth & Permissions**
2. In the **Redirect URLs** section, click **Add New Redirect URL**
3. Enter your callback URL:
   ```
   https://your-domain.com/auth/slack/callback
   ```
4. Click **Add** then **Save URLs**

### Step 3: Set Permission Scopes
In the **Scopes** section, add the following permissions based on your needs:

**Bot Token Scopes (Recommended)**:
- `channels:read` - View public channel information
- `chat:write` - Send messages
- `users:read` - View user information
- `files:read` - Access file information

**User Token Scopes (Optional)**:
- `identity.basic` - View user identity information
- `identity.email` - Access user email

