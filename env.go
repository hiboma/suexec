package suexec

import (
	"os"
	"strings"
)

const AP_SAFE_PATH = "PATH=/usr/local/bin:/usr/bin:/bin"

var safe_env_lst = []string{
	/* variable name starts with */
	"HTTP_",
	"SSL_",

	/* variable name is */
	"AUTH_TYPE=",
	"CONTENT_LENGTH=",
	"CONTENT_TYPE=",
	"CONTEXT_DOCUMENT_ROOT=",
	"CONTEXT_PREFIX=",
	"DATE_GMT=",
	"DATE_LOCAL=",
	"DOCUMENT_NAME=",
	"DOCUMENT_PATH_INFO=",
	"DOCUMENT_ROOT=",
	"DOCUMENT_URI=",
	"GATEWAY_INTERFACE=",
	"HTTPS=",
	"LAST_MODIFIED=",
	"PATH_INFO=",
	"PATH_TRANSLATED=",
	"QUERY_STRING=",
	"QUERY_STRING_UNESCAPED=",
	"REMOTE_ADDR=",
	"REMOTE_HOST=",
	"REMOTE_IDENT=",
	"REMOTE_PORT=",
	"REMOTE_USER=",
	"REDIRECT_ERROR_NOTES=",
	"REDIRECT_HANDLER=",
	"REDIRECT_QUERY_STRING=",
	"REDIRECT_REMOTE_USER=",
	"REDIRECT_SCRIPT_FILENAME=",
	"REDIRECT_STATUS=",
	"REDIRECT_URL=",
	"REQUEST_METHOD=",
	"REQUEST_URI=",
	"REQUEST_SCHEME=",
	"SCRIPT_FILENAME=",
	"SCRIPT_NAME=",
	"SCRIPT_URI=",
	"SCRIPT_URL=",
	"SERVER_ADMIN=",
	"SERVER_NAME=",
	"SERVER_ADDR=",
	"SERVER_PORT=",
	"SERVER_PROTOCOL=",
	"SERVER_SIGNATURE=",
	"SERVER_SOFTWARE=",
	"UNIQUE_ID=",
	"USER_NAME=",
	"TZ=",
}

func CleanEnv() []string {
	return cleanupEnv(os.Environ(), safe_env_lst)
}

/* While cleaning the environment, the environment should be clean.
 * (e.g. malloc() may get the name of a file for writing debugging info.
 * Bad news if MALLOC_DEBUG_FILE is set to /etc/passwd.  Sprintf() may be
 * susceptible to bad locale settings....)
 * (from PR 2790)
 */
func cleanupEnv(environ []string, safe_env_lst []string) []string {
	cleanenv := []string{}
	for cidx := range environ {
		for idx := range safe_env_lst {
			if strings.HasPrefix(environ[cidx], safe_env_lst[idx]) {
				cleanenv = append(cleanenv, environ[cidx])
			}
		}
	}
	cleanenv = append(cleanenv, AP_SAFE_PATH)
	return cleanenv
}
