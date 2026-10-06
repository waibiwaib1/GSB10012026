defmodule Commanded.Aggregates.ExecutionContext do
  @moduledoc """
  Defines the arguments used to execute a command for an aggregate.

  The available options are:

    - `causation_id` - an optional UUID used to identify the cause of the
      command being dispatched. It is copied to the `causation_id` of any
      events created by the command.

    - `correlation_id` - a UUID used to correlate related commands and events.
      It is copied to the `correlation_id` of any events created by the
      command.

    - `command` - the command to execute, typically a struct
      (e.g. `%OpenBankAccount{...}`).

    - `metadata` - a map of key/value pairs containing the metadata to be
      associated with all events created by the command.

    - `handler` - the module that handles the command. It may be either the
      aggregate module itself or a separate command handler module.

    - `function` - the name of function, as an atom, that handles the command.
      The default value is `:execute`, used to support command dispatch directly
      to the aggregate module. For command handlers the `:handle` function is
      used.

    - `lifespan` - a module implementing the `Commanded.Aggregates.AggregateLifespan`
      behaviour to control the aggregate instance process lifespan. The default
      value, `Commanded.Aggregates.DefaultLifespan`, keeps the process running
      indefinitely.

  """

  alias Commanded.Aggregates.DefaultLifespan

  defstruct [
    causation_id: nil,
    command: nil,
    correlation_id: nil,
    metadata: %{},
    handler: nil,
    function: nil,
    lifespan: DefaultLifespan,
  ]
end
