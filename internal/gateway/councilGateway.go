package gateway

import (
	"bintracker/internal/bintracker"
	"encoding/json"
)

func BinController(config bintracker.Config) Response {
	return BuildMockResponse()
}

func BuildMockResponse() Response {
	jsonData := `{
    "status": "CollectionDatesFound",
    "message": "4 collection date(s) found for property 123456789",
    "lastWeek": [
        {
            "date": "2026-07-30T00:00:00",
            "dayOfWeekName": "Thursday",
            "bins": [
                {
                    "colourLabel": "Grey",
                    "colour": "Grey",
                    "name": "General waste bin",
                    "capacity": "240L",
                    "type": "BIN",
                    "uprn": "123456789"
                }
            ]
        }
    ],
    "thisWeek": [
        {
            "date": "2026-08-04T00:00:00",
            "dayOfWeekName": "Tuesday",
            "bins": [
                {
                    "colourLabel": "Brown/green",
                    "colour": "Brown",
                    "name": "Garden and food waste bin",
                    "capacity": "240L",
                    "type": "BIN",
                    "uprn": "123456789"
                },
                {
                    "colourLabel": "",
                    "colour": "Yellow",
                    "name": "Glass container",
                    "capacity": "240L",
                    "type": "CONTAINER",
                    "uprn": "123456789"
                }
            ]
        },
        {
            "date": "2026-08-06T00:00:00",
            "dayOfWeekName": "Thursday",
            "bins": [
                {
                    "colourLabel": "Blue",
                    "colour": "Blue",
                    "name": "Recycling bin",
                    "capacity": "240L",
                    "type": "BIN",
                    "uprn": "123456789"
                }
            ]
        }
    ],
    "nextWeek": [
        {
            "date": "2026-08-13T00:00:00",
            "dayOfWeekName": "Thursday",
            "bins": [
                {
                    "colourLabel": "Grey",
                    "colour": "Grey",
                    "name": "General waste bin",
                    "capacity": "240L",
                    "type": "BIN",
                    "uprn": "123456789"
                }
            ]
        }
    ]
}`

	var response Response
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		panic(err)
	}
	return response

}
