package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
)

func fetch(client *http.Client, url string) io.ReadCloser {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("error trying to prepare the request")
		panic(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0")
	response, err := client.Do(req)

	if err != nil {
		fmt.Println("fetch: error trying to get the page")
		panic(err)
	}

	return response.Body
}

func getTr(node *html.Node) *html.Node {
	return node.
		FirstChild.FirstChild.NextSibling.LastChild.
		PrevSibling.PrevSibling.PrevSibling.PrevSibling.
		PrevSibling.PrevSibling.PrevSibling.PrevSibling.
		PrevSibling.PrevSibling.PrevSibling.PrevSibling.
		PrevSibling.PrevSibling.FirstChild.NextSibling.
		FirstChild.NextSibling.NextSibling.NextSibling.
		NextSibling.NextSibling.NextSibling.NextSibling.
		NextSibling.NextSibling.NextSibling.NextSibling.
		NextSibling.NextSibling.NextSibling.NextSibling.
		NextSibling.NextSibling.FirstChild.NextSibling.
		FirstChild.FirstChild.NextSibling.NextSibling.
		NextSibling.NextSibling.NextSibling.FirstChild.
		FirstChild.NextSibling.FirstChild.NextSibling
}

func getDate(tr *html.Node) string {
	dateElem := tr.FirstChild.FirstChild.FirstChild

	if dateElem == nil {
		return tr.FirstChild.FirstChild.NextSibling.NextSibling.FirstChild.Data
	}

	return dateElem.Data
}

func getAircraft(tr *html.Node) string {
	aircraftElem := tr.NextSibling.NextSibling.FirstChild

	if aircraftElem == nil {
		return ""
	}

	return aircraftElem.Data
}

func getReg(tr *html.Node) string {
	regElem := tr.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild

	if regElem == nil {
		return ""
	}

	return regElem.Data
}

func getOperator(tr *html.Node) string {
	operatorElem := tr.NextSibling.NextSibling.NextSibling.NextSibling.
		NextSibling.NextSibling.FirstChild

	if operatorElem == nil {
		return ""
	}

	return operatorElem.Data
}

func getFat(tr *html.Node) string {
	fatElem := tr.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.
		NextSibling.NextSibling.NextSibling.FirstChild

	if fatElem == nil {
		return ""
	}

	return fatElem.Data
}

func getLocation(tr *html.Node) string {
	locationElem := tr.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.
		NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild

	if locationElem == nil {
		return ""
	}

	return locationElem.Data
}

func getDmg(tr *html.Node) string {
	dmgElem := tr.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.
		NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild

	if dmgElem == nil {
		return ""
	}

	return dmgElem.Data
}

func extractData(row *html.Node) string {
	date := getDate(row)
	aircraft := getAircraft(row)
	reg := getReg(row)
	operator := getOperator(row)
	fat := getFat(row)
	location := getLocation(row)
	dmg := getDmg(row)

	return date + "," + aircraft + "," + reg + "," + operator + "," + fat + "," + location + "," + dmg
}
