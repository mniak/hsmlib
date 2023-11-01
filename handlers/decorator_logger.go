package handlers

import (
	"fmt"

	"github.com/mniak/hsmlib"
)

func LoggerDecorator(logger hsmlib.Logger) hsmlib.CommandHandlerDecoratorFunc {
	return func(inner hsmlib.CommandHandler, cmd hsmlib.CommandWithHeader) (hsmlib.Response, error) {
		logger.Info("Command received",
			"command.header", fmt.Sprintf("[% 2X]", cmd.Header),
			"command.code", string(cmd.Code()),
			"command.data", cmd.Data,
		)
		resp, err := inner.Handle(cmd)
		logger.Info("Response sent",
			"response.error_code", string(resp.ErrorCode()),
			"response.data", resp.Data(),
		)
		return resp, err
	}
}
