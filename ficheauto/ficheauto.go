package ficheauto

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Car struct {
	Brand            string  // Brand name
	Model            string  // Model + Variant
	Weight           int     // Weight (kg)
	Layout           string  // Motor on front, rear wheels or both
	Power            int     // Power (ch)
	Acceleration     float32 // From 0 to 100km/h (s)
	MaxSpeed         int     // Max speed (km/h)
	TrunkSpace       int     // Available space in the trunk (dm3)
	BatteryCapacity  float32 // Battery capacity (kWh)
	MaxChargingPower int     // Maximum charging power (kW)
	RetailPrice      int     // Retail price (k€)
}

func (car Car) String() string {
	return fmt.Sprintf("%s %s (%d ch, %.1f kWh)", car.Brand, car.Model, car.Power, car.BatteryCapacity)
}

type Scrapper struct {
	url string
	c   *colly.Collector
}

func NewScrapper(baseUrl string) *Scrapper {
	scrapper := Scrapper{
		url: fmt.Sprintf("%s/articles-auto/electrique/s-852-comparatif-des-voitures-electriques.php", baseUrl),
		c:   colly.NewCollector(),
	}
	return &scrapper
}

func (s *Scrapper) Scrape() []Car {
	var cars []Car

	s.c.OnHTML("table", func(table *colly.HTMLElement) {
		if table.Attr("bordercolor") != "#C0C0C0" {
			return
		}
		table.ForEach("tbody", func(_ int, tbody *colly.HTMLElement) {
			var cols []string = []string{}
			tbody.ForEach("tr", func(i int, tr *colly.HTMLElement) {
				var vals []string = []string{}
				tr.ForEach("td", func(j int, td *colly.HTMLElement) {
					val := strings.TrimSpace(td.Text)
					//fmt.Println(val)
					if i == 0 {
						cols = append(cols, val)
					} else {
						vals = append(vals, val)
					}
				})
				if i > 0 {
					data, err := slice2map(cols, vals)
					if err != nil {
						fmt.Println(err)
						return
					}
					car, err := map2car(data)
					if err != nil {
						fmt.Println(err)
						return
					}
					cars = append(cars, car)
				}
			})
		})
	})

	s.c.Visit(s.url)

	return cars
}

var nameRegexp, miscRegexp *regexp.Regexp

func init() {
	nameRegexp = regexp.MustCompile(`^([A-Za-z0-9]+) +(.*) +\(([0-9]+) ?kg\)`)
	miscRegexp = regexp.MustCompile(`^([.0-9]+)`)
}

func getInt(s string) (int, error) {
	matches := miscRegexp.FindStringSubmatch(s)
	if len(matches) != 2 {
		return 0, fmt.Errorf("wrong number of matches: '%s'", s)
	}
	i, err := strconv.ParseInt(matches[1], 10, 32)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func getFloat(s string) (float32, error) {
	matches := miscRegexp.FindStringSubmatch(s)
	if len(matches) != 2 {
		return 0, fmt.Errorf("wrong number of matches: '%s'", s)
	}
	f, err := strconv.ParseFloat(matches[1], 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func map2car(data map[string]string) (Car, error) {
	var car Car
	matches := nameRegexp.FindStringSubmatch(data["Modèles"])
	if len(matches) != 4 {
		return Car{}, fmt.Errorf("wrong number of matches: '%s'", data["Modèles"])
	}
	car.Brand = strings.TrimSpace(matches[1])
	car.Model = strings.TrimSpace(matches[2])
	w, err := strconv.ParseInt(matches[3], 10, 32)
	if err != nil {
		return Car{}, err
	}
	car.Weight = int(w)
	car.Layout = data["Motr."]
	car.Acceleration, err = getFloat(data["0/100sec."])
	if err != nil {
		fmt.Println(err)
	}
	car.BatteryCapacity, err = getFloat(data["Bat.kWh"])
	if err != nil {
		fmt.Println(err)
	}
	car.MaxChargingPower, err = getInt(data["Puiss.ChargeMAX"])
	if err != nil {
		fmt.Println(err)
	}
	car.MaxSpeed, err = getInt(data["Vmaxkm/h"])
	if err != nil {
		fmt.Println(err)
	}
	car.Power, err = getInt(data["Puiss.ch"])
	if err != nil {
		fmt.Println(err)
	}
	car.RetailPrice, err = getInt(data["Prix"])
	if err != nil {
		fmt.Println(err)
	}
	car.TrunkSpace, err = getInt(data["CoffreLitres"])
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(data)

	return car, nil
}

func slice2map(cols, vals []string) (map[string]string, error) {
	if len(cols) != len(vals) {
		return nil, fmt.Errorf("length mismatch")
	}

	var ret map[string]string = make(map[string]string, len(cols))
	for i := range cols {
		ret[cols[i]] = vals[i]
	}

	return ret, nil
}
