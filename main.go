package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/sony/sonyflake"
)

func machineID() (uint16, error) {
	ipStr := os.Getenv("MY_IP")
	if len(ipStr) == 0 {
		return 0, errors.New("MY_IP env variable is not set")
	}
	ip := net.ParseIP(ipStr)
	if len(ip) < 1 {
		return 0, errors.New("Invalid IP")
	}

	return uint16(ip[2])<<8 + uint16(ip[2]), nil
}

func main() {
	st := sonyflake.Settings{}
	st.MachineID = machineID
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		log.Fatal("Failed to initializae sonyflake")
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		// Generate a new ID
		id, err := sf.NextID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {
			// Return ID as a string
			c.JSON(http.StatusOK, gin.H{
				"id": fmt.Sprint(id),
			})
		}
	})

	if err := r.Run(":3000"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
