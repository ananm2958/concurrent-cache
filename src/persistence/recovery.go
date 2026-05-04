package recovery

import (
	"os"
	"fmt"
	"time"
)

func main (){

	data, err := os.ReadFile("snapshot.json")

	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var entries []Entry
	json.Unmarshal()

	s.cache.mutex.Lock()
	defer s.cache.mutex.Unlock()

	for _, entry := range entries {
		if !entry.Expiry.IsZero() && time.Now().After(entry.Expiry) {
		continue

		}
		s.cache.Set(entry.key, entry.value)

		if node, ok := c.data[entry.Key]; ok {
		node.expiry = entry.Expiry
		}
	}

}