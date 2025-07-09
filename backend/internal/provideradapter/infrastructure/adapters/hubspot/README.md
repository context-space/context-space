# HubSpot Adapter Setup Guide

## Service Overview
The HubSpot adapter allows your application to connect to the HubSpot API through OAuth 2.0, managing contacts, companies, deals, marketing campaigns, and sales processes.

## Prerequisites
- Valid HubSpot account
- HubSpot developer account
- Basic understanding of OAuth 2.0 flow

## Application Creation Steps

### Step 1: Create HubSpot Developer Account
1. Visit [HubSpot Developer Platform](https://developers.hubspot.com/)
2. Click **Get started** or **Create Developer Account**
3. Log in with your HubSpot account or create a new account
4. Complete the developer account setup

### Step 2: Create HubSpot App
1. In the developer dashboard, click the **Create app** button
2. Select the app type:
   - **Public app**: Apps for multiple HubSpot accounts
   - **Private app**: Apps for use with a single HubSpot account only

### Step 3: Configure Basic Information
In the **Basic Info** tab, fill in:
- **App name**: Enter your app name
- **App description**: Describe the app functionality
- **App logo**: Optional, upload app icon

### Step 4: Configure OAuth Settings
1. Switch to the **Auth** tab
2. In the **Redirect URLs** section, add the callback URL:
   ```
   https://your-domain.com/auth/hubspot/callback
   ```
3. Select the required permission scopes

crm.objects.contacts.read
crm.objects.contacts.write
crm.objects.companies.read
crm.objects.companies.write
crm.objects.deals.read
crm.objects.deals.write
