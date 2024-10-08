package php2go

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
	"unicode/utf8"
)

func TestTime(t *testing.T) {
	gt(t, float64(Time()), 1522684800)

	// Ensure we take timezones into account
	_, offset := time.Now().Local().Zone()
	unix := int64(1524799394)
	equal(t, "27/04/2018 03:23:14 AM", Date("02/01/2006 15:04:05 PM", unix-int64(offset)))

	tStrtotime, _ := Strtotime("02/01/2006 15:04:05", "02/01/2016 15:04:05")
	equal(t, int64(1451747045), tStrtotime)
	tStrtotime1, _ := Strtotime("3 04 PM", "8 41 PM")
	equal(t, int64(-62167144740), tStrtotime1)

	equal(t, false, Checkdate(2, 29, 2018))
	equal(t, true, Checkdate(2, 29, 2020))
}

func TestString(t *testing.T) {
	equal(t, `\'wo\'中文\"emoji`, Addslashes("'wo'中文\"emoji"))

	equal(t, "e10adc3949ba59abbe56e057f20f883e", Md5("123456"))

	equal(t, "7c4a8d09ca3762af61e59520943dc26494f8941b", Sha1("123456"))

	equal(t, uint32(158520161), Crc32("123456"))

	equal(t, "中文中文中文", StrRepeat("中文", 3))

	equal(t, "ab", Substr("abc", 0, 2))

	equal(t, "@gmail.com", Strstr("xxx@gmail.com", "@"))

	equal(t, "Hello world", Ucfirst("hello world"))

	equal(t, "hello world", Lcfirst("Hello world"))

	equal(t, "Hello World", Ucwords("hello world"))

	tStrlen := Strlen("G中文")
	equal(t, 7, tStrlen)

	tMbStrlen := MbStrlen("G中文")
	equal(t, 3, tMbStrlen)

	equal(t, 6, Strpos("hello wworld", "w", -6))

	equal(t, -1, Stripos("hello Wworld", "w", 8))

	equal(t, 6, Strrpos("hello wworld", "w", -6))

	equal(t, 7, Strripos("hello wWorld", "w", 0))

	equal(t, "a,b,c", Implode(",", []string{"a", "b", "c"}))

	tAddslashes := Addslashes("f'oo b\"ar")
	equal(t, `f\'oo b\"ar`, tAddslashes)

	equal(t, `f'oo b"ar\a\\\`, Stripslashes("f\\'oo b\\\"ar\\\\a\\\\\\\\\\\\"))

	equal(t, 4, Levenshtein("golang", "google", 1, 1, 1))

	var percent float64
	equal(t, 3, SimilarText("golang", "google", &percent))
	equal(t, float64(50), percent)

	equal(t, "H416", Soundex("Heilbronn"))

	equal(t, 14, len(Uniqid("x")))
	equal(t, true, bytes.HasPrefix([]byte(Uniqid("x")), []byte("x")))

	equal(t, 5, utf8.RuneCountInString(StrShuffle("中˚abc")))

	equal(t, 20013, Ord("中文"))

	equal(t, "z", Chr(122))

	equal(t, 4, MbStrlen("中文 a"))

	equal(t, "<a><br />xxx<br />yy<br />中文<br />n<br />x", Nl2br("<a>\n\rxxx\nyy\r中文\r\nn\r\nx", true))

	equal(t, "---)(5.01文中% cin 	 cba", Strrev("abc \t nic %中文10.5()---"))

	equal(t, "abce 	 enice %中e文10e.5(e)--e-e", ChunkSplit("abc \t nic %中文10.5()---", 3, "e"))

	equal(t, `\.\+\?\[\$\]\(\*\)\^中文`, Quotemeta(".+?[$](*)^中文"))

	equal(t, `&lt;html&gt;hello world &lt;/html&gt;`, Htmlentities("<html>hello world </html>"))

	equal(t, "<html>hello world </html>", HTMLEntityDecode("&lt;html&gt;hello world &lt;/html&gt;"))

	equal(t, "abc\nhello\nworld\nxxx", Wordwrap("abc hello world xxx", 5, "\n", false))

	equal(t, []string{"a", "b", "c"}, StrWordCount("a b c"))

	equal(t, "1001", Strtr("baab", "ab", "01"))
	equal(t, "bccb", Strtr("baab", "ab", "c"))
	equal(t, "bccb", Strtr("baab", "a", "cd"))
	equal(t, "ba01", Strtr("baab", map[string]string{"ab": "01"}))

	tParseStr := make(map[string]interface{})
	_ = ParseStr("f[a][]=m&f[a][]=n", tParseStr)
	equal(t, "map[f:map[a:[m n]]]", fmt.Sprint(tParseStr))
}

func TestArray(t *testing.T) {
	var s1 = make([]interface{}, 3)
	s1[0] = "a"
	s1[1] = "b"
	s1[2] = "c"
	tArrayChunk := ArrayChunk(s1, 2)
	equal(t, [][]interface{}{{"a", "b"}, {"c"}}, tArrayChunk)

	var au = []string{"a", "a"}
	tau := ArrayUnique(au)
	equal(t, []string{"a"}, tau)

	var aui = []int{1, 1}
	taui := ArrayUniqueInt(aui)
	equal(t, []int{1}, taui)

	var ad1 = []string{"a", "c"}
	var ad2 = []string{"a", "b"}
	tad := ArrayDiff(ad1, ad2)
	equal(t, []string{"c", "b"}, tad)

	var m1 = make(map[interface{}]interface{}, 3)
	m1[1] = "a"
	m1["a"] = "b"
	m1[2.5] = 1
	tArrayKeyExists := ArrayKeyExists("a", m1)
	equal(t, true, tArrayKeyExists)

	tArrayUnshift := ArrayUnshift(&s1, "x", "y")
	equal(t, 5, tArrayUnshift)
	equal(t, []interface{}{"x", "y", "a", "b", "c"}, s1)

	equal(t, 7, ArrayPush(&s1, "u", "v"))
	equal(t, []interface{}{"x", "y", "a", "b", "c", "u", "v"}, s1)

	equal(t, "v", ArrayPop(&s1))
	equal(t, []interface{}{"x", "y", "a", "b", "c", "u"}, s1)

	tArrayShift := ArrayShift(&s1)
	equal(t, "x", tArrayShift)
	equal(t, []interface{}{"y", "a", "b", "c", "u"}, s1)

	equal(t, map[int]interface{}{-1: "aa", 0: "aa", 1: "aa", 2: "aa", -3: "aa", -2: "aa"}, ArrayFill(-3, 6, "aa"))

	equal(t, 3, len(ArrayRand([]interface{}{"a", "b", "c"})))

	var s2 = make([]interface{}, 3)
	s2[0] = "a"
	s2[1] = "b"
	s2[2] = "c"
	equal(t, []interface{}{"d", "d", "a", "b", "c"}, ArrayPad(s2, -5, "d"))

	var s3 = make([]interface{}, 3, 3)
	s3[0] = "x"
	s3[1] = "y"
	s3[2] = "z"
	equal(t, map[interface{}]interface{}{"a": "x", "b": "y", "c": "z"}, ArrayCombine(s2, s3))

	tInArray1 := InArray(1, [2]interface{}{"a", 1})                        // array
	tInArray2 := InArray(1, []interface{}{"a", 1})                         // slice
	tInArray3 := InArray(1, map[interface{}]interface{}{"a": "c", 1: "d"}) // map
	equal(t, true, tInArray1)
	equal(t, true, tInArray2)
	equal(t, false, tInArray3)
}

func TestUrl(t *testing.T) {
	tParseURL, _ := ParseURL("http://username:password@hostname:9090/path?arg=value#anchor", -1)
	equal(t, map[string]string{"pass": "password", "path": "/path", "query": "arg=value", "fragment": "anchor", "scheme": "http", "host": "hostname", "port": "9090", "user": "username"}, tParseURL)

	equal(t, "http%3A%2F%2Fgolang.org%3Fx+y", URLEncode("http://golang.org?x y"))

	tURLDecode, _ := URLDecode("http%3A%2F%2Fgolang.org%3Fx+y")
	equal(t, "http://golang.org?x y", tURLDecode)

	equal(t, "http%3A%2F%2Fgolang.org%3Fx%20y", Rawurlencode("http://golang.org?x y"))

	tRawurldecode, _ := Rawurldecode("http%3A%2F%2Fgolang.org%3Fx%20y")
	equal(t, "http://golang.org?x y", tRawurldecode)

	equal(t, "VGhpcyBpcyBhbiBlbmNvZGVkIHN0cmluZw==", Base64Encode("This is an encoded string"))

	tBase64Decode, _ := Base64Decode("VGhpcyBpcyBhbiBlbmNvZGVkIHN0cmluZw")
	equal(t, "This is an encoded string", tBase64Decode)

	equal(t, "first=value&multi=foo+bar&multi=baz", HTTPBuildQuery(map[string][]string{"first": {"value"}, "multi": {"foo bar", "baz"}}))
}

func TestMath(t *testing.T) {
	equal(t, float64(5), Max(2, 3.7, 5, 1.1))

	equal(t, 1.1, Min(2, 3.7, 5, 1.1))

	rangeValue(t, float64(2), float64(5), float64(Rand(2, 5)))

	tDecbin := Decbin(100)
	equal(t, "1100100", tDecbin)

	tBindec, _ := Bindec(tDecbin)
	equal(t, int64(100), tBindec)

	tBinhex, _ := Binhex(tDecbin)
	equal(t, "64", tBinhex)

	tHexdec, _ := Hexdec(tBinhex)
	equal(t, int64(100), tHexdec)

	tDechex := Dechex(tHexdec)
	equal(t, "64", tDechex)

	tHexbin, _ := Hexbin(tDechex)
	equal(t, "1100100", tHexbin)

	tBin2hex, _ := Bin2hex("你好世界")
	equal(t, "e4bda0e5a5bde4b896e7958c", tBin2hex)

	tHex2binStr, _ := Hex2bin(tBin2hex)
	equal(t, "你好世界", tHex2binStr)

	tDecoct := Decoct(tHexdec)
	equal(t, "144", tDecoct)

	tOctdec, _ := Octdec(tDecoct)
	equal(t, int64(100), tOctdec)

	tBaseConvert, _ := BaseConvert("64", 16, 2)
	equal(t, "1100100", tBaseConvert)

	equal(t, "1,234,567,890.78", NumberFormat(1234567890.777, 2, ".", ","))
}

func TestFile(t *testing.T) {
	tRealpath1, _ := Realpath("/home/go/../go/test/../")
	equal(t, "/home/go", tRealpath1)

	equal(t, "php2go.go", Basename("/home/go/src/pkg/php2go.go"))

	tPathinfo := Pathinfo("/home/go/php2go.go.go", -1)
	equal(t, map[string]string{"dirname": "/home/go", "basename": "php2go.go.go", "extension": "go", "filename": "php2go.go"}, tPathinfo)

	wd, _ := os.Getwd()
	tFilesize, _ := FileSize(wd)
	gt(t, float64(tFilesize), 0)
}

func TestVariable(t *testing.T) {
	equal(t, true, IsNumeric("-0xaF"))
	equal(t, true, IsNumeric("123456"))

	equal(t, true, Empty(nil))
	equal(t, true, Empty(false))
	equal(t, true, Empty(0))
	equal(t, true, Empty(""))
	equal(t, true, Empty(0.0))
	equal(t, true, Empty([]int{}))
	equal(t, true, Empty([0]int{}))
	equal(t, false, Empty([1]int{}))
	equal(t, true, Empty(map[int]int{}))
}

func TestProgramExecution(t *testing.T) {
	var output []string
	var retVal int
	gt(t, float64(len(Exec("/bin/bash -c \"ls -a|grep php\"", &output, &retVal))), 0)
	equal(t, 0, retVal)

	equal(t, 0, retVal)
	gt(t, float64(len(System("ls -l", &retVal))), 0)

	Passthru("echo hello", &retVal)
	equal(t, 0, retVal)
}

func TestNetwork(t *testing.T) {
	tGethostname, _ := Gethostname()
	gt(t, float64(len(tGethostname)), 0)

	equal(t, uint32(134744072), IP2long("8.8.8.8"))

	equal(t, "8.8.8.8", Long2ip(134744072))

	tGethostbyname, _ := Gethostbyname("localhost")
	equal(t, "127.0.0.1", tGethostbyname)

	tGethostbynamel, _ := Gethostbynamel("localhost")
	gt(t, float64(len(tGethostbynamel)), 0)

	tGethostbyaddr, _ := Gethostbyaddr("127.0.0.1")
	equal(t, "localhost", tGethostbyaddr)
}

func TestMisc(t *testing.T) {
	equal(t, true, VersionCompare("1.3-beta", "1.4Rc1", "<"))

	gt(t, float64(MemoryGetUsage(true)), 0)

	tPack := Pack([]string{"N2", "N4"}, 123, 45678)
	equal(t, "7b006eb20000", tPack)

	tUnpack := Unpack([]string{"N2", "N4"}, tPack)
	equal(t, int64(123), tUnpack[0])
	equal(t, int64(45678), tUnpack[1])

	var i interface{}
	i = "123"
	ti := GetInterfaceToInt(i)
	equal(t, 123, ti)

	i = "123.0"
	tf := GetInterfaceToFloat(i)
	equal(t, 123.0, tf)

	i = 123.0
	ts := GetInterfaceToString(i)
	equal(t, "123", ts)
}

func BenchmarkFn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ChunkSplit("abc", 2, "\r\n")
	}
}

func TestPack(t *testing.T) {
	i := int64(4000000000)
	hexEncode := HexEncode(i, 62)
	equal(t, "4mHAJ2", hexEncode)
	hexDecode := HexDecode(hexEncode, 62)
	equal(t, i, hexDecode)

	var i1 uint8 = 100
	var i2 uint32 = 1000000000
	var i3 uint64 = 100000000000000000
	p := &Protocol{}
	p.Format = []string{"N2", "N4", "N8"}
	id := p.Pack16(
		int64(i1),
		int64(i2),
		int64(i3),
	)
	log.Println(id)
	p2 := &Protocol{}
	p2.Format = p.Format
	s := p2.UnPack16(id)
	log.Println(s)
	equal(t, i1, uint8(s[0]))
	equal(t, i2, uint32(s[1]))
	equal(t, i3, uint64(s[2]))
}

func TestConnection(t *testing.T) {
	NewServer("127.0.0.1", 80)
	conn := NewConnection("127.0.0.2", 443)
	clientId := conn.Pack()
	log.Println("clientId", clientId)
	unConn := Connection{}
	unConn.UnPack(clientId)
	equal(t, "127.0.0.1:80", unConn.Server.GetAddress())
	equal(t, "127.0.0.2:443", unConn.Client.GetAddress())

	log.Println("server ip:", GetInBoundIP())
	log.Println("local ip", GetOutBoundIP())
}

func TestConsulApi_WatchKeyToPath(t *testing.T) {
	consul, err := NewConsul("", "")
	if err != nil {
		panic("consul connection error:" + err.Error())
	}
	//构造数据
	file := "test/sub/sub2/file.txt"
	value := []byte("test")
	var bt []byte
	consul.GetClient().KV().Put(&api.KVPair{Key: file, Value: value}, nil)

	go consul.WatchKeyToPath(file, "test", "key")
	time.Sleep(1 * time.Second)
	bt, _ = os.ReadFile("test/file.txt")
	equal(t, "test", string(bt))

	go consul.WatchKeyToPath(file, "test/file.txt", "file")
	time.Sleep(1 * time.Second)
	bt, _ = os.ReadFile("test/file.txt")
	equal(t, value, bt)
	log.Println(bt, 1111)

	go consul.WatchKeyToPath(file, "test/file.txt", "source")
	time.Sleep(1 * time.Second)
	bt, _ = os.ReadFile(file)
	equal(t, value, bt)

	go consul.WatchKeyToPath("test", "test", "")
	time.Sleep(1 * time.Second)
	bt, _ = os.ReadFile(file)
	equal(t, value, bt)

	go consul.WatchKeyToPath("test", "test/new/test", "")
	time.Sleep(1 * time.Second)
	bt, _ = os.ReadFile("test/new/" + file)
	equal(t, value, bt)

	//清理数据
	os.RemoveAll("test")
	consul.GetClient().KV().DeleteTree("test", nil)
	//ch := make(chan int, 1)
	//<-ch
	//log.Println(1111)
}

// Expected to be equal.
func equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
	}
}

// Expected to be unequal.
func unequal(t *testing.T, expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
	}
}

// Expect a greater than b.
func gt(t *testing.T, a, b float64) {
	if a <= b {
		t.Errorf("Expected %v (type %v) > Got %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
	}
}

// Expect a greater than or equal to b.
func gte(t *testing.T, a, b float64) {
	if a < b {
		t.Errorf("Expected %v (type %v) > Got %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
	}
}

// Expected value needs to be within range.
func rangeValue(t *testing.T, min, max, actual float64) {
	if actual < min || actual > max {
		t.Errorf("Expected range of %v-%v (type %v) > Got %v (type %v)", min, max, reflect.TypeOf(min), actual, reflect.TypeOf(actual))
	}
}
