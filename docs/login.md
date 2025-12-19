# Authentication

Clubs provides multiple secure authentication methods to access the platform. Choose the method that works best for you.

## First Time? Start Here

!!! tip "Brand New User?"
    If you've never used Clubs before, check out our [Complete Getting Started Guide](getting-started.md) for a full walkthrough!

## Accessing the Application

Before you can create an account, you need to know where to access Clubs:

1. **Get the URL**: Contact your organization or the person who invited you for your Clubs URL
   - It typically looks like: `https://clubs.yourorganization.com`
   - Or for local development: `http://localhost:3000` (see [Local Development Guide](../Documentation/LocalDev.md))

2. **Open in your browser**: Use a modern web browser (Chrome, Firefox, Safari, or Edge)

3. **Bookmark the page**: Save the URL for easy access later

## Creating Your Account

### Method 1: Single Sign-On (SSO) Registration

Clubs uses Keycloak for secure Single Sign-On authentication. This is the most common registration method.

**Step-by-step:**

1. **Navigate to the Login Page**: Open the Clubs application in your web browser
2. **Click "Sign Up"** or **"Create Account"** or **"Register"**
   - The exact text may vary, but look for an option to create a new account
   - You'll be redirected to the registration page
3. **Complete Registration Form**:
   - **Email address**: Use a valid email you can access
   - **Password**: Create a strong password that includes:
     - At least 8 characters
     - Uppercase letters (A-Z)
     - Lowercase letters (a-z)
     - Numbers (0-9)
     - Special characters (!@#$%^&*)
   - **First name**: Your first name
   - **Last name**: Your last name
4. **Submit the form**: Click "Register" or "Sign Up"
5. **Verify Your Email**: 
   - Check your email inbox for a verification link
   - **Check your spam folder** if you don't see it within 5 minutes
   - Click the verification link to activate your account
6. **Return to Clubs**: Go back to the login page
7. **Login**: Sign in with your email and password

!!! warning "Email Verification Required"
    You **must** verify your email address before you can use Clubs. If you don't receive the email:
    - Check your spam/junk folder
    - Verify you entered the correct email address
    - Request a new verification email from the login page

### Method 2: Organization SSO

Some organizations integrate Clubs with their existing authentication system (like Google Workspace, Microsoft Azure AD, etc.).

**If your organization uses this method:**

1. **Navigate to the Login Page**: Open the Clubs application
2. **Look for your organization's login button**: 
   - Examples: "Login with Google", "Login with Microsoft", "Login with [Your Organization]"
3. **Click the button**: You'll be redirected to your organization's login page
4. **Use your existing credentials**: Sign in with your work/organization account
5. **Grant permissions**: Allow Clubs to access your basic profile information
6. **Automatic account creation**: Your Clubs account will be created automatically

### Method 3: Magic Link Authentication

**⚠️ Important**: Magic Link login is for **existing users only**. You must first create an account using SSO before you can use Magic Links.

For passwordless login (after you have an account):

1. **Navigate to the Login Page**: Open the Clubs application
2. **Choose Magic Link**: 
   - Look for "Login with Magic Link" or "Passwordless Login"
   - Click this option
3. **Enter Your Email**: Provide the email address associated with your **existing** account
4. **Check Your Email**: You'll receive a secure login link within a few minutes
5. **Click the Link**: This will automatically log you into Clubs

!!! note "Magic Link Expiration"
    Magic Link emails are valid for 15-30 minutes for security purposes. If your link expires, simply request a new one.

!!! warning "First-Time Users Cannot Use Magic Link"
    If you don't have an account yet, you must first register using the SSO method above. Magic Links only work for existing accounts.

## Logging In (Returning Users)

## Logging In (Returning Users)

Once you have an account, here's how to log in:

### Using SSO (Username and Password)

1. **Open the Application**: Navigate to the Clubs platform URL
2. **Enter your credentials**:
   - **Email/Username**: The email address you registered with
   - **Password**: Your account password
3. **Click "Login"** or **"Sign In"**
4. **Access Your Dashboard**: After successful authentication, you'll be directed to your dashboard

### Using Magic Link (Passwordless)

1. **Open the Application**: Navigate to the Clubs platform
2. **Select "Login with Magic Link"** (or similar passwordless option)
3. **Enter your email address**: Use the email associated with your account
4. **Click "Send Magic Link"**
5. **Check your email**: A login link will arrive within a few minutes
6. **Click the link**: You'll be automatically logged in

### Using Organization SSO

1. **Open the Application**: Navigate to the Clubs platform
2. **Click your organization's login button**: (e.g., "Login with Google")
3. **Sign in**: Use your organization credentials
4. **Access Your Dashboard**: You'll be logged in automatically

## First Time Login

After your first successful login:

### 1. You'll See Your Dashboard

The dashboard is your home base in Clubs. It shows:
- **Upcoming Events**: Events you're invited to or have RSVP'd to
- **Recent News**: Latest posts from your clubs
- **My Clubs**: Clubs you're a member of
- **Notifications**: Important updates
- **Quick Actions**: Shortcuts to common tasks

### 2. Complete Your Profile (Recommended)

Take a few minutes to set up your profile:

1. **Click your profile icon** in the top navigation (usually top-right corner)
2. **Select "My Profile"** or **"Settings"**
3. **Click "Edit Profile"**
4. **Add information**:
   - Verify your name
   - Add a phone number (optional)
   - Add a profile picture (optional)
   - Write a short bio (optional)
5. **Click "Save"**

👉 **See the [Profile Management Guide](profile.md) for detailed instructions**

### 3. Join or Create a Club

You're now ready to join the club community!

**To Join an Existing Club:**
1. Click **"Clubs"** in the navigation menu
2. Browse available clubs
3. Click on a club to view details
4. Click **"Join Club"** or **"Request to Join"**
5. Wait for administrator approval

**To Create Your Own Club:**
1. Click **"Clubs"** in the navigation menu
2. Click **"Create Club"** or **"New Club"**
3. Fill in club details (name, description, privacy settings)
4. Click **"Save"** or **"Create"**

👉 **See the [Getting Started Guide](getting-started.md) for detailed step-by-step instructions**

### 4. Configure Notification Preferences (Optional)

Set how you want to receive updates:

1. **Go to Profile** → **Settings** → **Notifications**
2. **Choose your preferences**:
   - Email notifications (all, important only, daily digest, weekly, none)
   - Event notifications
   - Club notifications
   - Fine notifications
   - Shift notifications
3. **Save your preferences**

👉 **See the [Profile Management Guide](profile.md) for all notification options**

## Account Security

Clubs implements robust security measures to protect your account:

- **Secure Authentication**: All login sessions use encrypted tokens
- **Role-Based Access**: Different permission levels ensure appropriate access control
- **Session Management**: Monitor and manage your active sessions from your profile settings
- **Automatic Timeout**: Sessions expire after periods of inactivity for your protection

## Troubleshooting

### Can't Access Your Account?

- **Forgot Password**: Click "Forgot Password" on the login page to reset your credentials
- **Magic Link Not Received**: Check your spam/junk folder, or request a new link
- **Account Locked**: Contact your system administrator if you're unable to access your account

### Need Help?

If you encounter issues with authentication, please contact your organization's Clubs administrator or submit a support request through the help portal.
