package runtime

import (
	"EcdsServer/common"
	"regexp"
	"strconv"
	"strings"
)

type STRING_READER_INFO struct {
	Name  string   //功能名称
	Voice []string //播报功能语音
	Type  string   //对应IP地址 尾值 即：192.168.1.21 的 21
}

type Reader_Info struct {
	reader_map map[string]STRING_READER_INFO
}

func (ri *Reader_Info) Init() {
	ri.reader_map = make(map[string]STRING_READER_INFO)
	readerType := common.Config.GetString("reader.types")
	res := strings.Split(readerType, ",")

	for _, str := range res {
		reg := regexp.MustCompile(`[\p{Han}]+`)
		name := reg.FindString(str)
		reg = regexp.MustCompile(`[a-z]+`)
		voice := reg.FindString(str)

		readerName := common.Config.GetString("reader." + name)
		sub_name := strings.Split(readerName, ",")
		for _, sub_str := range sub_name {
			ri.reader_map[sub_str] = STRING_READER_INFO{name, []string{voice}, sub_str}
		}
	}
	readerXjj := common.Config.GetString("reader.洗镜机")
	res = strings.Split(readerXjj, ",")

	for i, str := range res {
		i++
		xjj_voice := []string{}
		if i <= 9 {
			xjj_voice = []string{"xjj", "N" + strconv.Itoa(i)}
		}
		if i == 10 {
			xjj_voice = []string{"xjj", "shi"}
		}
		if i > 10 && i <= 19 {
			xjj_voice = []string{"xjj", "shi", "N" + strconv.Itoa(i-10)}
		}
		if i == 20 {
			xjj_voice = []string{"xjj", "N2", "shi"}
		}
		if i > 20 && i <= 29 {
			xjj_voice = []string{"xjj", "N2", "shi", "N" + strconv.Itoa(i-20)}
		}
		if i == 30 {
			xjj_voice = []string{"xjj", "N3", "shi"}
		}
		if i > 30 && i <= 39 {
			xjj_voice = []string{"xjj", "N3", "shi", "N" + strconv.Itoa(i-30)}
		}
		if i == 40 {
			xjj_voice = []string{"xjj", "N4", "shi"}
		}
		if i > 40 && i <= 49 {
			xjj_voice = []string{"xjj", "N4", "shi", "N" + strconv.Itoa(i-40)}
		}
		if i == 50 {
			xjj_voice = []string{"xjj", "N5", "shi"}
		}

		ri.reader_map[str] = STRING_READER_INFO{"洗镜机" + strconv.Itoa(i), xjj_voice, str}
	}
}

func (ri *Reader_Info) Find(Ip_Add string) STRING_READER_INFO {
	//log.Println("ipaddr:", Ip_Add)

	reg := regexp.MustCompile(`[0-9.]+`)
	res := strings.Split(reg.FindString(Ip_Add), ".")
	return ri.reader_map[res[3]]
}

func (ri *Reader_Info) FindByName(name string) STRING_READER_INFO {
	for _, v := range ri.reader_map {
		if v.Name == "结束" {
			return v
		}
	}
	return STRING_READER_INFO{"", []string{""}, ""}
}
