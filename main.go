package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

// Replace with your own credentials from the Twitter Developer Portal
var (
	apiKey            string
	apiSecretKey      string
	accessToken       string
	accessTokenSecret string
	baseURI           string = "https://api.twitter.com/2"
)

// Create new struct
type tweet struct {
	Message string `json:"message"`
}

// Load environment variables from .env file
func loadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Ignore empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in .env file: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// Set the environment variable
		os.Setenv(key, value)
	}

	// Fetch environment variables
	apiKey = os.Getenv("TWITTER_API_KEY")
	apiSecretKey = os.Getenv("TWITTER_API_SECRET_KEY")
	accessToken = os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	return scanner.Err()
}

// Generate OAuth signature
func generateSignature(method, baseURL string, params map[string]string, tokenSecret string) string {
	// Step 1: Sort the parameters
	var paramStr string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, key := range keys {
		paramStr += url.QueryEscape(key) + "=" + url.QueryEscape(params[key])
		if i < len(keys)-1 {
			paramStr += "&"
		}
	}

	// Step 2: Create the base signature string
	signatureBaseStr := method + "&" + url.QueryEscape(baseURL) + "&" + url.QueryEscape(paramStr)

	// Step 3: Create the signing key
	signingKey := url.QueryEscape(apiSecretKey) + "&" + url.QueryEscape(tokenSecret)

	// Step 4: Hash the signature base string with the signing key
	hash := hmac.New(sha1.New, []byte(signingKey))
	hash.Write([]byte(signatureBaseStr))
	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	return signature
}

// Generate OAuth 1.0a Header
func generateOAuthHeader(method, baseURL string) string {

	// Prepare OAuth parameters
	oauthParams := map[string]string{
		"oauth_consumer_key":     apiKey,
		"oauth_nonce":            fmt.Sprintf("%d", time.Now().UnixNano()),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_token":            accessToken,
		"oauth_version":          "1.0",
	}

	// Generate OAuth signature
	allParams := make(map[string]string)
	for k, v := range oauthParams {
		allParams[k] = v
	}
	oauthParams["oauth_signature"] = generateSignature(method, baseURL, allParams, accessTokenSecret)

	header := "OAuth "
	for k, v := range oauthParams {
		header += fmt.Sprintf(`%s="%s", `, url.QueryEscape(k), url.QueryEscape(v))
	}
	return strings.TrimSuffix(header, ", ")
}

func postTweet(w http.ResponseWriter, r *http.Request) {
	apiURL := baseURI + "/tweets"
	method := http.MethodPost

	//Get the message from the request
	var tweetContent tweet
	json.NewDecoder(r.Body).Decode(&tweetContent)

	// Prepare request body (tweet content)
	data := map[string]string{
		"text": tweetContent.Message,
	}
	jsonData, _ := json.Marshal(data)

	// Generate the OAuth header
	oauthHeader := generateOAuthHeader(method, apiURL)
	// Create a new HTTP request
	req, err := http.NewRequest(method, apiURL, strings.NewReader(string(jsonData)))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error found: %s", err), http.StatusNotAcceptable)
	}

	// Set the necessary headers
	req.Header.Set("Authorization", oauthHeader)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error found: %s", err), http.StatusNotAcceptable)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error found: %s", err), http.StatusNotAcceptable)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(body))
}

func deleteTweet(w http.ResponseWriter, r *http.Request) {

	//Get the ID from the Request URI
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid Path", http.StatusNotAcceptable)
	}
	id := parts[2]

	method := http.MethodDelete
	apiURL := baseURI + "/tweets/" + id

	// Get the OAuth header
	oauthHeader := generateOAuthHeader(method, apiURL)

	// Create a new HTTP request
	req, err := http.NewRequest(method, apiURL, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error found: %s", err), http.StatusNotAcceptable)
	}

	// Set the necessary headers
	req.Header.Set("Authorization", oauthHeader)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error found: %s", err), http.StatusNotAcceptable)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error found: %s", err), http.StatusNotAcceptable)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(body))
}

func main() {
	//Load .env variables
	loadEnv()

	// Set up HTTP server
	http.HandleFunc("/tweet", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postTweet(w, r)
		default:
			http.Error(w, fmt.Sprintf("Invalid request method %s", r.Method), http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tweet/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			deleteTweet(w, r)
		default:
			http.Error(w, fmt.Sprintf("Invalid request method %s", r.Method), http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
