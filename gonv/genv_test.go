package gonv

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func envFactory(t *testing.T, values map[string]string) string {
	envFile := filepath.Join(t.TempDir(), randomString(4)+".env")

	str := strings.Builder{}

	for k, v := range values {
		str.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	err := os.WriteFile(envFile, []byte(str.String()), 0644)
	if err != nil {
		t.Fatal(err)
	}

	return envFile
}

func Test_New(t *testing.T) {
	envFile := envFactory(t, map[string]string{
		"FOO": "bar",
	})

	g, err := New(envFile, "test.json")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, g.EnvFile, envFile)
	assert.Equal(t, g.ConfigFile, "test.json")
	assert.Equal(t, len(g.Vars), 1)
	assert.Equal(t, g.Vars["FOO"], "bar")
}

func TestGenv_Set_Clobber(t *testing.T) {
	type fields struct {
		EnvFile    string
		ConfigFile string
		Vars       map[string]string
	}
	type args struct {
		key     string
		value   string
		clobber bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{

			name: "success",
			fields: fields{
				EnvFile:    "test.env",
				ConfigFile: "test.json",
				Vars: map[string]string{
					"FOO": "bar",
				},
			},
			args: args{
				key:     "FOO",
				value:   "baz",
				clobber: false,
			},
			wantErr: true,
		},
		{
			name: "success clobber",
			fields: fields{
				EnvFile:    "test.env",
				ConfigFile: "test.json",
				Vars: map[string]string{
					"FOO": "bar",
				},
			},
			args: args{
				key:     "FOO",
				value:   "baz",
				clobber: true,
			},
			wantErr: false,
		},
		{
			name: "no conflicting keys",
			fields: fields{
				EnvFile:    "test.env",
				ConfigFile: "test.json",
				Vars: map[string]string{
					"FOO": "bar",
				},
			},
			args: args{
				key:     "BAR",
				value:   "baz",
				clobber: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Genv{
				EnvFile:    tt.fields.EnvFile,
				ConfigFile: tt.fields.ConfigFile,
				Vars:       tt.fields.Vars,
			}
			if err := g.Set(tt.args.key, tt.args.value, tt.args.clobber); (err != nil) != tt.wantErr {
				t.Errorf("Genv.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				assert.Equal(t, g.Vars[tt.args.key], tt.args.value)
			}
		})
	}
}

func TestGenv_Remove(t *testing.T) {
	type fields struct {
		EnvFile    string
		ConfigFile string
		Vars       map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				EnvFile:    "test.env",
				ConfigFile: "test.json",
				Vars: map[string]string{
					"FOO": "bar",
				},
			},
			args: args{
				key: "FOO",
			},
		},
		{
			name: "success",
			fields: fields{
				EnvFile:    "test.env",
				ConfigFile: "test.json",
				Vars: map[string]string{
					"FOO": "bar",
				},
			},
			args: args{
				key: "NOT_EXIST",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Genv{
				EnvFile:    tt.fields.EnvFile,
				ConfigFile: tt.fields.ConfigFile,
				Vars:       tt.fields.Vars,
			}

			g.Remove(tt.args.key)
			assert.NotContains(t, g.Vars, tt.args.key)
		})
	}
}
