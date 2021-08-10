package seslog

import (
	"encoding/json"
	"os"
)

func ReadOptions() (options Options, err error) {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	options = Options{}
	err = decoder.Decode(&options)
	return
}
