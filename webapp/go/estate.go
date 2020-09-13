package main

import (
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	estateCache     = cache.New(5*time.Minute, 10*time.Minute)
	estateCachefile = "./chair.gob"
)

func loadEstates() error {
	var estates []Estate
	query := `SELECT id,thumbnail,name,description,latitude,longitude,address,rent,door_height,door_width,features,popularity FROM estate`
	if err := db.Select(&estates, query); err != nil {
		return err
	}
	for _, v := range estates {
		estateCache.Set(strconv.FormatInt(v.ID, 10), v, cache.DefaultExpiration)
	}
	return nil
}

func getEstateByID(id int) (Estate, error) {
	if val, found := estateCache.Get(strconv.Itoa(id)); found {
		return val.(Estate), nil
	}
	estate := Estate{}
	query := `SELECT id,thumbnail,name,description,latitude,longitude,address,rent,door_height,door_width,features,popularity FROM estate WHERE id = ?`
	err := db.Get(&estate, query, id)
	if err != nil {
		estateCache.Set(strconv.Itoa(id), estate, cache.DefaultExpiration)
	}
	return estate, err
}

func setEstatesOnMem(records [][]string) {
	for _, row := range records {
		rm := RecordMapper{Record: row}

		id := rm.NextInt()
		name := rm.NextString()
		description := rm.NextString()
		thumbnail := rm.NextString()
		address := rm.NextString()
		latitude := rm.NextFloat()
		longitude := rm.NextFloat()
		rent := rm.NextInt()
		doorHeight := rm.NextInt()
		doorWidth := rm.NextInt()
		features := rm.NextString()
		popularity := rm.NextInt()

		estate := Estate{
			ID:          int64(id),
			Name:        name,
			Description: description,
			Thumbnail:   thumbnail,
			Address:     address,
			Latitude:    latitude,
			Longitude:   longitude,
			Rent:        int64(rent),
			DoorHeight:  int64(doorHeight),
			DoorWidth:   int64(doorWidth),
			Features:    features,
			Popularity:  int64(popularity),
		}

		estateCache.Set(strconv.Itoa(id), estate, cache.DefaultExpiration)
	}
}

func SaveEstateAll(path string) error {
	items := estateCache.Items()
	estates := make([]Estate, 0)
	for _, v := range items {
		ch := v.Object.(Estate)
		estates = append(estates, ch)
	}
	return Save(path, estates)
}

func RecoverEstateAll(path string) error {
	var newItems []Estate
	if err := Load(path, &newItems); err != nil {
		return err
	}
	for _, v := range newItems {
		estateCache.Set(strconv.FormatInt(v.ID, 10), v, cache.DefaultExpiration)
	}
	return nil
}
