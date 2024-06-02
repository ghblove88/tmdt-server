package common

import (
	"crypto/md5"
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

/**
 * 缓存管理
 */
var _intance map[string]interface{}
var once sync.Once

func CacheIntance() map[string]interface{} {
	once.Do(func() {
		_intance = make(map[string]interface{})
	})
	return _intance
}

// Strtomd5 create md5 string
func Strtomd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

// Pwdhash password hash function
func PwdHash(str string) string {
	return Strtomd5(str)
}

// seeded indicates if math/rand has been seeded
var seeded bool = false

// uuidRegex matches the UUID string
var uuidRegex *regexp.Regexp = regexp.MustCompile(`^\{?([a-fA-F0-9]{8})-?([a-fA-F0-9]{4})-?([a-fA-F0-9]{4})-?([a-fA-F0-9]{4})-?([a-fA-F0-9]{12})\}?$`)

// UUID type.
type UUID [16]byte

// Hex returns a hex string representation of the UUID in xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx format.
func (this UUID) Hex() string {
	x := [16]byte(this)
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		x[0], x[1], x[2], x[3], x[4],
		x[5], x[6],
		x[7], x[8],
		x[9], x[10], x[11], x[12], x[13], x[14], x[15])

}

// Rand generates a new version 4 UUID.
func Rand() UUID {
	var x [16]byte
	randBytes(x[:])
	x[6] = (x[6] & 0x0F) | 0x40
	x[8] = (x[8] & 0x3F) | 0x80
	return x
}

// FromStr returns a UUID based on a string.
// The string could be in the following format:
//
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//
// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//
// {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
//
// If the string is not in one of these formats, it'll return an error.
func FromStr(s string) (id UUID, err error) {
	if s == "" {
		err = errors.New("Empty string")
		return
	}

	parts := uuidRegex.FindStringSubmatch(s)
	if parts == nil {
		err = errors.New("Invalid string format")
		return
	}

	var array [16]byte
	slice, _ := hex.DecodeString(strings.Join(parts[1:], ""))
	copy(array[:], slice)
	id = array
	return
}

// MustFromStr behaves similarly to FromStr except that it'll panic instead of
// returning an error.
func MustFromStr(s string) UUID {
	id, err := FromStr(s)
	if err != nil {
		panic(err)
	}
	return id
}

// randBytes uses crypto random to get random numbers. If fails then it uses math random.
func randBytes(x []byte) {

	length := len(x)
	n, err := crand.Read(x)

	if n != length || err != nil {
		if !seeded {
			mrand.Seed(time.Now().UnixNano())
		}

		for length > 0 {
			length--
			x[length] = byte(mrand.Int31n(256))
		}
	}
}

func Substring(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
func GetDeviceVoice(deviceid string, deviceinfo string) (voice string) {
	if deviceinfo != "" && len(deviceinfo) >= 4 {
		reg := regexp.MustCompile(`([0-9][0-9][0-9][0-9])$`)
		strVoiceNo := reg.FindString(deviceinfo)
		val, err := strconv.Atoi(strVoiceNo)
		if err == nil && val > 0 && val < 9999 {
			//return strVoiceNo //返回完整的 4位内镜编号
			return strconv.Itoa(val)
		}
	}

	if len(deviceid) >= 4 {
		//return Substring(deviceid, len(deviceid)-4, 4)//返回完整的 4位内镜编号

		tmp := Substring(deviceid, len(deviceid)-4, 4)
		val, err := strconv.Atoi(tmp)
		if err == nil && val > 0 && val < 9999 {
			return strconv.Itoa(val)
		}
	}
	return "0000"
}
func Slice_Remove(slice []string, start, end int) []string {
	return append(slice[:start], slice[end:]...)
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func TermExec(cmdstr string, args []string) (err error) {
	log.Println(cmdstr, args)
	cmd := exec.Command(cmdstr, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println("StderrPipe: ", err.Error())
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Println("Start: ", err.Error())
		return err
	}

	bytesErr, err := io.ReadAll(stderr)
	if err != nil {
		log.Println("ReadAll stderr: ", err.Error())
		return err
	}

	if len(bytesErr) != 0 {
		log.Printf("stderr is not nil: %s", bytesErr)
		return err
	}
	if err := cmd.Wait(); err != nil {
		log.Println("Wait: ", err.Error())
		return err
	}
	return cmd.Run()
}

func CopyFile(src, dst string) (err error) {
	exec.Command("cp", []string{"-f", src, dst}...).Run()
	return
}
func SecToStr(sec int64) (str string) {
	hour := sec / 3600
	minute := sec % 3600 / 60
	second := sec % 3600 % 60

	hour_str := strconv.FormatInt(hour, 10)
	if len(hour_str) <= 1 {
		hour_str = "0" + hour_str
	}
	minute_str := strconv.FormatInt(minute, 10)
	if len(minute_str) <= 1 {
		minute_str = "0" + minute_str
	}
	second_str := strconv.FormatInt(second, 10)
	if len(second_str) <= 1 {
		second_str = "0" + second_str
	}
	return hour_str + ":" + minute_str + ":" + second_str
}
