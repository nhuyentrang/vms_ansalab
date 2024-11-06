package hikivision

import (
	"encoding/xml"
	"fmt"
	"net"
	"time"
)

var (
// Host     string = "http://192.168.100.14:8888"
// Username string = "admin"
// Password string = "123456aA"
)

var ClientHik *Client

func ConnectToHikvisionDevice(host string, username string, password string) (*Client, error) {
	url := fmt.Sprintf("%s:%s@%s", username, password, host)
	// client, err := NewClient("admin:123456aA@192.168.100.28", "", "")
	client, err := NewClient(url, "", "")
	if err != nil {
		return nil, err
	}
	ClientHik = client
	fmt.Println("======> Hikvision connect !!! ")
	return client, nil
}

func ScanDeviceHikvision() {
	// Create a UDP socket and bind it to a local IP address and port number
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Send a Hikvision discovery message over the UDP socket to the target network
	dstAddr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:37020")
	if err != nil {
		panic(err)
	}
	msg := `<Probe><types>inquiry</types></Probe>`
	_, err = conn.WriteTo([]byte(msg), dstAddr)
	if err != nil {
		panic(err)
	}

	// Receive and parse the responses from the Hikvision IP cameras that reply to the multicast message
	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(15 * time.Second)) // Set a timeout for reading responses
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Timeout, no more responses
				break
			}
			panic(err)
		}
		fmt.Printf("Received response from %v:\n%s\n", addr, string(buf[:n]))
		// Parse the response XML and extract the device information
		var response struct {
			XMLName xml.Name `xml:"Envelope"`
			Body    struct {
				ProbeMatches []struct {
					XMLName  xml.Name `xml:"ProbeMatch"`
					Endpoint struct {
						Address string `xml:"Address,attr"`
						Port    string `xml:"Port,attr"`
					} `xml:"XAddrs"`
					Types string `xml:"Types"`
				} `xml:"ProbeMatches>ProbeMatch"`
			} `xml:"Body"`
		}
		err = xml.Unmarshal(buf[:n], &response)
		if err != nil {
			panic(err)
		}
		for _, match := range response.Body.ProbeMatches {
			fmt.Printf("Device: %s:%s (Types: %s)\n", match.Endpoint.Address, match.Endpoint.Port, match.Types)
		}
	}
}
