package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
)

type Region struct {
	xMid   float64
	yMid   float64
	width  float64
	height float64

	index int
}

type RegionList struct {
	Regions  []Region
	filename string
}

func LoadRegionList(filename string) (RegionList, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file %s: %s", filename, err)
		return RegionList{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var retval []Region
	for scanner.Scan() {
		columns := strings.Fields(scanner.Text())
		if len(columns) != 5 {
			log.Printf("Invalid line: %s in %s", scanner.Text(), filename)
			continue
		}
		region := Region{}
		var err error
		if region.index, err = strconv.Atoi(columns[0]); err != nil {
			log.Printf("Invalid index: %s in %s", columns[0], filename)
			continue
		}
		if region.xMid, err = strconv.ParseFloat(columns[1], 64); err != nil {
			log.Printf("Invalid xMid: %s in %s", columns[1], filename)
			continue
		}
		if region.yMid, err = strconv.ParseFloat(columns[2], 64); err != nil {
			log.Printf("Invalid yMid: %s in %s", columns[2], filename)
			continue
		}
		if region.width, err = strconv.ParseFloat(columns[3], 64); err != nil {
			log.Printf("Invalid width: %s in %s", columns[3], filename)
			continue
		}
		if region.height, err = strconv.ParseFloat(columns[4], 64); err != nil {
			log.Printf("Invalid height: %s in %s", columns[4], filename)
			continue
		}

		if !region.IsValid() {
			log.Printf("Invalid region: %#v in %s", columns, filename)
			continue
		}
		retval = append(retval, region)
	}

	return RegionList{Regions: retval, filename: filename}, nil
}

func (r RegionList) Save() error {
	log.Printf("Saving regions to %s", r.filename)
	file, err := os.OpenFile(r.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("Error creating file %s: %s", r.filename, err)
		return err
	}
	for _, region := range r.Regions {
		fmt.Fprintf(file, "%d %f %f %f %f\n", region.index, region.xMid, region.yMid, region.width, region.height)
	}
	return file.Close()
}

func (r Region) Color() color.Color {
	switch r.index {
	case 0:
		return color.RGBA{255, 128, 64, 255}
	case 1:
		return color.RGBA{255, 0, 0, 255}
	case 2:
		return color.RGBA{0, 255, 0, 255}
	case 3:
		return color.RGBA{0, 0, 255, 255}
	case 4:
		return color.RGBA{255, 255, 0, 255}
	case 5:
		return color.RGBA{255, 0, 255, 255}
	case 6:
		return color.RGBA{0, 255, 255, 255}
	default:
		return color.YCbCr{255, uint8(r.index * 16), uint8(r.index * 16)}
	}
}

func (r Region) IsValid() bool {
	if r.width <= 0 || r.height <= 0 || r.width > 1 || r.height > 1 {
		log.Printf("Invalid width/height: %#v", r)
		return false
	}
	if r.xMid < 0 || r.xMid > 1 || r.yMid < 0 || r.yMid > 1 {
		log.Printf("Invalid x/y mid: %#v", r)
		return false
	}
	if r.xMid-r.width/2 < 0 || r.xMid+r.width/2 > 1 {
		log.Printf("Invalid x range: %#v", r)
		return false
	}
	if r.yMid-r.height/2 < 0 || r.yMid+r.height/2 > 1 {
		log.Printf("Invalid y range: %#v %f %f", r, r.yMid-r.height/2, r.yMid+r.height/2)
		return false
	}

	// TODO: Is this legit? These are too small to be useful
	if r.width < 0.001 || r.height < 0.001 {
		return false
	}
	return true
}

func (r *RegionList) AddRegion(region Region) {
	if region.IsValid() {
		log.Printf("Added new region %#v", region)
		r.Regions = append(r.Regions, region)
		r.Save()
	} else {
		log.Printf("Invalid region: %#v", region)
	}
}

func (r *RegionList) Remove(index int) {
	if index >= 0 && index < len(r.Regions) {
		r.Regions = append(r.Regions[:index], r.Regions[index+1:]...)
		log.Printf("Removed region %d: %#v", index, r.Regions)
		r.Save()
	} else {
		log.Printf("Invalid index: %d", index)
	}
}
