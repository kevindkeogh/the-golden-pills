package main

import (
	"encoding/base64"
	"fmt"
	"html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/jzelinskie/geddit"
	"github.com/paddycarey/gophy"
)

// giphyLogin returns the Giphy connected API Client
func giphyLogin() *gophy.Client {
	giphyOptions := &gophy.ClientOptions{
		ApiKey: giphyApiKey}
	client := gophy.NewClient(giphyOptions)
	return client
}

// gifString takes the Giphy API Client and returns a base 64 encoded string
// of a random Golden Girls gif. Note that the function searches the term
// "golden girls" and returns a random result from the top 400 results.
func gifString(client *gophy.Client) string {
	// Get gif object from giphy
	randOffset := rand.Intn(400)
	gifs, _, err := client.SearchGifs("golden girls", "", 1, randOffset)
	if err != nil {
		fmt.Println(err)
	}

	// Download the image from the URL
	resp, err := http.Get(gifs[0].Images.Original.URL)
	if err != nil {
		fmt.Println(err)
	}

	// Convert to base64
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded
}

// twitterLogin returns the Twitter connected API Client
func twitterLogin() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(twitterConsumerKey)
	anaconda.SetConsumerSecret(twitterConsumerSecret)
	api := anaconda.NewTwitterApi(twitterAccessToken, twitterAccessSecret)
	return api
}

// sendTweet takes the Twitter API client, the text of the tweet to be sent,
// and the ID of the previously uploaded GIF..
func sendTweet(client *anaconda.TwitterApi, text string, mediaID string) {
	v := url.Values{}
	v.Set("media_ids", mediaID)
	_, err := client.PostTweet(text, v)
	if err != nil {
		fmt.Println(err)
	}
}

// redditLogin returns the Reddit connected API Client
func redditLogin() *geddit.OAuthSession {
	client, err := geddit.NewOAuthSession(
		redditID,
		redditSecret,
		"@the-golden-pills",
		"http://reddit.com",
	)
	if err != nil {
		fmt.Println(err)
	}
	err = client.LoginAuth(redditUsername, redditPassword)
	if err != nil {
		fmt.Println(err)
	}
	return client
}

// redPillComment takes a Reddit API Client and returns the string of a random
// recent comment. Note that the string is stripped to ensure that it is <140
// characters. It is also HTML unescaped, and included the /u/{username}
// signature of the commenter
func redPillComment(client *geddit.OAuthSession) string {
	comments, err := client.SubredditComments("TheRedPill")
	if err != nil {
		fmt.Println(err)
	}

	var comment string
	var author string
	var text string
	var numGuesses int
	numComments := len(comments)
	for {
		numGuesses = 0
		num := rand.Intn(numComments)
		comment = comments[num].Body
		author = comments[num].Author
		text = html.UnescapeString(comment + "&#13;&#10;&#13;&#10;/u/" + author)
		for len(text) > 140 && numGuesses < 10 {
			numGuesses++
			splitComment := strings.Split(comment, ".")
			comment = strings.Join(splitComment[:len(splitComment)-1], ".") + "."
			text = html.UnescapeString(comment + "&#13;&#10;&#13;&#10;/u/" + author)
		}
		if len(text) <= 140 && len(comment) > 1 {
			break
		}
	}
	return text
}

// makePost creates new API clients for Giphy, Twitter, and Reddit, and creates
// a single tweet
func makePost() {
	giphyClient := giphyLogin()
	twitterClient := twitterLogin()
	redditClient := redditLogin()

	text := redPillComment(redditClient)
	img := gifString(giphyClient)
	media, _ := twitterClient.UploadMedia(img)

	sendTweet(twitterClient, text, media.MediaIDString)
}

// main runs forever and executes the makePost function every 30 minutes
func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	for {
		makePost()
		time.Sleep(30 * time.Minute)
	}
}
