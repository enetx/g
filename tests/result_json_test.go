package g_test

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"errors"
	"fmt"
	"strings"
	"testing"

	. "github.com/enetx/g"
)

func TestResultMarshalJSON_OkInt(t *testing.T) {
	data, err := json.Marshal(Ok(42))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"ok":42}` {
		t.Errorf("got %s, want {\"ok\":42}", data)
	}
}

func TestResultMarshalJSON_OkString(t *testing.T) {
	data, err := json.Marshal(Ok("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"ok":"hello"}` {
		t.Errorf("got %s, want {\"ok\":\"hello\"}", data)
	}
}

func TestResultMarshalJSON_OkStruct(t *testing.T) {
	type Point struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	data, err := json.Marshal(Ok(Point{X: 1, Y: 2}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"ok":{"x":1,"y":2}}` {
		t.Errorf("got %s, want {\"ok\":{\"x\":1,\"y\":2}}", data)
	}
}

func TestResultMarshalJSON_OkSlice(t *testing.T) {
	data, err := json.Marshal(Ok([]int{1, 2, 3}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"ok":[1,2,3]}` {
		t.Errorf("got %s, want {\"ok\":[1,2,3]}", data)
	}
}

func TestResultMarshalJSON_OkGSlice(t *testing.T) {
	data, err := json.Marshal(Ok(SliceOf[Int](10, 20, 30)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"ok":[10,20,30]}` {
		t.Errorf("got %s, want {\"ok\":[10,20,30]}", data)
	}
}

func TestResultMarshalJSON_OkNilPointer(t *testing.T) {
	data, err := json.Marshal(Ok[*int](nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"ok":null}` {
		t.Errorf("got %s, want {\"ok\":null}", data)
	}
}

func TestResultMarshalJSON_OkUnsupportedValue(t *testing.T) {
	if _, err := json.Marshal(Ok(make(chan int))); err == nil {
		t.Fatal("expected error for unmarshalable Ok value")
	}
}

func TestResultMarshalJSON_Err(t *testing.T) {
	data, err := json.Marshal(Err[int](errors.New("boom")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"err":"boom"}` {
		t.Errorf("got %s, want {\"err\":\"boom\"}", data)
	}
}

func TestResultMarshalJSON_ErrWrapped(t *testing.T) {
	base := errors.New("io failure")
	wrapped := fmt.Errorf("read config: %w", base)

	data, err := json.Marshal(Err[string](wrapped))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"err":"read config: io failure"}` {
		t.Errorf("got %s, want {\"err\":\"read config: io failure\"}", data)
	}
}

func TestResultMarshalJSON_ErrEscaping(t *testing.T) {
	data, err := json.Marshal(Err[int](errors.New(`bad "input"` + "\n")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"err":"bad \"input\"\n"}` {
		t.Errorf("got %s, want {\"err\":\"bad \\\"input\\\"\\n\"}", data)
	}
}

func TestResultUnmarshalJSON_OkInt(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`{"ok":42}`), &res); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.IsErr() || res.Ok() != 42 {
		t.Errorf("got %v, want Ok(42)", res)
	}
}

func TestResultUnmarshalJSON_Err(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`{"err":"boom"}`), &res); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.IsOk() {
		t.Fatalf("expected Err, got Ok(%d)", res.Ok())
	}

	if res.Err().Error() != "boom" {
		t.Errorf("got %q, want %q", res.Err().Error(), "boom")
	}
}

func TestResultUnmarshalJSON_OkNullPointer(t *testing.T) {
	var res Result[*int]
	if err := json.Unmarshal([]byte(`{"ok":null}`), &res); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.IsErr() {
		t.Fatalf("expected Ok, got Err(%v)", res.Err())
	}

	if res.Ok() != nil {
		t.Errorf("got %v, want Ok(nil)", res.Ok())
	}
}

func TestResultUnmarshalJSON_OkNullInt(t *testing.T) {
	res := Ok(99)
	if err := json.Unmarshal([]byte(`{"ok":null}`), &res); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.IsErr() {
		t.Fatalf("expected Ok, got Err(%v)", res.Err())
	}

	// Per encoding/json rules, null into an int is a no-op on the target;
	// a fresh zero value is decoded into, so the result is Ok(0).
	if res.Ok() != 0 {
		t.Errorf("got Ok(%d), want Ok(0)", res.Ok())
	}
}

func TestResultUnmarshalJSON_TypeMismatch(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`{"ok":"not an int"}`), &res); err == nil {
		t.Fatal("expected error for type mismatch in ok value")
	}
}

func TestResultUnmarshalJSON_EmptyObject(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`{}`), &res); err == nil {
		t.Fatal("expected error for empty object")
	}
}

func TestResultUnmarshalJSON_BothKeys(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`{"ok":1,"err":"boom"}`), &res); err == nil {
		t.Fatal("expected error for object with both keys")
	}
}

func TestResultUnmarshalJSON_DuplicateKeys(t *testing.T) {
	// BREAKING vs the old encoding/json implementation: duplicate keys were
	// previously accepted with last-wins semantics; the v2-backed decoder
	// rejects them.
	var res Result[int]
	if err := json.Unmarshal([]byte(`{"ok":1,"ok":2}`), &res); err == nil {
		t.Fatal("expected error for duplicate ok keys")
	}

	if err := json.Unmarshal([]byte(`{"err":"a","err":"b"}`), &res); err == nil {
		t.Fatal("expected error for duplicate err keys")
	}
}

func TestResultUnmarshalJSON_UnknownKey(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`{"value":1}`), &res); err == nil {
		t.Fatal("expected error for object with unknown key")
	}
}

func TestResultUnmarshalJSON_ExtraKey(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`{"ok":1,"extra":2}`), &res); err == nil {
		t.Fatal("expected error for object with extra key")
	}
}

func TestResultUnmarshalJSON_Array(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`[1,2,3]`), &res); err == nil {
		t.Fatal("expected error for JSON array")
	}
}

func TestResultUnmarshalJSON_Null(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`null`), &res); err == nil {
		t.Fatal("expected error for JSON null")
	}
}

func TestResultUnmarshalJSON_Garbage(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`not_json`), &res); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestResultUnmarshalJSON_ErrNotString(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`{"err":42}`), &res); err == nil {
		t.Fatal("expected error for non-string err value")
	}
}

func TestResultUnmarshalJSON_ErrNull(t *testing.T) {
	var res Result[int]
	if err := json.Unmarshal([]byte(`{"err":null}`), &res); err == nil {
		t.Fatal("expected error for null err value")
	}
}

func TestResultJSON_RoundTripOk(t *testing.T) {
	type Point struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	original := Ok(Point{X: 3, Y: 7})
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded Result[Point]
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.IsErr() || decoded.Ok() != original.Ok() {
		t.Errorf("round trip failed: got %v, want %v", decoded, original)
	}
}

func TestResultJSON_RoundTripErr(t *testing.T) {
	base := errors.New("io failure")
	original := Err[int](fmt.Errorf("read config: %w", base))

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded Result[int]
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.IsOk() {
		t.Fatalf("expected Err, got Ok(%d)", decoded.Ok())
	}

	if decoded.Err().Error() != original.Err().Error() {
		t.Errorf("message: got %q, want %q", decoded.Err().Error(), original.Err().Error())
	}

	// Only the message survives the round trip: the decoded error is a plain
	// errors.New value, so the errors.Is chain against base is lost.
	if decoded.ErrIs(base) {
		t.Error("expected errors.Is chain to be lost after round trip")
	}
}

func TestResultJSON_RoundTripOkNilPointer(t *testing.T) {
	original := Ok[*string](nil)
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded Result[*string]
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.IsErr() || decoded.Ok() != nil {
		t.Errorf("round trip failed: got %v, want Ok(nil)", decoded)
	}
}

func TestResultJSON_StructField(t *testing.T) {
	type Response struct {
		Name   string      `json:"name"`
		Lookup Result[int] `json:"lookup"`
	}

	data, err := json.Marshal(Response{Name: "x", Lookup: Ok(5)})
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"name":"x","lookup":{"ok":5}}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}

	data, err = json.Marshal(Response{Name: "x", Lookup: Err[int](errors.New("missing"))})
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected = `{"name":"x","lookup":{"err":"missing"}}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}

	var decoded Response
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.Lookup.IsOk() || decoded.Lookup.Err().Error() != "missing" {
		t.Errorf("got %v, want Err(missing)", decoded.Lookup)
	}
}

func TestResultJSON_ResultInsideOption(t *testing.T) {
	data, err := json.Marshal(Some(Ok(7)))
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if string(data) != `{"ok":7}` {
		t.Errorf("got %s, want {\"ok\":7}", data)
	}

	data, err = json.Marshal(None[Result[int]]())
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if string(data) != "null" {
		t.Errorf("got %s, want null", data)
	}

	var decoded Option[Result[int]]
	if err := json.Unmarshal([]byte(`{"err":"boom"}`), &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.IsNone() || decoded.Some().IsOk() || decoded.Some().Err().Error() != "boom" {
		t.Errorf("got %v, want Some(Err(boom))", decoded)
	}
}

func TestResultJSON_OptionInsideResult(t *testing.T) {
	data, err := json.Marshal(Ok(Some(5)))
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if string(data) != `{"ok":5}` {
		t.Errorf("got %s, want {\"ok\":5}", data)
	}

	data, err = json.Marshal(Ok(None[int]()))
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if string(data) != `{"ok":null}` {
		t.Errorf("got %s, want {\"ok\":null}", data)
	}

	var decoded Result[Option[int]]
	if err := json.Unmarshal([]byte(`{"ok":null}`), &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.IsErr() || decoded.Ok().IsSome() {
		t.Errorf("got %v, want Ok(None)", decoded)
	}

	if err := json.Unmarshal([]byte(`{"ok":9}`), &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.IsErr() || decoded.Ok().IsNone() || decoded.Ok().Some() != 9 {
		t.Errorf("got %v, want Ok(Some(9))", decoded)
	}
}

func TestResultJSON_NestedResult(t *testing.T) {
	original := Ok(Ok(3))
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if string(data) != `{"ok":{"ok":3}}` {
		t.Errorf("got %s, want {\"ok\":{\"ok\":3}}", data)
	}

	var decoded Result[Result[int]]
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.IsErr() || decoded.Ok().IsErr() || decoded.Ok().Ok() != 3 {
		t.Errorf("round trip failed: got %v, want Ok(Ok(3))", decoded)
	}
}

func TestResultJSONV2_MarshalOk(t *testing.T) {
	data, err := jsonv2.Marshal(Ok(42))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"ok":42}` {
		t.Errorf("got %s, want {\"ok\":42}", data)
	}
}

func TestResultJSONV2_MarshalErr(t *testing.T) {
	data, err := jsonv2.Marshal(Err[int](errors.New("boom")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"err":"boom"}` {
		t.Errorf("got %s, want {\"err\":\"boom\"}", data)
	}
}

func TestResultJSONV2_MarshalOkNilSlice(t *testing.T) {
	// v2 semantics: a nil slice marshals as [], not null.
	data, err := jsonv2.Marshal(Ok[[]int](nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"ok":[]}` {
		t.Errorf("got %s, want {\"ok\":[]}", data)
	}
}

func TestResultJSONV2_MarshalOkInvalidUTF8(t *testing.T) {
	// v2 strictness: strings with invalid UTF-8 are a marshal error instead of
	// being silently replaced with U+FFFD.
	if _, err := jsonv2.Marshal(Ok("\xff\xfe")); err == nil {
		t.Fatal("expected error for invalid UTF-8 in Ok string")
	}
}

func TestResultJSONV2_MarshalErrInvalidUTF8(t *testing.T) {
	if _, err := jsonv2.Marshal(Err[int](errors.New("bad \xff message"))); err == nil {
		t.Fatal("expected error for invalid UTF-8 in error message")
	}
}

func TestResultJSONV2_UnmarshalOk(t *testing.T) {
	var res Result[int]
	if err := jsonv2.Unmarshal([]byte(`{"ok":42}`), &res); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.IsErr() || res.Ok() != 42 {
		t.Errorf("got %v, want Ok(42)", res)
	}
}

func TestResultJSONV2_UnmarshalErr(t *testing.T) {
	var res Result[int]
	if err := jsonv2.Unmarshal([]byte(`{"err":"boom"}`), &res); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.IsOk() || res.Err().Error() != "boom" {
		t.Errorf("got %v, want Err(boom)", res)
	}
}

func TestResultJSONV2_UnmarshalDuplicateKeys(t *testing.T) {
	var res Result[int]
	if err := jsonv2.Unmarshal([]byte(`{"ok":1,"ok":2}`), &res); err == nil {
		t.Fatal("expected error for duplicate ok keys")
	}

	if err := jsonv2.Unmarshal([]byte(`{"err":"a","err":"b"}`), &res); err == nil {
		t.Fatal("expected error for duplicate err keys")
	}
}

func TestResultJSONV2_UnmarshalShapeErrors(t *testing.T) {
	cases := []string{
		`{}`,
		`{"ok":1,"err":"boom"}`,
		`{"ok":1,"extra":2}`,
		`{"value":1}`,
		`{"err":42}`,
		`{"err":null}`,
		`[1,2,3]`,
		`null`,
		`"ok"`,
		`not_json`,
	}

	for _, input := range cases {
		var res Result[int]
		if err := jsonv2.Unmarshal([]byte(input), &res); err == nil {
			t.Errorf("input %s: expected error, got %v", input, res)
		}
	}
}

func TestResultJSONV2_RoundTripOptionResult(t *testing.T) {
	original := Ok(Some(SliceOf[Int](1, 2)))
	data, err := jsonv2.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if string(data) != `{"ok":[1,2]}` {
		t.Errorf("got %s, want {\"ok\":[1,2]}", data)
	}

	var decoded Result[Option[Slice[Int]]]
	if err := jsonv2.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.IsErr() || decoded.Ok().IsNone() {
		t.Fatalf("got %v, want Ok(Some([1, 2]))", decoded)
	}

	s := decoded.Ok().Some()
	if s.Len() != 2 || s[0] != 1 || s[1] != 2 {
		t.Errorf("got %v, want Slice[1, 2]", s)
	}

	if err := jsonv2.Unmarshal([]byte(`{"ok":null}`), &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.IsErr() || decoded.Ok().IsSome() {
		t.Errorf("got %v, want Ok(None)", decoded)
	}
}

func TestResultJSON_UnmarshalStrictMessage(t *testing.T) {
	var res Result[int]
	err := json.Unmarshal([]byte(`{}`), &res)
	if err == nil {
		t.Fatal("expected error for empty object")
	}

	if !strings.Contains(err.Error(), "exactly one") {
		t.Errorf("got %q, want message mentioning exactly one key", err.Error())
	}
}
