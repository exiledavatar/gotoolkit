package meta

import (
	"os"
	"reflect"
	"regexp"
	"strings"
)

// ExpandEnvVars substitutes environment variables of the form ${ENV_VAR_NAME}
// if you have characters that need to be escaped, they should be surrounded in
// quotes in the source string.
func ExpandEnvVars(s string) string {
	re := regexp.MustCompile(`\$\{[^}]+\}`)

	envvars := map[string]string{}
	for _, m := range re.FindAllString(s, -1) {
		mre := regexp.MustCompile(`[${}]`)
		mtrimmed := mre.ReplaceAllString(m, "")
		// fmt.Printf("%s:\t%s\n", mtrimmed, os.Getenv(mtrimmed))
		envvars[m] = os.Getenv(mtrimmed)
	}

	for k, v := range envvars {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

// ExpandFileContents substitutes the placeholder with the contents on the first
// line of a file. It only accepts the pattern {file:/path/to/file}
// if you have characters that need to be escaped, they should be surrounded in
// quotes in the source string.
func ExpandFileContents(s string) string {
	re := regexp.MustCompile(`\{file:[^}]+\}`)

	files := map[string]string{}
	for _, filename := range re.FindAllString(s, -1) {
		idpattern := regexp.MustCompile(`(^\{file:)|(\}$)`)
		fn := idpattern.ReplaceAllString(filename, "")
		// fmt.Printf("%s:\t%s\n", mtrimmed, os.Getenv(mtrimmed))
		b, err := os.ReadFile(fn)
		if err != nil {
			panic(err)
		}
		fc := string(b)
		fca := strings.Split(fc, "\n")
		files[filename] = fca[0]
	}

	for k, v := range files {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

// ImplementsInterface is a simple wrapper for checking if a value implements and interface
// it is primarily useful in switch statements
// note: Reference must be passed as a type parameter
func ImplementsInterface[Reference any](value any) bool {
	_, ok := value.(Reference)
	return ok
}

func CanSet(a any) bool {
	switch v, err := ToValue(a); {
	case err != nil:
		panic(err)
	case v.Kind() == reflect.Invalid:
		return false
	default:
		return v.CanSet()
	}
}
