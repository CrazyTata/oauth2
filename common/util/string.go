package util

import (
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func MD5(str string) string {
	md5Data := md5.Sum([]byte(str))
	return hex.EncodeToString(md5Data[:])
}

func Unique[T comparable](s []T) []T {
	m := make(map[T]bool)
	uniq := make([]T, 0)
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = true
			uniq = append(uniq, v)
		}
	}
	return uniq
}

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// 将任意类型转string
func convertAnyToStr(v interface{}) string {
	if v == nil {
		return ""
	}
	switch d := v.(type) {
	case string:
		return d
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
	case []byte:
		return string(d)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(d)
	default:
		return fmt.Sprint(v)
	}
}

func CheckMobile(mobile string) bool {
	reg := `^1\d{10}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(mobile)
}

// GetRandomBoth 获取随机数  数字和文字
func GetRandomBoth(n int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// GetRandomNum 获取随机数  纯数字
func GetRandomNum(n int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// Sha1En sha1加密
func Sha1En(data string) string {
	t := sha1.New() ///产生一个散列值得方式
	_, _ = io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func Desensitize(str string) string {
	strLen := len(str)
	if strLen < 8 {
		var asterisks string
		for i := 1; i < strLen; i++ {
			asterisks += "*"
		}
		return str[:1] + asterisks
	}

	maskedStr := str[:4]
	for i := 4; i < strLen-4; i++ {
		if len(maskedStr) > 8 {
			break
		}
		maskedStr += "*"
	}
	maskedStr += "****" + str[strLen-4:]

	return maskedStr
}

func GenerateRandomNumber(min, max int) uint32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rInt := r.Intn(max-min+1) + min
	return uint32(rInt)
}

func NowTimeFormat() string {
	return time.Now().Format("2006-01-02")
}

/** 加密方式 **/

func Md5ByString(str string) string {
	m := md5.New()
	io.WriteString(m, str)
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
}

func Md5ByBytes(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

func ConvertToInt64(v interface{}) int64 {
	switch value := v.(type) {
	case string:
		a, _ := strconv.Atoi(value)
		return int64(a)

	}
	return 0
}

func GetGender(g int64) string {
	if g == 1 {
		return "男生"
	} else if g == 2 {
		return "女生"
	}
	return "未知"
}

func CheckNo(s string) bool {
	if s == "N" {
		return true
	}
	return false
}

func CheckYes(s string) bool {
	if s == "Y" || s == "y" {
		return true
	}
	//if strings.Contains(s, "Y") {
	//	return true
	//}
	return false
}

func ContainYes(s string) bool {
	if strings.Contains(s, "Y") {
		return true
	}
	return false
}

func GetAge(birthDate string) int {
	birth, _ := time.Parse("2006-01-02", birthDate)
	return time.Now().Year() - birth.Year()
}

func SplitByN(s string, n int) []string {
	var result []string
	runes := []rune(s) // 将字符串转换为rune切片
	j := 0
	count := 0
	for i := 0; i < len(runes); i++ {
		count++
		if string(runes[i]) == "\n" || count == n || i == len(runes)-1 {
			result = append(result, string(runes[j:i+1]))
			j = i + 1
			count = 0
		}
	}
	return result
}

func CovertMultilineStr(str string, num int) string {
	textLines := SplitByN(str, num)
	newTxtStr := ""
	for _, textLine := range textLines {
		if string(textLine[len(textLine)-1]) == "\n" {
			newTxtStr += textLine
		} else {
			newTxtStr += textLine + "\n"
		}
	}
	return newTxtStr
}
func GetAgeInMonths(birthdayStr string) int {
	birthday, _ := time.Parse("2006-01-02", birthdayStr)

	currentDate := time.Now()
	months := ((currentDate.Year() - birthday.Year()) * 12) + int(currentDate.Month()) - int(birthday.Month())

	if currentDate.Day() < birthday.Day() {
		months--
	}

	return months
}

func ExtractNumberFromStr(str string) int {
	// 尝试将字符串直接转换为数字
	if num, err := strconv.Atoi(str); err == nil {
		if num > 5 {
			return 5
		}
		return num
	}
	// 尝试使用正则表达式提取字符串中的数字
	//先尝试匹配带分的数字
	reScore := regexp.MustCompile(`\d+分`)
	match := reScore.FindString(str)
	if match != "" {
		// 提取分数数字
		reNum := regexp.MustCompile(`\d+`)
		numStr := reNum.FindString(match)
		num, err := strconv.Atoi(numStr)
		if err == nil {
			if num > 5 {
				return 5
			}
			return num
		}
	}
	re := regexp.MustCompile(`\d+`)
	numStr := re.FindString(str)
	if numStr != "" {
		num, err := strconv.Atoi(numStr)
		if err == nil {
			if num > 5 {
				return 5
			}
			return num
		}
	}
	// 默认返回值1
	return 1
}

func CalculateAgeAndMonth(birthdate string) (int, int, error) {
	// 解析出生日期
	layout := "2006-01-02" // Go 的时间格式化布局
	birthTime, err := time.Parse(layout, birthdate)
	if err != nil {
		return 0, 0, err
	}

	// 获取当前时间
	now := time.Now()

	// 计算年份差异
	years := now.Year() - birthTime.Year()

	// 如果今年的生日还没到，则年龄减一
	if now.Month() < birthTime.Month() || (now.Month() == birthTime.Month() && now.Day() < birthTime.Day()) {
		years--
	}

	// 计算月份差异
	months := int(now.Month()) - int(birthTime.Month())
	if now.Day() < birthTime.Day() {
		months--
	}
	if months < 0 {
		months += 12
	}

	return years, months, nil
}

func FilterYN(s string) string {
	s = strings.ReplaceAll(s, "Y", "")
	s = strings.ReplaceAll(s, "N", "")
	return s
}

func GenerateUniqueStrings(count, l int) []string {
	strings := make([]string, count)
	for i := 0; i < count; i++ {
		// 生成16个字节的随机数据
		randomBytes := make([]byte, l)
		_, err := rand.Read(randomBytes)
		if err != nil {
			fmt.Println("Error generating random bytes:", err)
			return nil
		}

		// 将随机数据转换为32位的十六进制字符串
		strings[i] = hex.EncodeToString(randomBytes)
	}
	return strings
}
func CheckStringLength(s string) bool {
	if len(s) > 300 {
		return true
	}
	return false
}

func SqlToTime(now sql.NullTime) (res time.Time) {
	if now.Valid {
		return now.Time
	}
	return
}

func TimeToString(now sql.NullTime) string {
	if now.Valid {
		return TimeFormat(now.Time)
	}
	return ""
}

func StringToInt64(s string) int64 {
	res, _ := strconv.Atoi(s)
	return int64(res)
}

func InArray(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func StructPrompt(prompt, userBackground, chatSummary string) string {
	prompt = strings.Replace(strings.Replace(prompt, "{{user_background}}", userBackground, -1), "{{chat_summary}}", chatSummary, -1)
	return prompt
}

func JsonDecodeString(s sql.NullString) (res []string) {
	json.Unmarshal([]byte(SqlToString(s)), &res)
	return res
}
