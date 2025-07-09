# Notion Adapter Setup Guide

Reference: https://developers.notion.com/docs/create-a-notion-integration

## Service Overview
The Notion adapter allows your application to connect to the Notion API through OAuth 2.0, enabling programmatic interaction with your Notion workspace. You can create new databases, add pages to databases, add content to pages, and add comments to page content.

## Prerequisites
- A valid Notion account
- "Workspace Owner" permissions for the workspace you're using
- Knowledge of HTML and JavaScript
- npm and Node.js installed locally

## Application Creation Steps

### Step 1: Create Your Integration in Notion
1. Visit [Notion's Integration Dashboard](https://www.notion.com/my-integrations)
2. Click **+ New integration**
3. Enter the integration name and select the associated workspace for the new integration

### Step 2: Obtain Your API Key
1. Visit the **Configuration** tab to get your integration's API key (i.e., "Internal Integration Secret")
2. **Important**: Keep your API key secure and do not expose it to others or commit it to version control systems

### Step 3: Grant Page Permissions for Your Integration
1. In your Notion workspace, select or create a page to serve as the parent page for the integration
2. Click the **...** (more) menu in the top right corner of the page
3. Scroll down and click **+ Add Connections**
4. Search for the integration name you just created and select it
5. Confirm the authorization, allowing the integration to access this page and all its sub-pages

## OAuth Permission Scopes
According to the `manifest.json` file, the `oauth_scopes` array in the `permissions` section of this integration is empty, so no specific OAuth scopes need to be requested from users during the authorization flow.
