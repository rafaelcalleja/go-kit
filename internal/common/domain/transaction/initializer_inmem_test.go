package transaction

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"unsafe"
)

func TestInMemTransaction_Rollback(t *testing.T) {
	ctx := context.Background()
	demo := map[string]string{
		"foo": "bar",
	}
	//actual := make(map[string]string)
	v := reflect.ValueOf(demo)
	mapType := reflect.MapOf(v.Type().Elem(), v.Type().Key())
	actual := reflect.MakeMap(mapType).Interface()

	expected := make(map[string]string)

	ctx = context.WithValue(ctx, "repository", NewCloneable(&demo, &actual))

	transaction := NewInitializerInMem(ctx)
	tx, err := transaction.Begin()
	demo["foo"] = "foo"

	require.NoError(t, err)

	err = tx.Rollback()
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}

func TestInMemTransaction_Success(t *testing.T) {
	ctx := context.Background()
	demo := map[string]string{
		"foo": "bar",
	}

	v := reflect.ValueOf(demo)
	mapType := reflect.MapOf(v.Type().Elem(), v.Type().Key())
	actual := reflect.MakeMap(mapType).Interface()

	expected := map[string]string{
		"foo": "foo",
	}

	ctx = context.WithValue(ctx, "repository", NewCloneable(&demo, &actual))

	transaction := NewInitializerInMem(ctx)
	tx, err := transaction.Begin()
	demo["foo"] = "foo"

	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestInMemTransaction_Commit(t *testing.T) {
	//ctx := context.Background()
	convert := func(pointer unsafe.Pointer, value interface{}) interface{} {
		switch fmt.Sprintf("%T", value) {
		case "string":
			return *(*string)(pointer)
		}

		return *(*string)(pointer)
	}

	convertA := func(pointer AnyPointer, value interface{}) interface{} {
		switch fmt.Sprintf("%T", value) {
		case "string":
			return *(*string)(pointer)
		}

		return *(*string)(pointer)
	}

	update := func(pointer unsafe.Pointer, value interface{}) {
		newP := uintptr(unsafe.Pointer(&value))
		switch fmt.Sprintf("%T", value) {
		case "string":
			*(*uintptr)(pointer) = newP
		}

		fmt.Printf("%v\n", *(*string)(pointer))
	}

	myValue := "foo"

	var p1 = unsafe.Pointer(&myValue)
	fmt.Printf("c %T %v\n", convert(p1, myValue), convert(p1, myValue))

	update(p1, "bar")

	fmt.Printf("v %s %s\n", myValue, *(*string)(p1))

	fmt.Printf("c %T %v\n", convertA(AnyPointer(p1), myValue), convertA(AnyPointer(p1), myValue))

	fmt.Printf("p1 %T %v\n", &p1, &p1)
	fmt.Printf("v %T %v\n", &myValue, &myValue)

	fmt.Printf("p1 %s\n", *(*string)(p1))
	myValue = "bar"
	fmt.Printf("p1 %s\n", *(*string)(p1))

	/*expected := NewCloneableFunc(func() *interface{} {
		var i interface{}
		i = myValue

		var p1 = unsafe.Pointer(&myValue)
		var p2 = uintptr(unsafe.Pointer(&i))

		fmt.Printf("%T %v\n", &p1, &p1)
		fmt.Printf("%T %v\n", &p2, &p2)

		fmt.Printf("\n%s\n", *(*string)(p1))

		return &i
	})

	ctx = context.WithValue(ctx, "repository", expected)
	fmt.Printf("%T %v\n", &myValue, &myValue)

	transaction := NewInitializerInMem(ctx)
	tx, err := transaction.Begin()
	myValue = "bar"
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	actual := ctx.Value("repository")
	require.Equal(t, "bar", actual)*/
}

func TestOther(t *testing.T) {

	val := "foo"
	immutable := val
	p1 := &val
	*p1 = "bar"

	require.Equal(t, "bar", val)

	val = "foo2"
	require.Equal(t, "foo2", *p1)
	require.Equal(t, &val, p1)

	require.Equal(t, "foo", immutable)

	val = immutable

	require.Equal(t, val, immutable)
	require.Equal(t, &val, p1)

	pointer := unsafe.Pointer(&val)
	s := *(*string)(pointer)
	fmt.Printf("%s\n", s)

	m := new(fns)
	m.SaveFn = func() unsafe.Pointer {
		ptr := unsafe.Pointer(&val)

		defer func() {
			other := *(*string)(ptr)
			m.orig = unsafe.Pointer(&other)
			fmt.Println("ORIG v1 ", *(*string)(m.orig))
		}()

		return ptr
	}

	m.RecoverFn = func(i *interface{}) {
		fmt.Println(*(*string)(m.Get()))
		fmt.Println("ORIG v2 ", *(*string)(m.orig))
		val = *(*string)(m.orig)
	}

	m.Save()
	v2 := *(*string)(m.Get())
	require.Equal(t, v2, val)

	val = "test"
	require.Equal(t, "foo", v2)

	v2 = *(*string)(m.Get())
	require.Equal(t, "test", v2)

	m.Recover(val)
	require.Equal(t, "foo", val)

	p := SaveString(val)
	s2 := GetString(p)
	fmt.Printf("%s\n", s2)

	require.Equal(t, val, s2)
	val = "bar2"
	require.Equal(t, GetString(p), s2)
	val = GetString(p)

	require.Equal(t, val, s2)

	//fmt.Printf("%s\n", Convert(&val))
	//fmt.Printf("%s\n", GetString(saveAny))
	//getValue := GetString(saveAny)

	//require.Equal(t, val, getValue)
}

func createInt() *int {
	return new(int)
}

func TestUnsafe(t *testing.T) {
	p0, y, z := createInt(), createInt(), createInt()
	var p1 = unsafe.Pointer(y)
	var p2 = uintptr(unsafe.Pointer(z))

	x := 4
	pointer := unsafe.Pointer(&x)
	var p3 = uintptr(pointer)
	var p4 = (*int)(unsafe.Pointer(uintptr(pointer)))

	var p5 = unsafe.Pointer(&p3)
	x2 := *p4

	fmt.Println(*p4, x2, (*int)(p5), *(*int)(p5))

	*p0 = 1
	*(*int)(p1) = 2
	*(*int)(unsafe.Pointer(p2)) = 3

	*(*int)(unsafe.Pointer(p3)) = 3

	fmt.Println(p0, *p0, p1, *(*int)(p1), p2, *(*int)(unsafe.Pointer(p2)))

	x = 9
	p6 := uintptr(pointer)
	fmt.Println(x, *(*int)(unsafe.Pointer(uintptr(p6))))
	fmt.Println(&x, (*int)(unsafe.Pointer(uintptr(p6))))

	*(*int)(unsafe.Pointer(uintptr(pointer))) = 5

	fmt.Println(x, *(*int)(unsafe.Pointer(uintptr(p6))))

}

type MyInt int

func TestCasting(t *testing.T) {
	x := 111
	pointer := unsafe.Pointer(&x)
	up1 := uintptr(pointer)

	fmt.Println(*(*int)(unsafe.Pointer(uintptr(up1))))
	fmt.Println(*(*MyInt)(unsafe.Pointer(uintptr(up1))))
	fmt.Println(*(*MyInt)(unsafe.Pointer(&x)))

	pps := new(pointerS)

	modifyByVal(unsafe.Pointer(&x), pps)

	//fmt.Println(x, &x)
	//require.Equal(t, 5, x)
}

type pointerS struct {
	src  unsafe.Pointer
	copy unsafe.Pointer
}

func (p pointerS) Get() []byte {
	i := *(*int)(p.src)
	arr := []byte(strconv.Itoa(i))
	size := len(arr)
	p1 := uintptr(unsafe.Pointer(&arr))

	var data []byte

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = p1
	sh.Len = size
	sh.Cap = size

	//hdr := (*reflect.StringHeader)(p.src)
	//pbyte := (*byte)(unsafe.Add(unsafe.Pointer(hdr.Data), 2))
	runtime.KeepAlive(arr)

	fmt.Println(data)
	return data
}

func modifyByVal(src unsafe.Pointer, pps *pointerS) {
	pps.src = src
	var data []byte = pps.Get()
	fmt.Println(data)

	//memAddress := uintptr(pointer)

	//pointer := unsafe.Pointer(value)
	//up1 := uintptr(pointer)
	//fmt.Println(*(*int)(unsafe.Pointer(uintptr(up1))))
	//fmt.Println(&value, (*int)(unsafe.Pointer(uintptr(up1))))
	//fmt.Println(&value, &pointer)

	//*(*int)(unsafe.Pointer(uintptr(pointer))) = 6

	//*(*int)(pointer) = 5
}

func unsafeGetBytesWRONG(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s)) // WRONG!!!!
}

