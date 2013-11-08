package suexec

import (
	"bufio"
	"os"
	"path"
	"strings"
)

func ReadRubyVersion(dir string) (string, error) {
	ruby_version := pathToRubyVersion(dir)
	file, err := os.Open(ruby_version)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	/* need 1st line only */
	scanner.Scan()
	line := scanner.Text()
	return strings.TrimSpace(line), nil
}

func pathToRubyVersion(dir string) string {
	return path.Join(dir, ".ruby-version")
}

func InjectRubyPathEnv(path string, environ []string) []string {
	for i, env := range environ {
		if strings.HasPrefix(env, "PATH=") {
			tokens := strings.Split(env, "=")
			environ[i] = tokens[0] + "=" + path + ":" + tokens[1]
			break
		}
	}
	return environ
}
