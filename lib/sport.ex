defmodule Sport do
  @moduledoc false
  defguard is_byte(n) when is_integer(n) and n >= 0 and n <= 0xFF
  defguard is_word(n) when is_integer(n) and n >= 0 and n <= 0xFFFF

  defguard is_hundred(n)
           when is_integer(n) and n >= 0 and rem(n, 100) == 0 and div(n, 100) <= 0xFF

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
          {^port, {:exit_status, status}} -> {:exit, status}
          {^port, {:data, "d"}} -> true
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
          {^port, {:exit_status, status}} -> {:exit, status}
          {^port, {:data, "D"}} -> true
        end
    end
  end

  def write(port, data, sync \\ false) do
    case sync do
      false ->
        Port.command(port, ['w', 0, data])

      true ->
        true = Port.command(port, ['w', 1, data])

        receive do
          {^port, {:exit_status, status}} -> {:exit, status}
          {^port, {:data, "w"}} -> true
        end
    end
  end

  def read(port, size \\ 0, timeout \\ 0) when is_word(size) and is_hundred(timeout) do
    true = Port.command(port, ['r', div(size, 256), rem(size, 256), div(timeout, 100)])

    receive do
      {^port, {:exit_status, status}} -> {:exit, status}
      {^port, {:data, <<"r", data::binary>>}} -> data
    end
  end

  def receive(port, timeout \\ -1) do
    cond do
      timeout < 0 ->
        receive do
          {^port, {:exit_status, status}} -> {:exit, status}
          {^port, {:data, <<"r", data::binary>>}} -> data
        end

      true ->
        receive do
          {^port, {:exit_status, status}} -> {:exit, status}
          {^port, {:data, <<"r", data::binary>>}} -> data
        after
          timeout ->
            :timeout
        end
    end
  end

  def packetn(port, size) when is_word(size) do
    true = Port.command(port, ['p', 'n', div(size, 256), rem(size, 256)])
  end

  def packetc(port, delim) when is_byte(delim) do
    true = Port.command(port, ['p', 'c', delim])
  end
end