func unsafeGetBytes(s interface{}) []byte {
	l := len(s.(string))
	fmt.Println(l, s.(string), unsafe.Sizeof(s), unsafe.Alignof(s))
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data),
	)[:l:l]
}
func unsafeGetBytesA(s string) []byte {
	fmt.Println(len(s), s, unsafe.Sizeof(s), unsafe.Alignof(s))
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data),
	)[:len(s):len(s)]
}

func ByteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func TestSlicing(t *testing.T) {
	c := "hello world"

	fmt.Println(unsafeGetBytes(c))
	fmt.Println(unsafeGetBytesA(c))
	fmt.Println(unsafeGetBytesWRONG(c))
	fmt.Println(ByteSlice2String(unsafeGetBytes(c)))
	fmt.Println(ByteSlice2String(unsafeGetBytesA(c)))
	fmt.Println(ByteSlice2String(unsafeGetBytesWRONG(c)))
}

func TestByting(t *testing.T) {
	a := [...]byte{'G', 'o', 'l', 'a', 'n', 'g'}
	fmt.Println(a)
	s := "Java"
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = uintptr(unsafe.Pointer(&a))
	hdr.Len = len(a)
	fmt.Println(s) // Golang
	// Now s and a share the same byte sequence, which
	// makes the bytes in the string s become mutable.
	a[2], a[3], a[4], a[5] = 'o', 'g', 'l', 'e'
	fmt.Println(s) // Google
}

