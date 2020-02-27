package command


func parseAndValidatePropagateEnv(propagateEnv interface{}, newCmd *Command) error {
	if propagate, ok := propagateEnv.(bool); ok {
		newCmd.PropagateEnv = propagate
	} else {
		return newCommandError(
			"must be a boolean",
			newCmd.Name,
			PROPAGATE_ENV,
			"",
		)
	}
	return nil
}
