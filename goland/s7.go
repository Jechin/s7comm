package main

import (
	"encoding/binary"
	"fmt"
	"gos7/gos7"
	"log"
	"os"
	"time"
)

const (
	tcpDevice = "127.0.0.1"
	rack      = 0
	slot      = 2
)

func main() {
	// TCPClient
	handler := gos7.NewTCPClientHandler(tcpDevice, rack, slot)
	handler.Timeout = 200 * time.Second
	handler.IdleTimeout = 200 * time.Second
	handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	handler.Connect()
	fmt.Println("Connect Successfully!")
	fmt.Println("-----------------------------------------------")
	defer handler.Close()
	//init client
	client := gos7.NewClient(handler)

	//AGWriteDB to address DB1 with value 111, start from position 8 with size = 2 (for an integer)
	address := 0001
	start := 8
	size := 2
	buffer := make([]byte, 255)
	value := 111
	var helper gos7.Helper
	helper.SetValueAt(buffer, 0, uint16(value))
	err := client.AGWriteDB(address, start, size, buffer)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Write DB Successfully!")
		fmt.Println("-----------------------------------------------")
	}

	// AGReadDB to address DB1, start from position 8 with size = 2
	buf := make([]byte, 255)
	err = client.AGReadDB(address, start, size, buf)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Read DB Successfully!")
		var result uint16
		helper.GetValueAt(buf, 0, &result)
		fmt.Println(result)
		fmt.Println("-----------------------------------------------")
	}

	// Get PLC Status
	status, err := client.PLCGetStatus()
	if err != nil {
		fmt.Println(err)
	} else {
		// 8=running, 4=stop, 0=unknown
		switch status {
		case 8:
			fmt.Println("PLC status: running")
		case 4:
			fmt.Println("PLC status: stop")
		case 0:
			fmt.Println("PLC status: unknown")
		}
		fmt.Println("-----------------------------------------------")
	}

	// List Blocks
	bl, err := client.PGListBlocks()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("List Blocks: ")
		fmt.Println("OB: ", bl.OBList)
		fmt.Println("FB: ", bl.FBList)
		fmt.Println("FC: ", bl.FCList)
		fmt.Println("SFB: ", bl.SFBList)
		fmt.Println("SFC: ", bl.SFCList)
		fmt.Println("DB: ", bl.DBList)
		fmt.Println("SDB: ", bl.SDBList)
		fmt.Println("-----------------------------------------------")
	}

	info, err := client.GetCPUInfo()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("CPU Info: ")
		fmt.Println("ModuleTypeName: ", info.ModuleTypeName)
		fmt.Println("SerialNumber: ", info.SerialNumber)
		fmt.Println("ASName: ", info.ASName)
		fmt.Println("CopyRight: ", info.Copyright)
		fmt.Println("ModuleName: ", info.ModuleName)
		fmt.Println("-----------------------------------------------")
	}

	// GetAGBlockInfo DB1 ; blocktype 65 means DB
	blockinfo, err := client.GetAgBlockInfo(65, 1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Block Info: ")
		fmt.Println("BlkType: ", blockinfo.BlkType)
		fmt.Println("BlkNumber: ", blockinfo.BlkNumber)
		fmt.Println("BlkLang: ", blockinfo.BlkLang)
		fmt.Println("BlkFlags: ", blockinfo.BlkFlags)
		fmt.Println("MC7Size: ", blockinfo.MC7Size)
		fmt.Println("LoadSize: ", blockinfo.LoadSize)
		fmt.Println("LocalData: ", blockinfo.LocalData)
		fmt.Println("SBBLength: ", blockinfo.SBBLength)
		fmt.Println("CheckSum: ", blockinfo.CheckSum)
		fmt.Println("Version: ", blockinfo.Version)
		fmt.Println("CodeDate: ", blockinfo.CodeDate)
		fmt.Println("IntfDate: ", blockinfo.IntfDate)
		fmt.Println("Author: ", blockinfo.Author)
		fmt.Println("Family: ", blockinfo.Family)
		fmt.Println("Header: ", blockinfo.Header)
		fmt.Println("-----------------------------------------------")
	}

	// MultiWrite
	data1 := make([]byte, 1024)
	data2 := make([]byte, 1024)
	data3 := make([]byte, 1024)

	for i := 0; i < 16; i++ {
		data1[i] = 0x01
		data2[i] = 0x02
		data3[i] = 0x03
	}
	var error1, error2, error3 string

	var items = []gos7.S7DataItem{
		gos7.S7DataItem{
			Area:     0x84,
			WordLen:  0x02,
			DBNumber: 1,
			Start:    0,
			Amount:   16,
			Data:     data1,
			Error:    error1,
		},
		gos7.S7DataItem{
			Area:     0x84,
			WordLen:  0x02,
			DBNumber: 2,
			Start:    0,
			Amount:   16,
			Data:     data2,
			Error:    error2,
		},
		gos7.S7DataItem{
			Area:     0x84,
			WordLen:  0x02,
			DBNumber: 3,
			Start:    0,
			Amount:   16,
			Data:     data3,
			Error:    error3,
		},
	}
	err = client.AGWriteMulti(items, 3)
	if err != nil {
		fmt.Println(err)
	} else if error1 != "" || error2 != "" || error3 != "" {
		fmt.Println("Error1: ", error1)
		fmt.Println("Error2: ", error2)
		fmt.Println("Error3: ", error3)
	} else {
		fmt.Println("MultiWrite Successfully!")
		fmt.Println("-----------------------------------------------")
	}

	for i := 0; i < 16; i++ {
		data1[i] = 0x00
		data2[i] = 0x00
		data3[i] = 0x00
	}
	err = client.AGReadMulti(items, 3)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("MultiRead Successfully!")
		fmt.Println("Data1: ", binary.BigEndian.Uint16(data1[8:]))
		fmt.Println("Data2: ", binary.BigEndian.Uint16(data2[8:]))
		fmt.Println("Data3: ", binary.BigEndian.Uint16(data3[8:]))
		fmt.Println("-----------------------------------------------")
	}
}
