package alert

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/mochammadshenna/arch-pba-template/internal/state"
	"github.com/mochammadshenna/arch-pba-template/internal/util/logger"
)

func Error(ctx context.Context, err error, webhookUrl, alertName, payload string, additionalData []byte) {
	// requestId := ctx.Value(state.HttpHeaders().RequestId)
	platformType := ctx.Value(state.HttpHeaders().PlatformType)
	version := ctx.Value(state.HttpHeaders().Version)

	if platformType == nil || platformType == "" {
		platformType = ""
	}

	if version == nil || version == "" {
		version = ""
	}

	// requestIdMessage := ""
	// if requestId != nil && requestId != "" {
	// 	requestIdMessage = fmt.Sprintf(requestIdTemplate, requestId)
	// }

	// additionalDataStr := strings.ReplaceAll(string(additionalData), `"`, `\"`)
	// payload := fmt.Sprintf(payload, time.Now().Unix(), requestIdMessage, platformType, version, err.Error(), additionalDataStr)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		defer func() {
			if r := recover(); r != nil {
				// logger.Error(ctx, r)
				return
			}
		}()
		_, err := http.DefaultClient.Post(webhookUrl, "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			logger.Errorf(ctx, "got an error while sending google chat alerting on alert.Error(); err=%+v", err)
		}
	}()
}

// const requestIdTemplate = `<https:https://ap-southeast-1.console.aws.amazon.com/cloudwatch/home?region=ap-southeast-1#logsV2:log-groups/log-group/$252Faws$252Flambda$252FPBA-API-Template/log-events/2024$252F04$252F15$252F$255B$2524LATEST$%s`
