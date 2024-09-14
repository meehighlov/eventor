package common

import (
	"log/slog"
)

func FSM(logger *slog.Logger, chatCahe *ChatCache, handlers map[string]CommandStepHandler) HandlerType {
	return func(event Event) error {
		ctx := event.GetContext()
		stepTODO := ctx.stepTODO

		logger.Debug("FSM", "setting command in progress:", event.GetCommand())
		ctx.SetCommandInProgress(event.GetCommand())

		nextStep := STEPS_DONE

		stepHandler, found := handlers[stepTODO]

		if !found {
			logger.Error("FSM: handler not found, resetting context", "step", stepTODO, "command", event.GetCommand())
			ctx.Reset()
			return nil
		}

		logger.Debug("FSM called", "command", event.GetCommand(), "handling step", stepTODO)

		nextStep, _ = stepHandler(event)

		if nextStep == STEPS_DONE {
			logger.Debug("FSM resetting context - termination step reached", "command", event.GetCommand(), "handling step", stepTODO)
			ctx.Reset()
			return nil
		}

		ctx.SetStepTODO(nextStep)

		logger.Debug("FSM", "command in progress after processing", ctx.GetCommandInProgress())

		return nil
	}
}
