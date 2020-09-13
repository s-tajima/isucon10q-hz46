package main

import (
	"encoding/gob"
	"os"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	chairCache     = cache.New(5*time.Minute, 10*time.Minute)
	chairCachefile = "./chair.gob"
)

func loadChairs() error {
	var chairs []Chair
	query := `SELECT * FROM chair`
	if err := db.Select(&chairs, query); err != nil {
		return err
	}
	for _, v := range chairs {
		chairCache.Set(strconv.FormatInt(v.ID, 10), v, cache.DefaultExpiration)
	}
	return nil
}

func getChairByID(id int) (Chair, error) {
	if val, found := chairCache.Get(strconv.Itoa(id)); found {
		return val.(Chair), nil
	}
	chair := Chair{}
	query := `SELECT * FROM chair WHERE id = ?`
	err := db.Get(&chair, query, id)
	if err != nil {
		chairCache.Set(strconv.Itoa(id), chair, cache.DefaultExpiration)
	}
	return chair, err
}

func setChairsOnMem(records [][]string) {
	for _, row := range records {
		rm := RecordMapper{Record: row}

		id := rm.NextInt()
		name := rm.NextString()
		description := rm.NextString()
		thumbnail := rm.NextString()
		price := rm.NextInt()
		height := rm.NextInt()
		width := rm.NextInt()
		depth := rm.NextInt()
		color := rm.NextString()
		features := rm.NextString()
		kind := rm.NextString()
		popularity := rm.NextInt()
		stock := rm.NextInt()

		chair := Chair{
			ID:          int64(id),
			Name:        name,
			Description: description,
			Thumbnail:   thumbnail,
			Price:       int64(price),
			Height:      int64(height),
			Width:       int64(width),
			Depth:       int64(depth),
			Color:       color,
			Features:    features,
			Kind:        kind,
			Popularity:  int64(popularity),
			Stock:       int64(stock),
		}

		chairCache.Set(strconv.Itoa(id), chair, cache.DefaultExpiration)
	}
}

func SaveChairAll(path string) error {
	items := chairCache.Items()
	chairs := make([]Chair, 0)
	for _, v := range items {
		ch := v.Object.(Chair)
		chairs = append(chairs, ch)
	}
	return Save(path, chairs)
}

func RecoverChairAll(path string) error {
	var newItems []Chair
	if err := Load(path, &newItems); err != nil {
		return err
	}
	for _, v := range newItems {
		chairCache.Set(strconv.FormatInt(v.ID, 10), v, cache.DefaultExpiration)
	}
	return nil
}

// Encode via Gob to file
func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Decode Gob file
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
