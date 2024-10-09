package sample

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type RankingData struct {
	name  string
	vote  string
	image string
}

func getData() []RankingData {
	targetPage := "https://gakumado.mynavi.jp/contests/mascot"
	resp, err := http.Get(targetPage)
	if err != nil {
		fmt.Println("Error")
	}
	defer resp.Body.Close()
	if resp.StatusCode > 200 {
		fmt.Printf("Status Code: %d", resp.StatusCode)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("failed to load")
		return nil
	}
	var result []RankingData
	doc.Find(".voteRanking-list").Children().Slice(0, 5).Each(func(i int, selection *goquery.Selection) {
		var newData RankingData
		newData.name = selection.Find(".voteRanking-list_name").Contents().First().Text()
		newData.name = strings.Replace(newData.name, "\n", "", -1)
		newData.name = strings.Replace(newData.name, " ", "", -1)
		newData.vote = selection.Find(".voteRanking-list_vote").Text()
		newData.image = selection.Find(".voteRanking-list_img img").AttrOr("src", "")
		result = append(result, newData)
	})
	return result
}
