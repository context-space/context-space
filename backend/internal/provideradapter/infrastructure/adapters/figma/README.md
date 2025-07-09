# Figma Adapter Setup Guide

## Service Overview
Registering a Figma OAuth application allows you to establish a connection between existing Figma user accounts and your application, enabling you to access data on behalf of users through Figma's API.

## Prerequisites
- A valid Figma account

## Application Creation Steps

### Step 1: Access App Management Page
1. Click **My apps** in the Figma page top toolbar, or visit [figma.com/developers/apps](https://www.figma.com/developers/apps) directly.

### Step 2: Create New App
1. In the top right corner of the page, click the **Create a new app** button.

### Step 3: Configure Basic Information
Fill in the application creation page:
- **Name and Website**: Enter your app name and website.
- **Logo**: Upload a logo for your app (users will see this).
- Click **Save** to save the app.

### Step 4: Obtain Client ID and Client Secret
1. After successful creation, you will receive a **Client ID** and **Client Secret**.
2. **Note**: The Client Secret will only be displayed once, so copy and store it securely immediately.

### Step 5: Configure Redirect URL
1. In your app list, click on the app you just created.
2. In the modal that appears, click **OAuth 2.0**, then click **Add a redirect URL**.
3. Add the callback URL. This is an endpoint in your application for receiving authentication information.
   ```
   https://your-domain.com/auth/figma/callback
   ```

## User Authentication
After configuration is complete, direct users to the following URL for authentication:
```
GET https://www.figma.com/oauth?client_id=:client_id&redirect_uri=:callback&scope=:scope&state=:state&response_type=code
```
**Important**: Users must access this URL through a browser, not through an embedded webview within an app. Figma does not support OAuth flow in webviews.

## Required Permission Scopes
When requesting authorization, you can apply for the following permission scopes as needed:
- `files:read`
- `file_variables:read`
- `file_variables:write`
- `library_analytics:read`
