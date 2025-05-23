package config

import (
	"fmt"
	"testing"
)

func TestIfaceSeelct(t *testing.T) {
	pload := PayloadSender{
		Interface: "Ethernet 4",
		// Local_addr: "192.168.60.227",
		RTimeout: "100ms",
	}

	fmt.Println(pload.ReadTimeout())
	fmt.Println(pload.LocalAddr())
}