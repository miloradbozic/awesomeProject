package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

const cutoffDate = "2023-06-01"

type Input struct {
	Data struct {
		ContractsByOpstinaIDAndDatumU struct {
			Items []struct {
				ID      string `json:"id"`
				DatumU  string `json:"datumU"`
				PpNaziv string `json:"ppNaziv"`
				CenaEUR int    `json:"cenaEUR"`
				N       string `json:"n"`
			} `json:"items"`
		} `json:"contractsByOpstinaIDAndDatumU"`
	} `json:"data"`
}

type Placemark struct {
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Point       struct {
		Coordinates string `xml:"coordinates"`
	} `xml:"Point"`
}

type KML struct {
	XMLName  xml.Name  `xml:"kml"`
	Document *Document `xml:"Document"`
}

type Document struct {
	Placemarks []Placemark `xml:"Placemark"`
}

func main() {
	// Read the JSON data
	jsonData, err := ioutil.ReadFile("input.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Unmarshal the JSON data
	var input Input
	err = json.Unmarshal(jsonData, &input)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Convert to XML format
	var kml KML
	kml.Document = &Document{}
	for _, item := range input.Data.ContractsByOpstinaIDAndDatumU.Items {
		var nItems []struct {
			LatLon struct {
				Lon float64 `json:"Lon"`
				Lat float64 `json:"Lat"`
			} `json:"latlon"`
			PvNepNaziv string `json:"pvNepNaziv"`
			Pov        int    `json:"pov"`
		}
		err := json.Unmarshal([]byte(item.N), &nItems)
		if err != nil {
			fmt.Println("Error unmarshalling inner JSON:", err)
			return
		}
		for _, nItem := range nItems {
			if item.DatumU > cutoffDate {
				coordinates := fmt.Sprintf("%f,%f,0", nItem.LatLon.Lon, nItem.LatLon.Lat)
				placemark := Placemark{
					Name:        nItem.PvNepNaziv,
					Description: fmt.Sprintf("Datum prodaje: %s, kvadrata: %d, cena/m2: %.2fЄ ", item.DatumU, nItem.Pov, float64(item.CenaEUR)/100.0),
					Point: struct {
						Coordinates string `xml:"coordinates"`
					}{Coordinates: coordinates},
				}
				kml.Document.Placemarks = append(kml.Document.Placemarks, placemark)
			} else {
				fmt.Println(item.DatumU)
			}
		}
	}

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(kml, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to XML:", err)
		return
	}

	// Write to XML file
	err = ioutil.WriteFile("output.kml", xmlData, 0644)
	if err != nil {
		fmt.Println("Error writing KML file:", err)
		return
	}

	fmt.Println("XML data written to output.kml successfully!")
}
