package short

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"short-url/bolt"
	"short-url/conf"
	"short-url/dto"
	"strings"
)

const BaseString = "UgtC0K9wX5IdJcGYBh4QVil1oPakrpMEzAye87Ds3FujOTb6nx2LZNvWmHqSfR"
const BaseStringLength = uint64(len(BaseString))
const root = "shorturl"

//const BaseStringLength = uint64(len("Ds3K9ZNvWmHcakr1oPnxh4qpMEzAye8wX5IdJ2LFujUgtC07lOTb6GYBQViSfR"))

// Int2String converts an unsigned 64bit integer to a string.
func Int2String(seq uint64) (shortURL string) {
	charSeq := []rune{}
	if seq != 0 {
		for seq != 0 {
			mod := seq % BaseStringLength
			div := seq / BaseStringLength
			charSeq = append(charSeq, rune(BaseString[mod]))
			seq = div
		}
	} else {
		charSeq = append(charSeq, rune(BaseString[seq]))
	}

	tmpShortURL := string(charSeq)
	shortURL = reverse(tmpShortURL)
	return
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// String2Int converts a short URL string to an unsigned 64bit integer.
func String2Int(shortURL string) (seq uint64) {
	shortURL = reverse(shortURL)
	for index, char := range shortURL {
		base := uint64(math.Pow(float64(BaseStringLength), float64(index)))
		seq += uint64(strings.Index(BaseString, string(char))) * base
	}
	return
}

func Short(longURL string) (string, error) {
	boltClient := bolt.NewBoltClient(conf.Conf.DBName, 0600)
	seq, err := boltClient.NextSequence(root)
	if err != nil {
		log.Printf("get next sequence error. %v", err)
		return "", errors.New("get next sequence error")
	}
	shortURL := Int2String(seq)
	compress, err := json.Marshal(dto.CompressDTO{ShortURL: shortURL, LongURL: longURL, ClickTime: 0})
	if err != nil {
		return "", err
	}
	err = boltClient.Set(root, shortURL, compress)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}
