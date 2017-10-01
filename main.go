package main

import (
	"encoding/base64"
	"fmt"
	"html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
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
	randOffset := rand.Intn(200)
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
	numComments := len(comments)
	for {
		num := rand.Intn(numComments)
		comment = html.UnescapeString(comments[num].Body)
		if len(comment) <= 140 {
			break
		}
	}
	return comment
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
