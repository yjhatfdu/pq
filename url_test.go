package pq

import (
	"database/sql"
	"testing"
)

func TestSimpleParseURL(t *testing.T) {
	expected := "host=hostname.remote"
	str, err := ParseURL("postgres://hostname.remote")
	if err != nil {
		t.Fatal(err)
	}

	if str != expected {
		t.Fatalf("unexpected result from ParseURL:\n+ %v\n- %v", str, expected)
	}
}

func TestIPv6LoopbackParseURL(t *testing.T) {
	expected := "host=::1 port=1234"
	str, err := ParseURL("postgres://[::1]:1234")
	if err != nil {
		t.Fatal(err)
	}

	if str != expected {
		t.Fatalf("unexpected result from ParseURL:\n+ %v\n- %v", str, expected)
	}
}

func TestFullParseURL(t *testing.T) {
	expected := `dbname=database host=hostname.remote password=top\ secret port=1234 user=username`
	str, err := ParseURL("postgres://username:top%20secret@hostname.remote:1234/database")
	if err != nil {
		t.Fatal(err)
	}

	if str != expected {
		t.Fatalf("unexpected result from ParseURL:\n+ %s\n- %s", str, expected)
	}
}

func TestInvalidProtocolParseURL(t *testing.T) {
	_, err := ParseURL("http://hostname.remote")
	switch err {
	case nil:
		t.Fatal("Expected an error from parsing invalid protocol")
	default:
		msg := "invalid connection protocol: http"
		if err.Error() != msg {
			t.Fatalf("Unexpected error message:\n+ %s\n- %s",
				err.Error(), msg)
		}
	}
}

func TestMinimalURL(t *testing.T) {
	cs, err := ParseURL("postgres://")
	if err != nil {
		t.Fatal(err)
	}

	if cs != "" {
		t.Fatalf("expected blank connection string, got: %q", cs)
	}
}

func TestMultiHostURL(t *testing.T) {
	expected := `attr=true dbname=database host=host1,host2 password=top\ secret user=username`
	str, err := ParseURL("postgres://username:top%20secret@host1,host2/database?attr=true")
	if err != nil {
		t.Fatal(err)
	}
	if str != expected {
		t.Fatalf("unexpected result from ParseURL:\n+ %s\n- %s", str, expected)
	}
}

func TestSplitMultiHostURL(t *testing.T) {
	urls, err := SplitMultiHostUrl("postgres://username:top%20secret@host1,host2/database?attr=true")
	if err != nil {
		t.Fatal(err)
	}
	if urls[0] != "postgres://username:top%20secret@host1/database?attr=true" {
		t.Fatalf("unexpected result from ParseURL:\n+ %s\n-", urls)
	}
	if urls[1] != "postgres://username:top%20secret@host2/database?attr=true" {
		t.Fatalf("unexpected result from ParseURL:\n+ %s\n-", urls)
	}
}
func TestSplitMultiHostURL2(t *testing.T) {
	urls, err := SplitMultiHostUrl("postgres://username:top%20secret@host1:5432,host2:5433/database?attr=true")
	t.Log(urls)
	if err != nil {
		t.Fatal(err)
	}
	if urls[0] != "postgres://username:top%20secret@host1:5432/database?attr=true" {
		t.Fatalf("unexpected result from ParseURL:\n+ %s\n-", urls)
	}
	if urls[1] != "postgres://username:top%20secret@host2:5433/database?attr=true" {
		t.Fatalf("unexpected result from ParseURL:\n+ %s\n-", urls)
	}
}

func TestConnMultiHost(t *testing.T) {
	url := "postgres://localhost,localhost/postgres?target_session_attrs=read_write&sslmode=disable"
	conn, err := sql.Open("postgres", url)
	t.Log(conn)
	t.Log(err)
	_, err = conn.Query("select 1")
	t.Log(err)
}
