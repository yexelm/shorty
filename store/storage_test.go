package store_test

import (
	"bytes"
	"crypto/rand"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"

	"github.com/yexelm/shorty/store"
)

var db *store.Storage

func TestMain(m *testing.M) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	db, _ = store.New(s.Addr(), 0)
	code := m.Run()
	conn := db.Pool.Get()
	_, _ = conn.Do("FLUSHDB")
	os.Exit(code)
}

// TestShorter checks if each new URL gets a new short alias and each URL already contained in Redis get the
// same short alias.
func Test_Shorter(t *testing.T) {
	tests := []struct {
		longURL []byte
		want    []byte
		wantErr bool
	}{
		{
			longURL: []byte("ya.ru"),
			want:    []byte("b"),
			wantErr: false,
		},
		{
			longURL: []byte("google.com"),
			want:    []byte("c"),
			wantErr: false,
		},
		{
			longURL: []byte("ya.ru"),
			want:    []byte("b"),
			wantErr: false,
		},
		{
			longURL: []byte("https://ya.ru"),
			want:    []byte("d"),
			wantErr: false,
		},
		{
			longURL: []byte("google.com"),
			want:    []byte("c"),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(string(tc.longURL), func(t *testing.T) {
			got, gotErr := db.Shorter(tc.longURL)
			if !bytes.Equal(got, tc.want) {
				t.Errorf("\ngot:  %q\nwant: %q\n", got, tc.want)
			}
			if gotErr != nil && !tc.wantErr {
				t.Errorf("unexpected error: %v\n", gotErr)
			}
		})
	}
}

// TestLonger tries to insert <urlsNum> URLs into Redis and checks if each of them can be retrieved by its
// short alias correctly afterwards.
func Test_Longer(t *testing.T) {
	const (
		urlsNum = 1000
	)

	buf := make([]byte, 100)
	longToShortMap := make(map[string][]byte)

	for i := 0; i < urlsNum; i++ {
		rand.Read(buf)

		longURL := buf
		short, err := db.SaveFull(longURL)
		if err != nil {
			t.Fatal(err)
		}

		longToShortMap[string(buf)] = short
	}

	for l, s := range longToShortMap {
		long, err := db.Longer(s)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(long, []byte(l)) {
			t.Fatalf("long in map: %q, long in Redis: %q", long, l)
		}
	}
}
