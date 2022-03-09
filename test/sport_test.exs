defmodule SportTest do
  use ExUnit.Case, async: false
  doctest Sport

  test "basic test" do
    port0 = Sport.open("/tmp/tty.socat0", 9600, "8N1")
    port1 = Sport.open("/tmp/tty.socat1", 9600, "8N1")
    assert true == Sport.discard(port0, false)
    assert true == Sport.discard(port1, true)
    assert true == Sport.write(port0, "hello", true)
    assert true == Sport.drain(port0, true)
    assert "hello" == Sport.read(port1, 5)
    assert true == Sport.close(port0)
    assert true == Sport.close(port1)

    port0 = Sport.open("/tmp/tty.socat0", 9600, "8N1")
    port1 = Sport.open("/tmp/tty.socat1", 9600, "8N1")
    assert true == Sport.write(port1, "hello", false)
    assert true == Sport.drain(port1, false)
    assert "he" == Sport.read(port0, 2)
    assert "llo" == Sport.read(port0, 3)
    assert "" == Sport.read(port0)
    assert true == Sport.close(port0)
    assert true == Sport.close(port1)
  end

  test "raw test" do
    lo =
      Enum.reduce(0..127, [], fn i, list ->
        [<<i>> | list]
      end)
      |> Enum.reverse()
      |> :erlang.iolist_to_binary()

    hi =
      Enum.reduce(128..255, [], fn i, list ->
        [<<i>> | list]
      end)
      |> Enum.reverse()
      |> :erlang.iolist_to_binary()

    port0 = Sport.open("/tmp/tty.socat0", 9600, "8N1")
    port1 = Sport.open("/tmp/tty.socat1", 9600, "8N1")
    assert true == Sport.discard(port0)
    assert true == Sport.discard(port1)
    assert true == Sport.write(port0, lo)
    assert true == Sport.drain(port0)
    assert lo == Sport.read(port1, byte_size(lo))
    assert true == Sport.write(port0, hi)
    assert true == Sport.drain(port0)
    assert hi == Sport.read(port1, byte_size(hi))
    assert true == Sport.close(port0)
    assert true == Sport.close(port1)
  end
end
