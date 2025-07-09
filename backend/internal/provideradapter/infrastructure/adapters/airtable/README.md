# Airtable Adapter Setup Guide

### OAuth Application Creation Steps

#### Step 1: Create OAuth Application
1. Visit [Airtable OAuth Application Creation Page](https://airtable.com/create/oauth)
2. Log in to your Airtable account
3. Click the **Create new OAuth integration** button
4. Fill in the application information:
   - **Integration name**: Enter your application name
   - **Integration logo**: Optional, upload application icon

#### Step 2: Configure Redirect URL
1. In the **Redirect URLs** section, add your callback URL:
   ```
   https://your-domain.com/auth/airtable/callback
   ```

3. Select the required permission scopes

data.records:read

data.records:write

schema.bases:read

schema.bases:write

user.email:read

webhook:manage