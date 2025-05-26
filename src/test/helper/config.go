package testhelper

import (
	"os"
	"testing"

	"github.com/jackysum/go-template/src/utils/file"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

type config struct {
	DatabaseConnString string
}

func Config(t *testing.T) config {
	t.Helper()

	err := godotenv.Load(file.AbsolutePath(".env.test"))
	require.NoError(t, err)

	return config{
		DatabaseConnString: os.Getenv("DATABASE_CONN_STRING"),
	}
}
