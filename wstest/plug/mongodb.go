package Plugins

import (
	"fmt"
	"wstest/comn"
	"strings"
	"time"
)

func MongodbScan(info *common.HostInfo) error {
	if common.IsBrute {
		return nil
	}
	_, err := MongodbUnauth(info)
	if err != nil {
		errlog := fmt.Sprintf("[-] Mongodb %v:%v %v", info.Host, info.Ports, err)
		common.LogError(errlog)
	}
	return err
}

func MongodbUnauth(info *common.HostInfo) (flag bool, err error) {
	flag = false
	// op_msg
	packet1 := []byte{
		0x69, 0x00, 0x00, 0x00, // messageLength
		0x39, 0x00, 0x00, 0x00, // requestID
		0x00, 0x00, 0x00, 0x00, // responseTo
		0xdd, 0x07, 0x00, 0x00, // opCode OP_MSG
		0x00, 0x00, 0x00, 0x00, // flagBits
		// sections db.adminCommand({getLog: "startupWarnings"})
		0x00, 0x54, 0x00, 0x00, 0x00, 0x02, 0x67, 0x65, 0x74, 0x4c, 0x6f, 0x67, 0x00, 0x10, 0x00, 0x00, 0x00, 0x73, 0x74, 0x61, 0x72, 0x74, 0x75, 0x70, 0x57, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x73, 0x00, 0x02, 0x24, 0x64, 0x62, 0x00, 0x06, 0x00, 0x00, 0x00, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x00, 0x03, 0x6c, 0x73, 0x69, 0x64, 0x00, 0x1e, 0x00, 0x00, 0x00, 0x05, 0x69, 0x64, 0x00, 0x10, 0x00, 0x00, 0x00, 0x04, 0x6e, 0x81, 0xf8, 0x8e, 0x37, 0x7b, 0x4c, 0x97, 0x84, 0x4e, 0x90, 0x62, 0x5a, 0x54, 0x3c, 0x93, 0x00, 0x00,
	}
	//op_query
	packet2 := []byte{
		0x48, 0x00, 0x00, 0x00, // messageLength
		0x02, 0x00, 0x00, 0x00, // requestID
		0x00, 0x00, 0x00, 0x00, // responseTo
		0xd4, 0x07, 0x00, 0x00, // opCode OP_QUERY
		0x00, 0x00, 0x00, 0x00, // flags
		0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x24, 0x63, 0x6d, 0x64, 0x00, // fullCollectionName admin.$cmd
		0x00, 0x00, 0x00, 0x00, // numberToSkip
		0x01, 0x00, 0x00, 0x00, // numberToReturn
		// query db.adminCommand({getLog: "startupWarnings"})
		0x21, 0x00, 0x00, 0x00, 0x2, 0x67, 0x65, 0x74, 0x4c, 0x6f, 0x67, 0x00, 0x10, 0x00, 0x00, 0x00, 0x73, 0x74, 0x61, 0x72, 0x74, 0x75, 0x70, 0x57, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x73, 0x00, 0x00,
	}

	realhost := fmt.Sprintf("%s:%v", info.Host, info.Ports)

	checkUnAuth := func(address string, packet []byte) (string, error) {
		conn, err := common.WrapperTcpWithTimeout("tcp", realhost, time.Duration(common.Timeout)*time.Second)
		if err != nil {
			return "", err
		}
		defer conn.Close()
		err = conn.SetReadDeadline(time.Now().Add(time.Duration(common.Timeout) * time.Second))
		if err != nil {
			return "", err
		}
		_, err = conn.Write(packet)
		if err != nil {
			return "", err
		}
		reply := make([]byte, 1024)
		count, err := conn.Read(reply)
		if err != nil {
			return "", err
		}
		return string(reply[0:count]), nil
	}

	// send OP_MSG first
	reply, err := checkUnAuth(realhost, packet1)
	if err != nil {
		reply, err = checkUnAuth(realhost, packet2)
		if err != nil {
			return flag, err
		}
	}
	if strings.Contains(reply, "totalLinesWritten") {
		flag = true
		result := fmt.Sprintf("[+] Mongodb %v unauthorized", realhost)
		common.LogSuccess(result)
	}
	return flag, err
}
