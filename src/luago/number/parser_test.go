package number

import "fmt"
import "testing"

func TestParseInteger(t *testing.T) {
	testInteger(t, "+100", 100)
	testInteger(t, "-100", -100)
	testInteger(t, "0xffffffffffffffff", -1)
	testInteger(t, "0xfffffffffffffffe", -2)
	testInteger(t, "-0XFFFFFFFFFFFFFFFE", 2)
	// testInteger(t, "0xffffffffffffffff.0", 0)
}

func TestParseFloat(t *testing.T) {
	testFloat(t, "314.16e-2", "3.1416")
	testFloat(t, "0.31416E1", "3.1416")
	testFloat(t, "34e1", "340")
	testFloat(t, "0x0.1E", "0.1171875")
	testFloat(t, "0xA23p-4", "162.1875")
	testFloat(t, "0X1.921FB54442D18P+1", "3.141592653589793")
	testFloat(t, "0xffffffffffffffff.0", "1.8446744073709552e+19")
	testFloat(t, "3.", "3")
	testFloat(t, ".3", "0.3")
	testFloat(t, "3.e1", "30")
	testFloat(t, ".3e1", "3")
	testFloat(t, "0x1.p1", "2")
	testFloat(t, "0x.1p1", "0.125")
	testFloat(t, "+0x2", "2")
	testFloatFail(t, "0x")
	testFloatFail(t, "-0x")
	//testFloat(t, "0x.p1", "?") todo
}

func testInteger(t *testing.T, str string, i int64) {
	j, ok := ParseInteger(str)
	if !ok || j != i {
		t.Errorf("%d != %d", j, i)
	}
}

func testFloat(t *testing.T, str string, x string) {
	f, ok := ParseFloat(str)
	if !ok {
		t.Fail()
	}
	y := fmt.Sprintf("%g", f)
	if y != x {
		t.Errorf("%s != %s", y, x)
	}
}

func testFloatFail(t *testing.T, str string) {
	_, ok := ParseFloat(str)
	if ok {
		t.Fail()
	}
}
