defmodule SportTest do
  use ExUnit.Case
  doctest Sport

  test "basic API test" do
    port0 = Sport.open("/tmp/tty.socat0", 9600, "8N1")
    port1 = Sport.open("/tmp/tty.socat1", 9600, "8N1")
    assert true == Sport.discard(port0)
    assert true == Sport.discard(port1)
    assert true == Sport.write(port0, "hello")
    assert true == Sport.drain(port0)
    assert "hello" == Sport.read(port1, 5)
    assert true == Sport.write(port1, "hello")
    assert true == Sport.drain(port1)
    assert "he" == Sport.read(port0, 2)
    assert "llo" == Sport.read(port0, 3)
    assert "" == Sport.read(port0)
    assert true == Sport.close(port0)
    assert true == Sport.close(port1)
  end
end
