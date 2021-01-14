package store

import (
	"bytes"
	"crypto/rand"
	"os"
	"testing"

	"github.com/yexelm/shorty/config"
)

var db *Storage

func TestMain(m *testing.M) {
	cfg := config.New()

	db, _ = New(cfg.RedisURL, cfg.TestStorageDb)
	code := m.Run()
	conn := db.Pool.Get()
	_, _ = conn.Do("FLUSHDB")
	db.Close()
	os.Exit(code)
}

//TestShortByLong checks if each new URL gets a new short alias and each URL already contained in Redis get the
//same short alias.
func TestShortByLong(t *testing.T) {
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
			got, gotErr := db.ShortByLong(tc.longURL)
			if !bytes.Equal(got, tc.want) {
				t.Errorf("\ngot:  %q\nwant: %q\n", got, tc.want)
			}
			if gotErr != nil && !tc.wantErr {
				t.Errorf("unexpected error: %v\n", gotErr)
			}
		})
	}
}

// TestSaveFullLongByShort tries to insert <urlsNum> URLs into Redis and checks if each of them can be retrieved by its
// short alias correctly afterwards.
func TestSaveFullLongByShort(t *testing.T) {
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
		long, err := db.LongByShort(s)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(long, []byte(l)) {
			t.Fatalf("long in map: %q, long in Redis: %q", long, l)
		}
	}
}
