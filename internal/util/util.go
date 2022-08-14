package util

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"example.com/url-shortener/internal/config"
)

type Counter struct {
	Counter    uint64
	CounterEnd uint64
	Mu         sync.Mutex
}

func (c *Counter) GetAndIncrease() (uint64, uint64) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	curCounter := c.Counter
	curCounterEnd := c.CounterEnd
	// if counter is less-than the assigned end range, increment
	// if not, get a new range
	if curCounter < curCounterEnd {
		c.Counter += 1
	} else {
		c.GetNewRange()
	}
	return curCounter, curCounterEnd
}

func (c *Counter) GetNewRange() (uint64, uint64) {
	/*
		TODO:
		1. query ZooKeeper to get next available counter range
	*/
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Counter, c.CounterEnd = 1000000, 2000000
	return c.Counter, c.CounterEnd
}

func LoadCounterRange(fileName string) (uint64, uint64) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	dataString := strings.Split(string(data), " ")

	counter, err := strconv.Atoi(dataString[0])
	if err != nil {
		log.Fatalln(err)
	}
	counterEnd, err := strconv.Atoi(dataString[1])
	if err != nil {
		log.Fatalln(err)
	}

	// delete counter range file after loading
	// prevents loading inconsistent range after crash
	err = os.Remove(fileName)
	if err != nil {
		log.Printf("Error removing counter range file: %v", err)
	}

	return uint64(counter), uint64(counterEnd)
}

func SaveCounterRange(fileName string, c *Counter) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error saving counter range file: %v", err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%v %v", c.Counter, c.CounterEnd))
}

func FileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) || info.IsDir() {
		return false
	}
	return true
}

func GenerateUrlSlug(c *Counter) string {
	base := uint64(62)
	characterSet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	slug := ""

	// uses mutex to get a unique counter value
	counter, _ := c.GetAndIncrease()

	for counter > 0 {
		r := counter % base
		counter /= base
		slug = string(characterSet[r]) + slug

	}
	return slug
}

func IsValidSlug(slug string) bool {
	re := fmt.Sprintf("^[A-Za-z0-9]{0,%v}$", config.MaxSlugLen)
	isBase62 := regexp.MustCompile(re).MatchString
	return isBase62(slug)
}
