package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				brokenPipe := isBrokenPipeError(err)

				if brokenPipe {
					renderBrokenPipe(c, err)

					return
				}

				renderRequestError(c, err)
			}
		}()
		c.Next()
	}
}

func renderRequestError(c *gin.Context, err interface{}) {
	httpRequest := extractRequest(c)

	if gin.IsDebugging() {
		log.Error().
			Str("error", fmt.Sprintf("%s", err)).
			Str("HTTPRequest", httpRequest).
			Str("Headers", extractHeaders(httpRequest)).
			Str("Stack", stack(4)).
			Msg("Panic recovered")
	} else {
		log.Error().
			Str("error", fmt.Sprintf("%s", err)).
			Str("HTTPRequest", httpRequest).
			Str("Stack", stack(4)).
			Msg("Panic recovered")
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}

func renderBrokenPipe(c *gin.Context, err interface{}) {
	httpRequest := extractRequest(c)

	log.Error().
		Str("error", fmt.Sprintf("%s", err)).
		Str("HTTPRequest", httpRequest).
		Msg("Broken pipe")

	c.Error(err.(error)) // nolint: errcheck
	c.Abort()
}

func extractRequest(c *gin.Context) string {
	httpRequest, _ := httputil.DumpRequest(c.Request, false)

	return string(httpRequest)
}

func extractHeaders(httpRequest string) string {
	headers := strings.Split(httpRequest, "\r\n")
	for idx, header := range headers {
		current := strings.Split(header, ":")
		if current[0] == "Authorization" {
			headers[idx] = current[0] + ": *"
		}
	}
	return strings.Join(headers, "\r\n")
}

func isBrokenPipeError(err interface{}) bool {
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") {
				return true
			}

			if strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				return true
			}
		}
	}

	return false
}

func stack(skip int) string {
	buf := new(bytes.Buffer)
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.String()
}

func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())

	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
