package main

import (
	"fmt"

	hemgw16a "./HEMGW16A"
)

func main() {
	var gw hemgw16a.GW
	gw.Init("http://192.168.3.121", "root", "pass")

	fmt.Printf("NowWatt:%d\n", gw.GetNowWatt())
	fmt.Printf("TodayBuyWatt:%f\n", gw.GetTodayBuyWatt())

}
