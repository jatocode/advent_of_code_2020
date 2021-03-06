package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(path string) []string {
	file, err := os.Open(path)
	defer file.Close()
	check(err)

	result := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}

	return result
}

type machine struct {
	maskOr  uint
	maskAnd uint

	values map[uint]uint
}

func (m *machine) setMask(mask string) {
	m.maskAnd = ^uint(0)
	m.maskOr = 0

	for offset, bit := range mask {
		if bit == '1' {
			m.maskOr += 1 << (len(mask) - offset - 1)
		} else if bit == '0' {
			m.maskAnd -= 1 << (len(mask) - offset - 1)
		}
	}

}

func (m *machine) setValue(address, value uint) {
	if m.values == nil {
		m.values = map[uint]uint{}
	}

	m.values[address] = value&m.maskAnd | m.maskOr
}

func (m machine) sum() uint {
	sum := uint(0)
	for _, value := range m.values {
		sum += value
	}
	return sum
}

func main() {

	data := readFile("./data.txt")

	machine := machine{}

	for _, row := range data {

		rmask := regexp.MustCompile("^mask = ([X10]+)$")
		rmem := regexp.MustCompile("^mem\\[(\\d+)\\] = (\\d+)$")

		maskMatches := rmask.FindStringSubmatch(row)
		memMatches := rmem.FindStringSubmatch(row)

		if maskMatches != nil {
			mask := maskMatches[1]
			machine.setMask(mask)
		} else if memMatches != nil {
			address, _ := strconv.Atoi(memMatches[1])
			value, _ := strconv.Atoi(memMatches[2])
			machine.setValue(uint(address), uint(value))
		} else {
			panic("No match: " + row)
		}

	}

	fmt.Println(machine.sum())
}
