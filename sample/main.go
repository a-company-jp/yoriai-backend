package sample

import (
	"encoding/json"
	"fmt"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type LineInstance struct {
	Token  string `json:"LineMessagingToken"`
	Name   string `json:"name"`
	Secret string `json:"LineSecret"`
}

func loadConfig(file string) []LineInstance {
	configFile, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error loading config")
	}
	var instances []LineInstance
	err = json.Unmarshal(configFile, &instances)
	if err != nil {
		fmt.Println("Error on unmarshal")
	}
	return instances
}

func createRankingComponent(data RankingData, color string) linebot.FlexComponent {
	return &linebot.BoxComponent{
		Layout: linebot.FlexBoxLayoutTypeHorizontal,
		Contents: []linebot.FlexComponent{
			&linebot.BoxComponent{Layout: linebot.FlexBoxLayoutTypeVertical,
				Contents: []linebot.FlexComponent{
					&linebot.TextComponent{
						Text:    data.vote,
						Color:   "#ffffff",
						Align:   linebot.FlexComponentAlignTypeCenter,
						Gravity: linebot.FlexComponentGravityTypeCenter,
					},
				},
				CornerRadius:    linebot.FlexComponentCornerRadiusTypeMd,
				BackgroundColor: color,
				JustifyContent:  linebot.FlexComponentJustifyContentTypeCenter,
			},
			&linebot.ImageComponent{
				URL:         "https://gakumado.mynavi.jp" + data.image,
				Size:        linebot.FlexImageSizeTypeXxs,
				Align:       linebot.FlexComponentAlignTypeCenter,
				AspectRatio: linebot.FlexImageAspectRatioType1to1,
				AspectMode:  linebot.FlexImageAspectModeTypeCover,
			},
			&linebot.TextComponent{Text: data.name,
				Gravity: linebot.FlexComponentGravityTypeCenter},
		},
	}
}

func createMessage(last int, image string, ranks []RankingData) *linebot.BubbleContainer {
	var rankingList []linebot.FlexComponent

	for i, r := range ranks {
		var color string
		switch i {
		case 0:
			color = "#e3c644"
			break
		case 1:
			color = "#a3a39d"
			break
		case 2:
			color = "#8f590d"
			break
		case 3:
			color = "#45a7ed"
			break
		case 4:
			color = "#38f2a8"
			break
		}
		rankingList = append(rankingList, createRankingComponent(r, color))
	}

	fmt.Println(last)
	if image == "" {
		image = "https://i.postimg.cc/fyD6NZDf/190-20220331210838.png"
	}
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{Text: "ただいまの状況", Size: linebot.FlexTextSizeTypeLg, Weight: linebot.FlexTextWeightTypeBold, Align: linebot.FlexComponentAlignTypeCenter},
				&linebot.BoxComponent{
					Layout:   linebot.FlexBoxLayoutTypeVertical,
					Contents: rankingList,
				},
				&linebot.ImageComponent{
					URL:         image,
					Size:        linebot.FlexImageSizeTypeFull,
					AspectRatio: "20:13",
					Margin:      linebot.FlexComponentMarginTypeSm,
					AspectMode:  linebot.FlexImageAspectModeTypeFit,
				},
				&linebot.TextComponent{
					Text:   "投票お願いもめ❤️",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeXl,
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
				&linebot.TextComponent{
					Text:   "残り" + strconv.Itoa(last) + "日",
					Size:   linebot.FlexTextSizeTypeLg,
					Align:  linebot.FlexComponentAlignTypeCenter,
					Weight: linebot.FlexTextWeightTypeRegular,
					Margin: linebot.FlexComponentMarginTypeXl,
				},
				&linebot.TextComponent{
					Text:  "11/07まで",
					Size:  linebot.FlexTextSizeTypeXxs,
					Align: linebot.FlexComponentAlignTypeCenter,
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Action: linebot.NewURIAction("モメポチ", "https://gakumado.mynavi.jp/contests/mascot/entries/124"),
					Style:  linebot.FlexButtonStyleTypePrimary,
				},
			},
		},
	}
}

