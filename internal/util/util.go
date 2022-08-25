package util

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"example.com/url-shortener/internal/model"
)

type Counter struct {
	Counter    uint64
	CounterEnd uint64
	Mu         sync.Mutex
}

func (c *Counter) GetAndIncrease(f *os.File, debug bool) (uint64, uint64) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	curCounter := c.Counter
	curCounterEnd := c.CounterEnd
	// if counter is less-than the assigned end range, increment
	// if not, get a new range
	if curCounter < curCounterEnd {
		c.Counter += 1
	} else {
		c.GetNewRange(f, debug)
	}
	return curCounter, curCounterEnd
}

func (c *Counter) GetNewRange(f *os.File, debug bool) (uint64, uint64) {
	/*
		TODO:
		1. query ZooKeeper/etcd to get next available counter range
	*/
	log.SetOutput(f)
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Counter, c.CounterEnd = 1000000, 2000000

	if debug {
		log.Printf("[DEBUG] Generated new counter range (%v through %v)", c.Counter, c.CounterEnd)
	}

	return c.Counter, c.CounterEnd
}

func LoadCounterRange(f *os.File, debug bool, fileName string) (uint64, uint64) {
	log.SetOutput(f)
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

	if debug {
		log.Printf("[DEBUG] Loaded counter range from file (%v through %v)", counter, counterEnd)
	}

	return uint64(counter), uint64(counterEnd)
}

func SaveCounterRange(f *os.File, debug bool, fileName string, c *Counter) {
	log.SetOutput(f)
	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error saving counter range file: %v", err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%v %v", c.Counter, c.CounterEnd))

	if debug {
		log.Printf("[DEBUG] Saved counter range to file (%v through %v)", c.Counter, c.CounterEnd)
	}
}

func FileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) || info.IsDir() {
		return false
	}
	return true
}

func GenerateUrlSlug(f *os.File, debug bool, c *Counter) string {
	log.SetOutput(f)
	base := uint64(62)
	characterSet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	slug := ""

	// uses mutex to get a unique counter value
	counter, _ := c.GetAndIncrease(f, debug)

	for counter > 0 {
		r := counter % base
		counter /= base
		slug = string(characterSet[r]) + slug

	}
	return slug
}

func IsValidSlug(maxSlugLen int, slug string) bool {
	re := fmt.Sprintf("^[A-Za-z0-9]{0,%v}$", maxSlugLen)
	isBase62 := regexp.MustCompile(re).MatchString
	return isBase62(slug)
}

func IsValidUrl(u model.Url) bool {
	_, err := url.ParseRequestURI(u.Target)
	if err != nil {
		return false
	} else {
		return true
	}
}
