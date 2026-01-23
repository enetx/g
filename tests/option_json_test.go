package g_test

import (
	"encoding/json"
	"testing"

	. "github.com/enetx/g"
)

func TestOptionMarshalJSON_SomeInt(t *testing.T) {
	data, err := json.Marshal(Some(42))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "42" {
		t.Errorf("got %s, want 42", data)
	}
}

func TestOptionMarshalJSON_SomeString(t *testing.T) {
	data, err := json.Marshal(Some("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `"hello"` {
		t.Errorf("got %s, want %q", data, "hello")
	}
}

func TestOptionMarshalJSON_SomeBool(t *testing.T) {
	data, err := json.Marshal(Some(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "true" {
		t.Errorf("got %s, want true", data)
	}
}

func TestOptionMarshalJSON_SomeFloat(t *testing.T) {
	data, err := json.Marshal(Some(3.14))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "3.14" {
		t.Errorf("got %s, want 3.14", data)
	}
}

func TestOptionMarshalJSON_SomeSlice(t *testing.T) {
	data, err := json.Marshal(Some([]int{1, 2, 3}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "[1,2,3]" {
		t.Errorf("got %s, want [1,2,3]", data)
	}
}

func TestOptionMarshalJSON_SomeGSlice(t *testing.T) {
	data, err := json.Marshal(Some(SliceOf(10, 20, 30)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "[10,20,30]" {
		t.Errorf("got %s, want [10,20,30]", data)
	}
}

func TestOptionMarshalJSON_SomeMap(t *testing.T) {
	data, err := json.Marshal(Some(map[string]int{"a": 1}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `{"a":1}` {
		t.Errorf("got %s, want {\"a\":1}", data)
	}
}

func TestOptionMarshalJSON_None(t *testing.T) {
	data, err := json.Marshal(None[int]())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "null" {
		t.Errorf("got %s, want null", data)
	}
}

func TestOptionMarshalJSON_NoneSlice(t *testing.T) {
	data, err := json.Marshal(None[[]int]())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "null" {
		t.Errorf("got %s, want null", data)
	}
}

func TestOptionMarshalJSON_NestedOption(t *testing.T) {
	data, err := json.Marshal(Some(Some(5)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "5" {
		t.Errorf("got %s, want 5", data)
	}
}

func TestOptionUnmarshalJSON_Int(t *testing.T) {
	var opt Option[int]
	if err := json.Unmarshal([]byte("42"), &opt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opt.IsNone() || opt.Some() != 42 {
		t.Errorf("got %v, want Some(42)", opt)
	}
}

func TestOptionUnmarshalJSON_String(t *testing.T) {
	var opt Option[string]
	if err := json.Unmarshal([]byte(`"world"`), &opt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opt.IsNone() || opt.Some() != "world" {
		t.Errorf("got %v, want Some(world)", opt)
	}
}

func TestOptionUnmarshalJSON_Slice(t *testing.T) {
	var opt Option[[]int]
	if err := json.Unmarshal([]byte("[1,2,3]"), &opt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opt.IsNone() {
		t.Fatal("expected Some, got None")
	}

	s := opt.Some()
	if len(s) != 3 || s[0] != 1 || s[1] != 2 || s[2] != 3 {
		t.Errorf("got %v, want [1 2 3]", s)
	}
}

func TestOptionUnmarshalJSON_GSlice(t *testing.T) {
	var opt Option[Slice[Int]]
	if err := json.Unmarshal([]byte("[10,20,30]"), &opt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opt.IsNone() {
		t.Fatal("expected Some, got None")
	}

	s := opt.Some()
	if s.Len() != 3 || s[0] != 10 || s[1] != 20 || s[2] != 30 {
		t.Errorf("got %v, want Slice[10, 20, 30]", s)
	}
}

func TestOptionUnmarshalJSON_Null(t *testing.T) {
	var opt Option[int]
	if err := json.Unmarshal([]byte("null"), &opt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opt.IsSome() {
		t.Errorf("expected None, got Some(%d)", opt.Some())
	}
}

func TestOptionUnmarshalJSON_NullResetsValue(t *testing.T) {
	opt := Some(99)
	if err := json.Unmarshal([]byte("null"), &opt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opt.IsSome() {
		t.Errorf("expected None after null, got Some(%d)", opt.Some())
	}
}

func TestOptionUnmarshalJSON_InvalidJSON(t *testing.T) {
	var opt Option[int]
	if err := json.Unmarshal([]byte("not_json"), &opt); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestOptionJSON_StructFieldNoTag(t *testing.T) {
	type Item struct {
		Name  string      `json:"name"`
		Value Option[int] `json:"value"`
	}

	data, _ := json.Marshal(Item{Name: "x", Value: None[int]()})
	expected := `{"name":"x","value":null}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}

	data, _ = json.Marshal(Item{Name: "x", Value: Some(5)})
	expected = `{"name":"x","value":5}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}
}

func TestOptionJSON_StructFieldOmitzero(t *testing.T) {
	type Item struct {
		Name  string      `json:"name"`
		Value Option[int] `json:"value,omitzero"`
	}

	data, _ := json.Marshal(Item{Name: "x", Value: None[int]()})
	expected := `{"name":"x"}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}

	data, _ = json.Marshal(Item{Name: "x", Value: Some(5)})
	expected = `{"name":"x","value":5}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}
}

func TestOptionJSON_StructFieldOmitempty(t *testing.T) {
	type Item struct {
		Name  string      `json:"name"`
		Value Option[int] `json:"value,omitempty"`
	}

	data, _ := json.Marshal(Item{Name: "x", Value: None[int]()})
	expected := `{"name":"x","value":null}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}

	data, _ = json.Marshal(Item{Name: "x", Value: Some(5)})
	expected = `{"name":"x","value":5}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}
}

func TestOptionJSON_OmitzeroVsOmitempty(t *testing.T) {
	type Mixed struct {
		A Option[int]    `json:"a,omitzero"`
		B Option[string] `json:"b,omitempty"`
	}

	data, _ := json.Marshal(Mixed{A: None[int](), B: None[string]()})
	expected := `{"b":null}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}

	data, _ = json.Marshal(Mixed{A: Some(1), B: Some("hi")})
	expected = `{"a":1,"b":"hi"}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}
}

func TestOptionJSON_SliceFieldsOmitzeroVsOmitempty(t *testing.T) {
	type Data struct {
		Tags   Option[Slice[String]] `json:"tags,omitzero"`
		Scores Option[Slice[Int]]    `json:"scores,omitempty"`
	}

	data, _ := json.Marshal(Data{Tags: None[Slice[String]](), Scores: None[Slice[Int]]()})
	expected := `{"scores":null}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}

	data, _ = json.Marshal(Data{
		Tags:   Some(SliceOf[String]("a", "b")),
		Scores: Some(SliceOf[Int](1, 2)),
	})
	expected = `{"tags":["a","b"],"scores":[1,2]}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", data, expected)
	}
}

func TestOptionJSON_UnmarshalStructFull(t *testing.T) {
	type User struct {
		Name   string             `json:"name"`
		Age    Option[int]        `json:"age"`
		Email  Option[string]     `json:"email"`
		Scores Option[Slice[Int]] `json:"scores"`
	}

	input := `{"name":"Alice","age":30,"email":"a@b.c","scores":[95,87]}`
	var u User
	if err := json.Unmarshal([]byte(input), &u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Name != "Alice" {
		t.Errorf("Name: got %s, want Alice", u.Name)
	}

	if u.Age.IsNone() || u.Age.Some() != 30 {
		t.Errorf("Age: got %v, want Some(30)", u.Age)
	}

	if u.Email.IsNone() || u.Email.Some() != "a@b.c" {
		t.Errorf("Email: got %v, want Some(a@b.c)", u.Email)
	}

	if u.Scores.IsNone() {
		t.Fatal("Scores: expected Some, got None")
	}

	scores := u.Scores.Some()
	if scores.Len() != 2 || scores[0] != 95 || scores[1] != 87 {
		t.Errorf("Scores: got %v, want [95, 87]", scores)
	}
}

func TestOptionJSON_UnmarshalStructNulls(t *testing.T) {
	type User struct {
		Name  string         `json:"name"`
		Age   Option[int]    `json:"age"`
		Email Option[string] `json:"email"`
	}

	input := `{"name":"Bob","age":null,"email":null}`
	var u User
	if err := json.Unmarshal([]byte(input), &u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Age.IsSome() {
		t.Errorf("Age: expected None, got Some(%d)", u.Age.Some())
	}

	if u.Email.IsSome() {
		t.Errorf("Email: expected None, got Some(%s)", u.Email.Some())
	}
}

func TestOptionJSON_UnmarshalStructMissingFields(t *testing.T) {
	type User struct {
		Name  string         `json:"name"`
		Age   Option[int]    `json:"age"`
		Email Option[string] `json:"email"`
	}

	input := `{"name":"Carol"}`
	var u User
	if err := json.Unmarshal([]byte(input), &u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Age.IsSome() {
		t.Errorf("Age: expected None for missing field, got Some(%d)", u.Age.Some())
	}

	if u.Email.IsSome() {
		t.Errorf("Email: expected None for missing field, got Some(%s)", u.Email.Some())
	}
}

func TestOptionJSON_RoundTripSome(t *testing.T) {
	original := Some(SliceOf(1, 2, 3))
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded Option[Slice[int]]
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.IsNone() {
		t.Fatal("expected Some, got None")
	}

	s := decoded.Some()
	if s.Len() != 3 || s[0] != 1 || s[1] != 2 || s[2] != 3 {
		t.Errorf("round trip failed: got %v", s)
	}
}

func TestOptionJSON_RoundTripNone(t *testing.T) {
	original := None[string]()
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded Option[string]
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.IsSome() {
		t.Errorf("expected None, got Some(%s)", decoded.Some())
	}
}

func TestOptionJSON_RoundTripStruct(t *testing.T) {
	type Config struct {
		Host  string                `json:"host"`
		Port  Option[int]           `json:"port,omitzero"`
		Tags  Option[Slice[String]] `json:"tags,omitzero"`
		Debug Option[bool]          `json:"debug,omitempty"`
	}

	original := Config{
		Host:  "localhost",
		Port:  Some(8080),
		Tags:  Some(SliceOf[String]("web", "api")),
		Debug: None[bool](),
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded Config
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.Host != "localhost" {
		t.Errorf("Host: got %s, want localhost", decoded.Host)
	}

	if decoded.Port.IsNone() || decoded.Port.Some() != 8080 {
		t.Errorf("Port: got %v, want Some(8080)", decoded.Port)
	}

	if decoded.Tags.IsNone() {
		t.Fatal("Tags: expected Some, got None")
	}

	tags := decoded.Tags.Some()
	if tags.Len() != 2 || tags[0] != "web" || tags[1] != "api" {
		t.Errorf("Tags: got %v, want [web, api]", tags)
	}
}
