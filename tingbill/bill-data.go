package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Bill struct {
	Link           string
	Title          string
	ID             string
	Error          error
	Minutes        map[string]float64
	Messages       map[string]float64
	Megabytes      map[string]float64
	MinutesCount   float64
	MessagesCount  float64
	MegabytesCount float64
	MinuteTotal    float64
	MessageTotal   float64
	MegabyteTotal  float64
	DeviceTotal    float64
	FeeTotal       float64
	GrandTotal     float64
}

var space = regexp.MustCompile("\\s+")
var price = regexp.MustCompile("[^0-9.]")
var nondigit = regexp.MustCompile("[^0-9]")
var brackets = regexp.MustCompile("\\(.*\\)")
var phone = regexp.MustCompile("^[0-9-]+$")

func GetBills() ([]Bill, error) {
	r, err := client.Get("https://ting.com/account/bill_history")
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(r)
	if err != nil {
		return nil, err
	}
	allBills := doc.Find("td[data=bill] > a")
	bills := make([]Bill, 0, allBills.Size())

	allBills.Each(func(i int, sel *goquery.Selection) {
		if i >= int(*billCount) {
			return
		}
		title := sel.Text()
		link, _ := sel.Attr("href")
		bills = append(bills, Bill{
			Link:      link,
			Title:     title,
			ID:        filepath.Base(link),
			Minutes:   make(map[string]float64, 10),
			Megabytes: make(map[string]float64, 10),
			Messages:  make(map[string]float64, 10),
		})
	})

	var wg sync.WaitGroup
	for k := range bills {
		wg.Add(2)
		go GetDetails(bills, k, &wg)
		go GetUsage(bills, k, &wg)
	}
	wg.Wait()
	return bills, nil
}

func GetUsage(bills []Bill, index int, wg *sync.WaitGroup) {
	b := &bills[index]
	defer func() {
		wg.Done()
	}()
	fmt.Println("Get Usage for:", b.Title)
	r, err := client.Get("https://ting.com/account/usage?period_id=" + b.ID)
	if err != nil {
		b.Error = err
		return
	}
	doc, err := goquery.NewDocumentFromResponse(r)
	if err != nil {
		b.Error = err
		return
	}
	doc.Find("#minutesTable > tbody > tr").Each(func(i int, sel *goquery.Selection) {
		cols := sel.Find("td")
		var device string
		if sel.HasClass("outgoing") {
			device = cols.Eq(2).Text()
		} else {
			device = cols.Eq(3).Text()
		}
		minStr := cols.Eq(4).Text()

		device = brackets.ReplaceAllString(device, "")
		device = strings.TrimSpace(device)
		minStr = nondigit.ReplaceAllString(minStr, "")
		mins, err := strconv.ParseFloat(minStr, 64)
		if err != nil {
			return
		}

		b.Minutes[device] += mins
		b.MinutesCount += mins
	})
	doc.Find("#messageTable > tbody > tr").Each(func(i int, sel *goquery.Selection) {
		cols := sel.Find("td")
		from := strings.TrimSpace(cols.Eq(2).Text())
		to := strings.TrimSpace(cols.Eq(3).Text())
		var device string
		if phone.MatchString(from) {
			device = to
		} else {
			device = from
		}

		b.Messages[device]++
		b.MessagesCount++
	})
	doc.Find("#megabytesTable > tbody > tr").Each(func(i int, sel *goquery.Selection) {
		cols := sel.Find("td")
		device := strings.TrimSpace(cols.Eq(1).Text())
		usageStr := nondigit.ReplaceAllString(cols.Eq(2).Text(), "")
		usage, err := strconv.ParseFloat(usageStr, 64)
		if err != nil {
			return
		}

		b.Megabytes[device] += usage
		b.MegabytesCount += usage
	})
}

func GetDetails(bills []Bill, index int, wg *sync.WaitGroup) {
	b := &bills[index]
	defer func() {
		wg.Done()
	}()
	fmt.Println("Get Details for:", b.Title)
	r, err := client.Get(b.Link)
	if err != nil {
		b.Error = err
		return
	}
	doc, err := goquery.NewDocumentFromResponse(r)
	if err != nil {
		b.Error = err
		return
	}
	doc.Find("#billingDetailsTable > tbody > tr").Each(func(i int, sel *goquery.Selection) {
		name := strings.TrimSpace(sel.Find("td").Eq(1).Text())
		amtStr := price.ReplaceAllString(sel.Find("td").Eq(3).Text(), "")
		amt, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			return
		}
		name = space.ReplaceAllString(name, " ")
		if strings.HasSuffix(name, "Minutes") {
			b.MinuteTotal = amt
		} else if strings.HasSuffix(name, "Messages") {
			b.MessageTotal = amt
		} else if strings.HasSuffix(name, "Megabytes") {
			b.MegabyteTotal = amt
		} else if strings.HasSuffix(name, "Devices") {
			b.DeviceTotal = amt
		}
	})
	feeStr := doc.Find("p.totalLine span.totalAmount").Eq(1).Text()
	fee, err := strconv.ParseFloat(price.ReplaceAllString(feeStr, ""), 64)
	if err == nil {
		b.FeeTotal = fee
	}
	b.GrandTotal = b.MinuteTotal + b.MessageTotal + b.MegabyteTotal + b.DeviceTotal + b.FeeTotal
}
