package logger

import (
	"bytes"
	"fmt"
	"os"

	"cloud.google.com/go/compute/metadata"
	"github.com/blendle/zapdriver"
	"github.com/golang/protobuf/jsonpb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/runtime/protoiface"
)

func NewLogger() *zap.Logger {
	var config zap.Config

	if metadata.OnGCE() {
		config = zapdriver.NewProductionConfig()
	} else {
		config = zapdriver.NewDevelopmentConfig()
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.Level = getLevel()

	zapLogger, err := config.Build()
	if err != nil {
		return zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zap.NewAtomicLevel(),
		))
	}
	return zapLogger
}

func getLevel() zap.AtomicLevel {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	}
}

func PrototoJson(proto protoiface.MessageV1) []byte {
	m := jsonpb.Marshaler{}
	var buf bytes.Buffer
	if err := m.Marshal(&buf, proto); err != nil {
		return []byte(fmt.Sprintf("%v", err))
	}

	return buf.Bytes()
}

var Logger = NewLogger()