func TestCapsule(t *testing.T) {
	s := "Java"
	fmt.Println(&s)
	capsule(unsafe.Pointer(&s))
	fmt.Println(s)
}

type News struct {
	A int
	b string
}

type SameSize struct {
	A int
	c string
	b bool
	d *int
}

func TestCapsuleA(t *testing.T) {

	n := new(News)
	n.A = 1
	n.b = "b"

	fmt.Println(n)
	bytes := []byte(fmt.Sprintf("%v", n))
	fmt.Println(bytes)
	newsR := ByteSlice2News(n, n)

	fmt.Println(newsR)
	fmt.Println([]byte(fmt.Sprintf("%v", newsR)))
	//capsule(unsafe.Pointer(&n))
	//fmt.Println(n)
}

func TestCoping(t *testing.T) {
	f := 2
	n := new(SameSize)
	n.A = 1
	n.c = "b"
	n.b = false
	n.d = &f

	s := NewSaveStruct(n)

	s.Copy()
	n.A = 2
	require.Equal(t, "b", n.c)
	require.Equal(t, n, (*SameSize)(s.save))
	require.NotEqual(t, n, (*SameSize)(s.copy))

	s.Recover()
	fmt.Println(n)
	require.Equal(t, 1, n.A)
	require.Equal(t, "b", n.c)
	require.Equal(t, n, (*SameSize)(s.save))
	//	require.Equal(t, *(*SameSize)(unsafe.Pointer(&s.raw[0])), *(*SameSize)(s.save))
	require.Equal(t, n, (*SameSize)(s.copy))
	require.Equal(t, n.A, (*SameSize)(s.copy).A)

	fmt.Println(*(*SameSize)(s.copy).d)

}

