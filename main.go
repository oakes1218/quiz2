package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Data struct {
	Pid      int     `json:"pid"`
	Ppid     int     `json:"ppid"`
	Cmd      string  `json:"cmd"`
	Children []*Data `json:"children,omitempty"`
}

func main() {
	file, err := os.Open("./example.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	datas := make(map[int]*Data)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		pid, ppid := parseIds(record[0], record[1])

		cmd := strings.TrimSpace(record[2])

		data, ok := datas[pid]
		if !ok {
			data = &Data{Pid: pid, Ppid: ppid, Cmd: cmd}
			datas[pid] = data
		} else {
			data.Cmd = cmd
		}

		if ppid != 0 {
			parent, ok := datas[ppid]
			if !ok {
				parent = &Data{Pid: ppid}
				datas[ppid] = parent
			}
			parent.Children = append(parent.Children, data)
		}
	}

	var reslut []*Data
	for _, data := range datas {
		if data.Ppid == 0 {
			reslut = append(reslut, data)
		}
	}

	jsonBytes, err := json.Marshal(reslut)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonBytes))
}

func parseIds(pidStr string, ppidStr string) (int, int) {
	var pid, ppid int
	fmt.Sscanf(pidStr, "%d", &pid)
	fmt.Sscanf(ppidStr, "%d", &ppid)
	return pid, ppid
}
