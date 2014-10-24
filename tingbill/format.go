package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
)

func contains(strs []string, val string) bool {
	for _, v := range strs {
		if v == val {
			return true
		}
	}
	return false
}

func (b *Bill) Print() {
	names := make([]string, 0, 100)
	for k := range b.Minutes {
		if !contains(names, k) {
			names = append(names, k)
		}
	}
	for k := range b.Messages {
		if !contains(names, k) {
			names = append(names, k)
		}
	}
	for k := range b.Megabytes {
		if !contains(names, k) {
			names = append(names, k)
		}
	}

	sort.Strings(names)
	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"Name", "Minutes", "Messages", "Megabytes", "Base+Tax", "Total"}
	if *btc {
		header = append(header, fmt.Sprintf("mBTC (%.0f%% Discount)", *btcDiscount))
	}

	table.SetHeader(header)
	var grand float64 = 0
	for _, n := range names {
		minShare := b.Minutes[n] / b.MinutesCount * b.MinuteTotal
		total := minShare
		msgShare := b.Messages[n] / b.MessagesCount * b.MessageTotal
		total += msgShare
		mbShare := b.Megabytes[n] / b.MegabytesCount * b.MegabyteTotal
		total += mbShare
		fee := (b.FeeTotal + b.DeviceTotal) / float64(len(names))
		total += fee
		row := []string{
			n,
			fmt.Sprintf("$%.2f", minShare),
			fmt.Sprintf("$%.2f", msgShare),
			fmt.Sprintf("$%.2f", mbShare),
			fmt.Sprintf("$%.2f", fee),
			fmt.Sprintf("$%.2f", total),
		}
		if *btc {
			row = append(row, mBtc(total))
		}
		table.Append(row)
		grand += total
	}
	row := []string{"", "", "", "", "", fmt.Sprintf("$%.2f", grand)}
	if *btc {
		row = append(row, "")
	}
	table.Append(row)
	table.Render()
}