type saveSameSize struct {
	save         unsafe.Pointer
	copy         unsafe.Pointer
	sizeOfStruct uintptr
	underline    reflect.Value
	rtype        reflect.Type
	raw          []byte
}

func (s *saveSameSize) Copy() {
	fmt.Println("src is ", &s.save, (*SameSize)(s.save))
	fmt.Println("copy is ", &s.copy, (*SameSize)(s.copy))
	var dest SameSize

	// Explicitly copy the contents of in into out by casting both into byte
	// arrays and then slicing the arrays. This will produce the correct packed
	// union structure, without relying on unsafe casting to a smaller type of a
	// larger type.
	const sizeof = unsafe.Sizeof(dest)

	copy(
		(*(*[sizeof]byte)(unsafe.Pointer(&dest)))[:],
		(*(*[sizeof]byte)(s.save))[:],
	)

	s.copy = unsafe.Pointer(&dest)
	fmt.Println("src is ", &s.save, (*SameSize)(s.save))
	fmt.Println("copy is ", &s.copy, (*SameSize)(s.copy))
}

func (s *saveSameSize) Recover() {
	var src = s.getDataFromPointer(s.copy)
	*(*SameSize)(s.save) = *(*SameSize)(unsafe.Pointer(&src[0]))
	fmt.Println("to recover", *(*SameSize)(s.save))
	/*data := (*(*[1<<31 - 1]byte)(s.copy))[:s.sizeOfStruct]
	*(*SameSize)(s.save) = *(*SameSize)(unsafe.Pointer(&data[0]))*/
	//*(*SameSize)(s.save) = *(*SameSize)(src)
	//s.raw = (*(*[1<<31 - 1]byte)(src))[:s.sizeOfStruct]
}

func saveRecover(s saveSameSize) {
	var dest SameSize
	var src = s.save

	fmt.Println(s.save)

	copy(
		(*(*[unsafe.Sizeof(dest)]byte)(unsafe.Pointer(&dest)))[:],
		(*(*[unsafe.Sizeof(dest)]byte)(unsafe.Pointer(&src)))[:],
	)

	fmt.Println(&dest)
}

func (s saveSameSize) getDataFromPointer(ptr unsafe.Pointer) []byte {
	data := (*(*[1<<31 - 1]byte)(ptr))[:s.sizeOfStruct]
	return data
	//return unsafe.Pointer(&data[0])
}

func NewSaveStruct(toSave *SameSize) saveSameSize {
	s := saveSameSize{
		underline:    reflect.ValueOf(toSave),
		save:         unsafe.Pointer(toSave),
		sizeOfStruct: unsafe.Sizeof(reflect.TypeOf(toSave)),
		rtype:        reflect.TypeOf(toSave),
	}

	fmt.Println("Tosave ", &*toSave, s.save, s.getDataFromPointer(s.save), *(*SameSize)(s.save))

	return s
}

func ByteSlice2News(sizeSrc interface{}, sizeDest interface{}) *News {
	n := new(SameSize)
	n.A = 1
	n.c = "b"
	n.b = false

	var dest News

	// Explicitly copy the contents of in into out by casting both into byte
	// arrays and then slicing the arrays. This will produce the correct packed
	// union structure, without relying on unsafe casting to a smaller type of a
	// larger type.
	copy(
		(*(*[unsafe.Sizeof(dest)]byte)(unsafe.Pointer(&dest)))[:],
		(*(*[unsafe.Sizeof(dest)]byte)(unsafe.Pointer(n)))[:],
	)

	return &dest

	//n := new(News)
	//size := unsafe.Sizeof(n)
	/*hdr := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	hdr.Data = uintptr(unsafe.Pointer(&a))
	hdr.Len = len(a)
	return *(*News)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&bs)).Data),
	)*/
}

