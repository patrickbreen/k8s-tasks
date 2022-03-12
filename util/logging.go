package util

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var Log zerolog.Logger

func InitLogger(r *gin.Engine) {
	// make my own logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}

	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("LEET APP | msg:\"%s\"", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	Log = zerolog.New(output).With().Timestamp().Caller().Logger()

	// request ID middleware
	r.Use(
		requestid.New(
			requestid.WithGenerator(func() string {
				return uuid.NewString()
			}),
		),
	)
	// set how gin logs request
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - %s [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.Request.Header.Get("X-Request-ID"),
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

}
