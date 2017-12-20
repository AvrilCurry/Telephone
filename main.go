package main

import (
	"Telephone/config"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/mux"
)

type telephone struct {
	ID           int      `json:"id"`
	Brand        string   `json:"brand"`
	LowestPrice  int      `json:"lowestPrice"`
	HighestPrice int      `json:"highestPrice"`
	Color        []string `json:"color"`
}

// Index The first page
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to access \"telephone.com\"!\n")
}

// AddTelephone add an telephone infomation to the back end
func AddTelephone(w http.ResponseWriter, r *http.Request) {
	var tele, dataTele telephone
	var data []byte
	var isDataExisted = false

	// 读取传入的数据， 1048576=1M, 限制传入数据的大小
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(2)
		panic(err)
	}
	if err = r.Body.Close(); err != nil {
		log.Println(3)
		panic(err)
	}

	err = json.Unmarshal(body, &tele)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err = json.NewEncoder(w).Encode(err); err != nil {
			log.Println(4)
			panic(err)
		}
	}

	// 将数据写入data.json文件
	fout, err := os.OpenFile("./data.json", os.O_RDWR|os.O_APPEND, os.ModePerm)

	if err != nil {
		log.Fatal(11)
		panic(err)
	}

	buffReader := bufio.NewReader(fout)
	buffWriter := bufio.NewWriter(fout)

	for true {
		data, err = buffReader.ReadBytes('\n')

		if err != nil {
			if err != io.EOF {
				log.Println(5)
				panic(err)
			}
			break
		}

		err = json.Unmarshal(data, &dataTele)

		// 判断数据是否已经存在
		if tele.Brand == dataTele.Brand {
			isDataExisted = true
			break
		}
	}

	if !isDataExisted {
		tele.ID = config.GetID()
		config.SetID()
		data, _ = json.Marshal(tele)
		buffWriter.Write(data)
		buffWriter.WriteString("\n")
		buffWriter.Flush()
		fout.Close()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8") // 设置response的头部信息
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(tele); err != nil {
			panic(err)
		}
	}
}

// listTelephone list all the telephone information in the back end
func listTelephone(w http.ResponseWriter, r *http.Request) {
	var data []byte
	var tele telephone

	fin, err := os.OpenFile("./data.json", os.O_RDWR, os.ModePerm)

	if err != nil {
		os.Exit(1)
	}

	buffReader := bufio.NewReader(fin)

	for true {
		data, err = buffReader.ReadBytes('\n')

		if err != nil {
			if err != io.EOF {
				os.Exit(2)
			}
			break
		}

		err = json.Unmarshal(data, &tele)

		CommandInfo := "\tID:\t%d\n\tBrand:\t%s\n\tLowest-Price:\t%d\n\tHighest_Price:\t%d\n\tColor:\t%v\n\n"
		fmt.Fprintf(w, CommandInfo, tele.ID, tele.Brand, tele.LowestPrice, tele.HighestPrice, tele.Color)
	}
}

// searchTelephone search the telephone information according to the searchInfo
func searchTelephone(w http.ResponseWriter, r *http.Request) {
	var tele telephone
	var data []byte
	var ColorArray []string

	// 判断传入的是什么参数
	vars := mux.Vars(r)
	brand := vars["brandName"]
	lowestPrice := vars["low"]
	highestPrice := vars["high"]
	color := vars["colorArray"]

	fmt.Println(brand, lowestPrice, highestPrice, color)

	if len(color) != 0 {
		ColorArray = strings.Split(color, ",")
	}

	// fmt.Println(brand, lowestPrice, highestPrice, len(color), ColorArray, len(ColorArray))

	fin, err := os.OpenFile("./data.json", os.O_RDWR, os.ModePerm)

	if err != nil {
		panic(err)
	}

	buffReader := bufio.NewReader(fin)

	for true {
		var isBrandExisted = false
		var isPriceExisted = false
		var isColorExisted = false
		data, err = buffReader.ReadBytes('\n')

		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}

		err = json.Unmarshal(data, &tele)

		if brand == tele.Brand {
			isBrandExisted = true
		}
		// 将字符串转为数字
		low, _ := strconv.Atoi(lowestPrice)
		high, _ := strconv.Atoi(highestPrice)
		if low <= tele.LowestPrice && high >= tele.HighestPrice {
			isPriceExisted = true
		}
		var length = 0
		for _, item := range ColorArray {
			for _, items := range tele.Color {
				if item == items {
					length = length + 1
					break
				}
			}
		}
		if len(ColorArray) != 0 && len(ColorArray) == length {
			isColorExisted = true
		}

		//fmt.Println(low, high, tele.LowestPrice, tele.HighestPrice)
		//fmt.Println(len(ColorArray), isBrandExisted, isPriceExisted, isColorExisted)

		if isBrandExisted || isPriceExisted || isColorExisted {
			CommandInfo := "\tID:\t%d\n\tBrand:\t%s\n\tLowest-Price:\t%d\n\tHighest_Price:\t%d\n\tColor:\t%v\n\n"
			fmt.Fprintf(w, CommandInfo, tele.ID, tele.Brand, tele.LowestPrice, tele.HighestPrice, tele.Color)
		}
	}

}

func main() {
	// 使用了 gorilla/mux 定义了一个路由器
	router := mux.NewRouter().StrictSlash(true)

	// 定义路由规则  Index,AddTelephone,listTelephone,searchTelephone都是处理器，处理对应的请求
	router.HandleFunc("/", Index).Methods("Get")
	router.HandleFunc("/add", AddTelephone).Methods("Post")
	router.HandleFunc("/search", listTelephone).Methods("Get")
	router.HandleFunc("/search/brand={brandName}", searchTelephone).Methods("Get")
	router.HandleFunc("/search/lowestPrice={low}&highestPrice={high}", searchTelephone).Methods("Get")
	router.HandleFunc("/search/color={colorArray}", searchTelephone).Methods("Get")

	fout, err := os.OpenFile("./logInfo.log", os.O_RDWR|os.O_APPEND, os.ModePerm)

	if err != nil {
		os.Exit(1)
	}

	MyLogger := log.New(fout, "[Log]: ", log.Ldate|log.Ltime)

	MyLogger.Fatal(http.ListenAndServe(":8181", router))
}
