package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const maxOccurrence = 21290

var rwMutex = &sync.RWMutex{}
var yearGroups = make(map[int][]string, 122)

func main() {
	start := time.Now()

	file, err := os.Create("data.csv")
	if err != nil {
		fmt.Println("error trying to open the file")
		panic(err)
	}

	_, err = file.WriteString("date,aircraft,reg,operator,fat,location,dmg\n")
	if err != nil {
		fmt.Println("error trying to write to the file")
		panic(err)
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 0
	t.MaxConnsPerHost = 0
	t.MaxIdleConnsPerHost = 100
	httpClient := &http.Client{
		Transport: t,
	}
	wg := &sync.WaitGroup{}
	url := "http://aviation-safety.net/wikibase/dblist.php?Year="

	for i := 1902; i < 2024; i++ {
		yearGroups[i] = make([]string, 0, maxOccurrence)
		wg.Add(1)
		go process(i, url, httpClient, wg)
	}

	wg.Wait()

	keys := make([]int, 0, 122)

	for k := range yearGroups {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		for _, v := range yearGroups[k] {
			_, err = file.WriteString(v + "\n")
			if err != nil {
				fmt.Println("error trying to write to the file")
				panic(err)
			}
		}
	}

	fmt.Println("time elapsed:", time.Since(start))
}

func process(year int, url string, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	yearStr := strconv.Itoa(year)

	body := fetch(client, url+yearStr)

	node, _ := html.Parse(body)

	tr := getTr(node)
	current := tr

	occurrences := strings.SplitN(current.Parent.Parent.Parent.PrevSibling.LastChild.FirstChild.Data, " ", 2)[0]
	occurrencesInt, _ := strconv.Atoi(occurrences)
	pages := (occurrencesInt / 100) + 1

	if occurrencesInt%100 == 0 {
		pages--
	}

	for current != nil {
		data := extractData(current.FirstChild.NextSibling)

		rwMutex.Lock()
		yearGroups[year] = append(yearGroups[year], data)
		rwMutex.Unlock()

		current = current.NextSibling
	}

	if pages > 1 {
		for i := 2; i <= pages; i++ {
			body := fetch(client, url+yearStr+"&page="+strconv.Itoa(i))

			node, _ := html.Parse(body)

			tr := getTr(node)
			current := tr

			for current != nil {
				data := extractData(current.FirstChild.NextSibling)

				rwMutex.Lock()
				yearGroups[year] = append(yearGroups[year], data)
				rwMutex.Unlock()

				current = current.NextSibling
			}
		}
	}
}
