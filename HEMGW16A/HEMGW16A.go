package hemgw16a

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
)

// GW HEM-GW16A
type GW struct {
	rootURL  string
	username string
	password string
	attr     smartMeterAttribute
}

type smartMeterAttribute struct {
	macAddr                             string
	modelname                           string
	compositeTransformationRatio        float64
	multiplyingForRatio                 float64
	effectiveDigitsForCumulativeAmounts float64
	unitForCumulativeAmounts            float64
}

func (hem *GW) loadSmartMeterAttribute() {
	url := hem.rootURL + "/php/smartmeter_setting_get.php"
	refererURL := hem.rootURL + "/smartmeter.html"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Referer", refererURL)
	req.SetBasicAuth(hem.username, hem.password)

	client := new(http.Client)
	resp, _ := client.Do(req)
	byteArray, _ := ioutil.ReadAll(resp.Body)

	var attr smartMeterAttribute
	var err error = nil

	var value []byte

	value, _, _, err = jsonparser.Get(byteArray, "smartmeter", "[0]", "hwaddr")
	attr.macAddr = string(value)

	value, _, _, err = jsonparser.Get(byteArray, "smartmeter", "[0]", "modelname")
	attr.modelname = string(value)

	value, _, _, err = jsonparser.Get(byteArray, "smartmeter", "[0]", "attribute", "composite_transformation_ratio")
	attr.compositeTransformationRatio, err = strconv.ParseFloat(string(value), 64)

	value, _, _, err = jsonparser.Get(byteArray, "smartmeter", "[0]", "attribute", "multiplying_for_ratio")
	attr.multiplyingForRatio, err = strconv.ParseFloat(string(value), 64)

	value, _, _, err = jsonparser.Get(byteArray, "smartmeter", "[0]", "attribute", "effective_digits_for_cumulative_amounts")
	attr.effectiveDigitsForCumulativeAmounts, err = strconv.ParseFloat(string(value), 64)

	value, _, _, err = jsonparser.Get(byteArray, "smartmeter", "[0]", "attribute", "unit_for_cumulative_amounts")
	attr.unitForCumulativeAmounts, err = strconv.ParseFloat(string(value), 64)

	if err != nil {
		panic(errors.New("json decode error"))
	}

	hem.attr = attr
}

func (hem GW) getPropertyM(epc string) []byte {
	url := hem.rootURL + "/php/get_property_m.php"
	refererURL := hem.rootURL + "/smartmeter.html"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Referer", refererURL)
	req.SetBasicAuth(hem.username, hem.password)

	params := req.URL.Query()
	params.Add("ba", hem.attr.macAddr)
	params.Add("eoj", "0x028801")
	params.Add("epc", epc)
	req.URL.RawQuery = params.Encode()

	client := new(http.Client)
	resp, err := client.Do(req)

	byteArray, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(errors.New("getPropertyM error :" + url + " " + string(byteArray)))
	}

	ret, _ := hex.DecodeString(strings.Replace(string(byteArray), "0x", "", 1))

	return ret
}

func (hem GW) setPropertyM(epc string, data string) bool {
	url := hem.rootURL + "/php/set_property_m.php"
	refererURL := hem.rootURL + "/smartmeter.html"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Referer", refererURL)
	req.SetBasicAuth(hem.username, hem.password)

	params := req.URL.Query()
	params.Add("ba", hem.attr.macAddr)
	params.Add("eoj", "0x028801")
	params.Add("epc", epc)
	params.Add("data", data)
	req.URL.RawQuery = params.Encode()

	client := new(http.Client)
	resp, _ := client.Do(req)

	byteArray, _ := ioutil.ReadAll(resp.Body)

	ret := string(byteArray) == "ok"

	return ret
}

// GetNowWatt get now watt
func (hem GW) GetNowWatt() uint64 {
	ret := hem.getPropertyM("0xE7")
	padding := make([]byte, 8-len(ret))
	i := binary.BigEndian.Uint64(append(padding, ret...))
	return i
}

// GetTodayBuyWatt get today buy watt
func (hem GW) GetTodayBuyWatt() float64 {
	// 収集日を0(当日)に設定
	hem.setPropertyM("0xE5", "0x00")

	// 収集日を確認
	var day uint16
	res := hem.getPropertyM("0xE2")
	binary.Read(bytes.NewBuffer(res[0:2]), binary.BigEndian, &day)
	if day != 0 {
		panic(errors.New("set date error"))
	}

	// 積算電力量履歴(正方向)を取得
	var zeroClockWattHour uint32
	binary.Read(bytes.NewBuffer(res[2:6]), binary.BigEndian, &zeroClockWattHour)

	// 現在の積算電力量(正方向)を取得
	var wattHourNow uint32
	res = hem.getPropertyM("0xE0")
	binary.Read(bytes.NewBuffer(res[0:4]), binary.BigEndian, &wattHourNow)

	// 計算
	var diff float64 = float64(wattHourNow - zeroClockWattHour)
	if diff < 0 {
		diff += math.Pow(10, hem.attr.effectiveDigitsForCumulativeAmounts)
	}

	var TodayBuyWatt float64
	TodayBuyWatt = float64(diff) * (hem.attr.compositeTransformationRatio * hem.attr.multiplyingForRatio) * hem.attr.unitForCumulativeAmounts * 1000

	return TodayBuyWatt
}

// Init setup HEMGW16A credential
func (hem *GW) Init(rootURL string, username string, password string) {
	hem.rootURL = rootURL
	hem.username = username
	hem.password = password
	hem.loadSmartMeterAttribute()
}
