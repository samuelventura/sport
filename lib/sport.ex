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

  def drain(port, sync \\ true) do
    case sync do
      false ->
        Port.command(port, ["d\x00"])

      true ->
        true = Port.command(port, ["d\x01"])

        receive do
          {^port, {:data, "d\x01"}} -> true
        end
    end
  end

  def discard(port, sync \\ true) do
    case sync do
      false ->
        Port.command(port, ["D\x00"])

      true ->
        true = Port.command(port, ["D\x01"])

        receive do
          {^port, {:data, "D\x01"}} -> true
        end
    end
  end

  def write(port, data) do
    Port.command(port, ['w', data])
  end

  # tenths of a second
  def read(port, vmin \\ 0, vtime \\ 0) when is_byte(vmin) and is_byte(vtime) do
    true = Port.command(port, ['r', vmin, vtime])

    receive do
      {^port, {:data, <<"r", data::binary>>}} -> data
    end
  end
end