func executeSend(configFile string, last int, rankingData []RankingData) {
	instances := loadConfig(configFile)
	for _, instance := range instances {
		client, err := linebot.New(instance.Secret, instance.Token)
		if err != nil {
			fmt.Print("Error1: ", instance.Name, err)
		}
		message := linebot.NewFlexMessage("投票してもめ～", createMessage(last, randomImage(), rankingData))
		if _, err := client.BroadcastMessage(message).Do(); err != nil {
			fmt.Println("Error2: ", instance.Name, err)
		}
	}
}

func randomImage() string {
	images := []string{
		"https://i.postimg.cc/fyD6NZDf/1.png",
		"https://i.postimg.cc/59pWdFcD/131-20220225000909.png",
		"https://i.postimg.cc/W3W8SVKY/IMG-0892.png",
		"https://i.postimg.cc/Hnr2hK1Y/IMG-0894.png",
		"https://i.postimg.cc/Xqt8S8nw/IMG-0896.png",
		"https://i.postimg.cc/2SbsZcqf/IMG-0897.png",
		"https://i.postimg.cc/wMjPT0KZ/IMG-0942.png",
		"https://i.postimg.cc/d1BpYgfG/IMG-0946.png",
		"https://i.postimg.cc/7ZFpM37J/126-20220227230130.png",
		"https://i.postimg.cc/zfT4dySt/128-20220224181530.png",
		"https://i.postimg.cc/Bb7z8H46/133-20220225182512.png",
		"https://i.postimg.cc/RVq8ksLn/135-20220302165519.png",
		"https://i.postimg.cc/FsqqNyTC/140-20220227224431.png",
		"https://i.postimg.cc/NjHnMKVs/142-20220329140827.png",
		"https://i.postimg.cc/Jn1YPPh0/143-20220301235859.png",
		"https://i.postimg.cc/zfYcRJZn/145-20220331214243.png",
		"https://i.postimg.cc/yNhrmws2/146-20220327204809.png",
		"https://i.postimg.cc/rmQnbxdJ/147-20220327210927.png",
		"https://i.postimg.cc/BbbYGbHG/148-20220302170534.png",
		"https://i.postimg.cc/fT3223Mc/149-20220302195435.png",
		"https://i.postimg.cc/Wzy9pVwq/150-20220302221403.png",
		"https://i.postimg.cc/3JPLPKwP/166-20220331225146.png",
		"https://i.postimg.cc/JhQPDZQL/167-20220329170258.png",
		"https://i.postimg.cc/SxCZ4w2j/175-20220329161436.png",
		"https://i.postimg.cc/RFVGchPn/180-20220331174313.png",
		"https://i.postimg.cc/N0tdjpmD/181-20220331180053.png",
		"https://i.postimg.cc/ncP8Mbhk/192-20220331221207.png",
	}
	rand.Seed(time.Now().Unix())
	return images[rand.Int()%len(images)]
}

func sendKamomeReminder(configFile string) {
	ranking := getData()
	const format = "2006-01-02 15:04:05 (MST)"
	limit, _ := time.Parse(format, "2022-11-07 23:00:00 (JST)")
	sub := limit.Sub(time.Now())
	remaining := int(sub.Hours()/24 + 1)
	if remaining > 0 {
		executeSend(configFile, remaining, ranking)
	}
}

func sendSimpleTextMessage(configFile string, msg string) {
	instances := loadConfig(configFile)
	for _, instance := range instances {
		client, err := linebot.New(instance.Secret, instance.Token)
		if err != nil {
			fmt.Print("Error1: ", instance.Name, err)
		}
		message := linebot.NewTextMessage(msg)
		if _, err := client.BroadcastMessage(message).Do(); err != nil {
			fmt.Println("Error2: ", instance.Name, err)
		}
	}
}

func main() {
	sendKamomeReminder(os.Args[1])
}
