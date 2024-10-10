Twitter API Integration using Go
================================

This Go application allows interaction with Twitter's API (v2) using OAuth 1.0a authentication. It enables sending and deleting tweets via a local HTTP server. Environment variables are loaded from a .env file for the necessary authentication credentials.

Features
--------

*   **Send Tweet**: POST request to post a tweet.
    
*   **Delete Tweet**: DELETE request to delete a tweet by ID.
    

Prerequisites
-------------

1.  **Go**: Install Go on your system by following the instructions from the official Go documentation.
    
2.  **Twitter Developer Account**: You need to have a Twitter developer account and a project/app with API keys and tokens.
    
Getting Started
---------------

### 1. Clone the Repository

Run the following commands to clone the repository:

```
git clone https://github.com/rsauri/w6_go_twitter
cd w6_go_twitter
```

### 2. Create a .env File

In the root of the project, create a .env file with your Twitter API credentials:

```
TWITTER_API_KEY=your-twitter-api-key
TWITTER_API_SECRET_KEY=your-twitter-api-secret-key
TWITTER_ACCESS_TOKEN=your-access-token
TWITTER_ACCESS_TOKEN_SECRET=your-access-token-secret
```

### 3. Run the Application

Run the following command:

```
go run main.go
```

The server will start running on http://localhost:8080.

API Endpoints
-------------

### 1. **POST /tweet**

This endpoint is used to send a tweet.

```
{
  "message": "Hello, World!"
}
```
    
*   **Response**:Returns the Twitter API response as JSON.
    

#### Example cURL Request:

```
curl -X POST http://localhost:8080/tweet -d '{"message": "Hello from Go!"}' -H "Content-Type: application/json"
```

### 2. **DELETE /tweet/{id}**

This endpoint deletes a tweet by ID.

*   **Request**:Replace {id} with the Tweet ID you want to delete.
    
*   **Response**:Returns the Twitter API response as JSON.
    

#### Example cURL Request:

```
curl -X DELETE http://localhost:8080/tweet/your-tweet-id
```


Error Handling
--------------

*   The application checks for required API credentials from the .env file.
    
*   If the credentials are missing or invalid, appropriate error messages will be returned.
    
*   HTTP error codes (400, 405) are returned for invalid requests.
    