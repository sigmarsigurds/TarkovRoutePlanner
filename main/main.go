package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"strings"
)

/*
	type textbox struct {
		MapImage            string        `json:"mapImage"`
		PageCategories      []interface{} `json:"pageCategories"`
		DefaultSort         string        `json:"defaultSort"`
		Description         string        `json:"description"`
		CoordinateOrder     string        `json:"coordinateOrder"`
		MapBounds           [][]int       `json:"mapBounds"`
		Origin              string        `json:"origin"`
		UseMarkerClustering bool          `json:"useMarkerClustering"`
		Categories          []struct {
			Id          string `json:"id"`
			ListId      int    `json:"listId"`
			Name        string `json:"name"`
			Color       string `json:"color"`
			Symbol      string `json:"symbol"`
			SymbolColor string `json:"symbolColor"`
			Icon        string `json:"icon"`
		} `json:"categories"`
		Markers []struct {
		}
	}
*/

type categories struct {
	Id          string `json:"id"`
	ListId      int    `json:"listId"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Symbol      string `json:"symbol"`
	SymbolColor string `json:"symbolColor"`
	Icon        string `json:"icon"`
}

type marker struct {
	CategoryId string    `json:"categoryId"`
	Position   []float64 `json:"position"`
	Popup      struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Link        struct {
			Url   string `json:"url"`
			Label string `json:"label"`
		} `json:"link"`
	} `json:"popup"`
	Config struct {
		SortMarkers      string   `json:"sortMarkers"`
		HiddenCategories []string `json:"hiddenCategories"`
	} `json:"config,omitempty"`
	Id string `json:"id"`
}
type textbox struct {
	MapImage            string        `json:"mapImage"`
	PageCategories      []interface{} `json:"pageCategories"`
	DefaultSort         string        `json:"defaultSort"`
	Description         string        `json:"description"`
	CoordinateOrder     string        `json:"coordinateOrder"`
	MapBounds           [][]int       `json:"mapBounds"`
	Origin              string        `json:"origin"`
	UseMarkerClustering bool          `json:"useMarkerClustering"`
	Categories          []categories  `json:"categories"`
	Markers             []marker      `json:"markers"`
}

const customs = "Customs_Interactive_Map"
const woods = "Woods_Interactive_Map"

func findMapPng(response *http.Response) (string, error) {
	prefix := "https://static.wikia.nocookie.net/escapefromtarkov_gamepedia/images"
	mapUrl := ""
	error := false

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}

	document.Find("script").Each(func(i int, s *goquery.Selection) {
		var scriptContent string = s.Text()

		if strings.Contains(scriptContent, prefix) { // finn backround url staðinn
			imageIndex := strings.Index(scriptContent, prefix) // find start
			if imageIndex == -1 {
				error = true
			}
			imageEndIndex := strings.Index(scriptContent[imageIndex:imageIndex+500], "\"") // find end

			if imageEndIndex == -1 {
				error = true
			}
			mapUrl = scriptContent[imageIndex : imageIndex+imageEndIndex]
		}
	})
	if error == true {
		return "", errors.New("map not found.")
	} else {
		return mapUrl, nil
	}

}
func getBaseMap(mapName string) {
	baseString := "https://escapefromtarkov.fandom.com/wiki/Map:"
	response, err := http.Get(baseString + mapName)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	mapUrl, mapError := findMapPng(response)

	if mapError != nil {
		panic(mapError)
	} else {
		println(mapUrl)
	}
	if downloadImage(mapName, mapUrl) != nil {
		panic(err)
	}

}

func downloadImage(imageName string, imageURL string) error {
	// get the image
	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// create a file
	imageFile, err := os.Create(imageName + ".png")
	if err != nil {
		panic(err)
	}
	defer imageFile.Close()

	// copy image into file
	_, err = io.Copy(imageFile, response.Body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	fmt.Println("Image downloaded successfully!")
	return nil
}

func getMapData() textbox {
	baseString := "https://escapefromtarkov.fandom.com/wiki/Map:Customs_Interactive_Map?action=edit"

	resp, err := http.Get(baseString)
	if err != nil {
		println("MARKERS: could not find map url")
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	var text textbox
	strMarker := ""
	doc.Find("#wpTextbox1").Each(func(i int, s *goquery.Selection) {
		strMarker = s.Text()
	})
	json.Unmarshal([]byte(strMarker), &text)

	return text
}

func printMarkerTitles(mapData textbox) {
	fmt.Printf("%v\n", mapData)
	for _, item := range mapData.Markers {
		fmt.Printf("%v \n", item.Popup.Title)
	}
}

func getLocationByTitle(mapData textbox, title string) [][]float64 {
	var positions [][]float64
	for _, item := range mapData.Markers {
		if item.Popup.Title == title {
			positions = append(positions, item.Position)
		}
	}
	return positions
}

func main() {
	//crawler()
	//getBaseMap(woods)
	mapData := getMapData()
	//printMarkerTitles(mapData)
	locations := getLocationByTitle(mapData, "PMC Spawn")
	for _, location := range locations {
		fmt.Printf("%v\n", location)
	}
	drawLine("Customs_Interactive_Map.png")

	//fmt.Printf("Categories: %v, Description: %s", mapData.Categories, mapData.PageCategories)

}
