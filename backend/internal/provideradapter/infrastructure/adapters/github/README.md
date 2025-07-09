Doc: https://docs.github.com/en/apps/creating-github-apps/about-creating-github-apps/about-creating-github-apps

# GitHub Adapter Setup Guide

## Service Overview
The GitHub adapter allows your application to connect to the GitHub API through OAuth 2.0, enabling it to perform operations on behalf of users, such as reading user information, accessing repositories, and managing organization information.

## Prerequisites
- A valid GitHub account
- Logged in to your GitHub account
- Basic understanding of OAuth 2.0 authorization flow

## Application Creation Steps

### Step 1: Navigate to Developer Settings
1. Log in to your GitHub account
2. Click on your profile picture in the top right corner, then select **Settings**
3. In the left sidebar, scroll down and click **Developer settings**

### Step 2: Register New OAuth Application
1. On the developer settings page, click **OAuth Apps**
2. Click the **New OAuth App** button
3. Fill in the application information:
   - **Application name**: Enter a descriptive name for your application
   - **Homepage URL**: The complete homepage URL of your application
   - **Application description**: Provide a brief description of your application
   - **Authorization callback URL**: This is the URL where users will be redirected after authorization. Enter the following URL:
     ```
     https://your-domain.com/auth/github/callback
     ```
4. Click **Register application**

### Step 3: Obtain Client ID and Client Secret
1. After registration, you will be redirected to the application's settings page
2. You will see the **Client ID**
3. Click **Generate a new client secret** to generate a client secret. Be sure to store this secret securely as it will only be displayed once

## User Authentication
After configuration is complete, direct users to the following URL for authentication:
```
GET https://github.com/login/oauth/authorize?client_id=:client_id&redirect_uri=:callback&scope=:scope&state=:state&response_type=code
```
**Note**: You can request multiple scopes by separating them with spaces (`%20`) in the `scope` parameter, for example `scope=user%20repo`.

### Application Requested Scopes
- **`user:read`** (`read_user`): Read user's profile data
- **`repo`** (`repo_access`): Full control of private repositories
- **`read:org`** (`read_org`): Read organization and team membership