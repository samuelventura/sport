defmodule Sport do
  @moduledoc false
  defguard is_byte(n) when is_integer(n) and n >= 0 and n <= 0xFF

  def open(device, speed, config) do
    exec =
      case :os.type() do
        {:unix, :darwin} -> :code.priv_dir(:sport) ++ '/sport_darwin'
        {:unix, :linux} -> :code.priv_dir(:sport) ++ '/sport_linux'
      end

    args = [device, to_string(speed), config]
    opts = [:binary, :exit_status, packet: 2, args: args]
    Port.open({:spawn_executable, exec}, opts)
  end

  def close(port) do
    Port.close(port)
  end

  def drain(port) do
    Port.command(port, ['d'])
  end

  def discard(port) do
    Port.command(port, ['D'])
  end

  def write(port, data) do
    Port.command(port, ['w', data])
  end

  # tenths of a second
  def read(port, vmin \\ 0, vtime \\ 0) when is_byte(vmin) and is_byte(vtime) do
    Port.command(port, ['r', vmin, vtime])

    receive do
      {^port, {:data, data}} -> data
    end
  end
end
