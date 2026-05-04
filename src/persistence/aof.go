package aof

import (
	"os"
	"fmt"
	"time"
)

func AppendSet (key T, value T, expiry T) {
	f, err := os.OpenFile("appendonly.aof", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}

	_, err = f.WriteString("SET", key, value, expiry)

}


func AppendDelete (key T) {
	f, err := os.OpenFile("appendonly.aof", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	
	_, err = f.WriteString("DELETE", key)
}