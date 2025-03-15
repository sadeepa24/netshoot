package config

import (
	"fmt"
	"testing"
)

func TestIfaceSeelct(t *testing.T) {
	pload := PayloadSender{
		Interface: "Ethernet 2",
		Local_addr: "192.168.60.227",
	}

	fmt.Println(pload.LocalAddr())
}