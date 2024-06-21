package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"net"
	"os"
)

type Address struct {
	PoolName string `yaml:"PoolName"`
	IPv4Pool string `yaml:"IPv4Pool"`
}

type AddressEntry struct {
	Address Address `yaml:"Address"`
}

type AddressList struct {
	Addresses []AddressEntry `yaml:"AddressList"`
}

func main() {
	// 定义不同的配置
	configs := []struct {
		totalAddresses int
		pools          int
		filename       string
	}{
		{2500, 50, "addresses_2500.yaml"},
		{5000, 50, "addresses_5000.yaml"},
		{10000, 100, "addresses_10000.yaml"},
		{50000, 100, "addresses_50000.yaml"},
	}

	poolCounter := 1
	baseIP := net.ParseIP("70.70.0.1")

	for _, config := range configs {
		baseIP = generateAddresses(config.totalAddresses, config.pools, config.filename, &poolCounter, baseIP)
	}
}

func generateAddresses(totalAddresses, numPools int, filename string, poolCounter *int, baseIP net.IP) net.IP {
	addressesPerBlock := totalAddresses / numPools
	subnetSize := nextPowerOf2(addressesPerBlock)

	var addressList AddressList

	for i := 0; i < numPools; i++ {
		ip := incrementIP(baseIP, i*subnetSize)
		mask := 32 - log2(subnetSize)
		subnet := fmt.Sprintf("%s/%d", ip.String(), mask)
		poolName := fmt.Sprintf("Pool-Lot-%d", *poolCounter)
		address := Address{
			PoolName: poolName,
			IPv4Pool: subnet,
		}
		addressEntry := AddressEntry{
			Address: address,
		}
		addressList.Addresses = append(addressList.Addresses, addressEntry)
		*poolCounter++
	}

	data, err := yaml.Marshal(&addressList)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return baseIP
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return baseIP
	}

	fmt.Printf("YAML data successfully written to %s\n", filename)

	return incrementIP(baseIP, numPools*subnetSize)
}

func incrementIP(ip net.IP, increment int) net.IP {
	ipv4 := ip.To4()
	val := (int(ipv4[0]) << 24) | (int(ipv4[1]) << 16) | (int(ipv4[2]) << 8) | int(ipv4[3])
	val += increment
	newIP := net.IPv4(byte(val>>24), byte((val>>16)&0xFF), byte((val>>8)&0xFF), byte(val&0xFF))
	return newIP
}

func nextPowerOf2(n int) int {
	if n <= 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

func log2(x int) int {
	l := 0
	for ; x > 1; x >>= 1 {
		l++
	}
	return l
}