func capsule(ptr unsafe.Pointer) {
	fmt.Println(*(*[]byte)(ptr))

	s := *(*string)(ptr)
	a := [...]byte{'G', 'o', 'l', 'a', 'n', 'g'}
	fmt.Println(a)
	hdr := (*reflect.StringHeader)(ptr)
	hdr.Data = uintptr(unsafe.Pointer(&a))
	hdr.Len = len(a)
	fmt.Println(s, ptr) // Golang
	// Now s and a share the same byte sequence, which
	// makes the bytes in the string s become mutable.
	a[2], a[3], a[4], a[5] = 'o', 'g', 'l', 'e'
	fmt.Println(s, ptr) // Google
}

func TestUnmars(t *testing.T) {
	b := "value"
	fmt.Println("b ", &b)
	s := new(savedata)

	acceptAny(&b, s)
	b = "other"
	recoverAny(s)

	require.Equal(t, "value", b)
	/*out := newTwo(one{v: 0xff})

	err := json.Unmarshal(out.b[:], &a)
	require.NoError(t, err)*/

}

type savedata struct {
	save    interface{}
	recover interface{}
}

func recoverAny(s *savedata) {
	s.recover = reflect.ValueOf(s.save)
	ri := reflect.ValueOf(s.save)
	ri.Elem().SetString(s.save.(string))
}

func acceptAny(i interface{}, s *savedata) {
	fmt.Println("acceptAnyI ", i, &i)
	s.save = i

	//ri := reflect.ValueOf(i)
	//var cache = i

	//ri.Set(reflect.ValueOf("newValue"))
}

func TestByte(t *testing.T) {
	n := 4

	// Create a slice of the correct size
	m := make([]int, n)

	// Use convoluted indirection to cast the first few bytes of the slice
	// to an unsafe uintptr
	mPtr := *(*uintptr)(unsafe.Pointer(&m))
	//mPtr := *(**int)(unsafe.Pointer(&m))

	// Check it worked
	m[0] = 987
	// (we have to recast the uintptr to a *int to examine it)
	fmt.Println(m, *(*int)(unsafe.Pointer(mPtr)))
}

// one is a typed Go structure containing structured data to pass to the kernel.
type one struct{ v uint64 }

// two mimics a C union type which passes a blob of data to the kernel.
type two struct{ b [32]byte }

// newTwo safely produces a two structure from an input one.
func newTwo(in one) *two {
	// Initialize out and its array.
	var out two

	// Explicitly copy the contents of in into out by casting both into byte
	// arrays and then slicing the arrays. This will produce the correct packed
	// union structure, without relying on unsafe casting to a smaller type of a
	// larger type.
	copy(
		(*(*[unsafe.Sizeof(two{})]byte)(unsafe.Pointer(&out)))[:],
		(*(*[unsafe.Sizeof(one{})]byte)(unsafe.Pointer(&in)))[:],
	)

	return &out
}

func TestCByte(t *testing.T) {
	// All is well! The two structure is appropriately initialized.
	out := newTwo(one{v: 0xff})

	fmt.Printf("%#v\n", out.b[:8])
	fmt.Println(unsafe.Sizeof(one{}))
}

type fns struct {
	orig      unsafe.Pointer
	ptr       unsafe.Pointer
	SaveFn    func() unsafe.Pointer
	GetFn     func()
	RecoverFn func(i *interface{})
}

func (f *fns) Recover(i interface{}) {
	f.RecoverFn(&i)
}

func (f *fns) Save() {
	f.ptr = f.SaveFn()
}

func (f fns) Get() unsafe.Pointer {
	return f.ptr
}

type Cache interface {
	Save()
	Get()
	Recover()
}

func SaveString(val string) unsafe.Pointer {
	return unsafe.Pointer(&val)
}

func GetString(unsafePointer unsafe.Pointer) string {
	return *(*string)(unsafePointer)
}

func Convert(value *interface{}) string {
	pointer := unsafe.Pointer(value)
	return *(*string)(pointer)
}
