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

//
func giphyLogin() *gophy.Client {
	giphyOptions := &gophy.ClientOptions{
		ApiKey: giphyApiKey}
	client := gophy.NewClient(giphyOptions)
	return client
}

func gifString(client *gophy.Client) string {
	// Get gif object from giphy
	randOffset := rand.Intn(500)
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

//
func twitterLogin() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(twitterConsumerKey)
	anaconda.SetConsumerSecret(twitterConsumerSecret)
	api := anaconda.NewTwitterApi(twitterAccessToken, twitterAccessSecret)
	return api
}

//
func sendTweet(client *anaconda.TwitterApi, text string, mediaID string) {
	v := url.Values{}
	v.Set("media_ids", mediaID)
	_, err := client.PostTweet(text, v)
	if err != nil {
		fmt.Println(err)
	}
}

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

//
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
		text = html.UnescapeString(comment + "&#13;&#10;&#13;&#10; /u/" + author)
		for len(text) > 140 && numGuesses < 10 {
			numGuesses++
			splitComment := strings.Split(comment, ".")
			comment = strings.Join(splitComment[:len(splitComment)-1], ".") + "."
			text = html.UnescapeString(comment + "&#13;&#10;&#13;&#10; /u/" + author)
		}
		if len(text) <= 140 && len(comment) > 1 {
			break
		}
	}
	return text
}

//
func makePost() {
	giphyClient := giphyLogin()
	twitterClient := twitterLogin()
	redditClient := redditLogin()

	text := redPillComment(redditClient)
	img := gifString(giphyClient)
	media, _ := twitterClient.UploadMedia(img)

	sendTweet(twitterClient, text, media.MediaIDString)
}

//
func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	for {
		makePost()
		time.Sleep(20 * time.Minute)
	}
}
