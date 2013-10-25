package env

import (
	. "github.com/r7kamura/gospel"
	"testing"
)

func TestCleanEnv(t *testing.T) {
	Describe(t, "cleanupEnv", func() {

		It("cleanup unsafe env", func() {

			environ := []string{
				"UNSAFE_ENV=1",
				"SAFE_ENV=1",
			}

			safe_env_lst := []string{
				"SAFE_ENV=",
			}

			cleanenv := cleanupEnv(environ, safe_env_lst)
			Expect(cleanenv[0]).To(Equal, "SAFE_ENV=1")
			Expect(cleanenv[1]).To(Equal, AP_SAFE_PATH)
		})
	})
}
