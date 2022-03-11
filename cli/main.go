package main

import (
	"fmt"

	"github.com/nmasse-itix/evdb/ficheauto"
)

func main() {
	scrapper := ficheauto.NewScrapper("https://www.fiches-auto.fr")
	cars := scrapper.Scrape()
	for _, car := range cars {
		fmt.Println(car)
	}
	// scrapper := ademe.NewScrapper("https://carlabelling.ademe.fr")
	// cars := scrapper.Scrape()
	// for _, car := range cars {
	// 	fmt.Println(car)
	// }
}
