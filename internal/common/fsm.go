package common

import (
	"context"
	"log/slog"
)

func FSM(logger *slog.Logger, chatCahe *ChatCache, handlers map[string]CommandStepHandler) HandlerType {
	return func(ctx context.Context, event Event) error {
		chatCtx := event.GetContext()
		stepTODO := chatCtx.stepTODO

		logger.Debug("FSM", "setting command in progress:", event.GetCommand())
		chatCtx.SetCommandInProgress(event.GetCommand())

		nextStep := STEPS_DONE

		stepHandler, found := handlers[stepTODO]

		if !found {
			logger.Error("FSM: handler not found, resetting context", "step", stepTODO, "command", event.GetCommand())
			chatCtx.Reset()
			return nil
		}

		logger.Debug("FSM called", "command", event.GetCommand(), "handling step", stepTODO)

		nextStep, _ = stepHandler(ctx, event)

		if nextStep == STEPS_DONE {
			logger.Debug("FSM resetting context - termination step reached", "command", event.GetCommand(), "handling step", stepTODO)
			chatCtx.Reset()
			return nil
		}

		chatCtx.SetStepTODO(nextStep)

		logger.Debug("FSM", "command in progress after processing", chatCtx.GetCommandInProgress())

		return nil
	}
}
