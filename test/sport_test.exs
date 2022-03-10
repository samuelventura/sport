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

  test "size=0 test" do
    all =
      Enum.reduce(0..511, [], fn _i, list ->
        [<<"0">> | list]
      end)
      |> Enum.reverse()
      |> :erlang.iolist_to_binary()

    port0 = Sport.open("/tmp/tty.socat0", 9600, "8N1")
    port1 = Sport.open("/tmp/tty.socat1", 9600, "8N1")
    assert true == Sport.discard(port0)
    assert true == Sport.discard(port1)
    assert true == Sport.write(port0, all)
    assert true == Sport.drain(port0)
    # sporadic avail=0
    :timer.sleep(100)
    assert all == Sport.read(port1, 0, 0)
    assert true == Sport.write(port1, all)
    assert true == Sport.drain(port1)
    # sporadic avail=0
    :timer.sleep(100)
    assert all == Sport.read(port0, 0, 100)
    assert true == Sport.close(port0)
    assert true == Sport.close(port1)
  end

  test "read > 255 test" do
    all =
      Enum.reduce(0..511, [], fn _i, list ->
        [<<"0">> | list]
      end)
      |> Enum.reverse()
      |> :erlang.iolist_to_binary()

    port0 = Sport.open("/tmp/tty.socat0", 9600, "8N1")
    port1 = Sport.open("/tmp/tty.socat1", 9600, "8N1")
    assert true == Sport.discard(port0)
    assert true == Sport.discard(port1)
    assert true == Sport.write(port0, all)
    assert true == Sport.drain(port0)
    assert all == Sport.read(port1, byte_size(all))
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

  test "async packet n test" do
    port0 = Sport.open("/tmp/tty.socat0", 9600, "8N1")
    port1 = Sport.open("/tmp/tty.socat1", 9600, "8N1")
    assert true == Sport.discard(port0)
    assert true == Sport.discard(port1)
    assert true == Sport.write(port0, "hellohellohello")
    assert true == Sport.drain(port0)
    assert true == Sport.packetn(port1, 5)
    assert "hello" == Sport.receive(port1)
    assert "hello" == Sport.receive(port1)
    assert "hello" == Sport.receive(port1)
    assert true == Sport.close(port0)
    assert true == Sport.close(port1)
  end

  test "async packet c test" do
    port0 = Sport.open("/tmp/tty.socat0", 9600, "8N1")
    port1 = Sport.open("/tmp/tty.socat1", 9600, "8N1")
    assert true == Sport.discard(port0)
    assert true == Sport.discard(port1)
    assert true == Sport.write(port0, "hellohellohello")
    assert true == Sport.drain(port0)
    assert true == Sport.packetc(port1, ?o)
    assert "hello" == Sport.receive(port1)
    assert "hello" == Sport.receive(port1)
    assert "hello" == Sport.receive(port1)
    assert true == Sport.close(port0)
    assert true == Sport.close(port1)
  end

  test "open config test" do
    for config <- ["8N1", "8E1", "8O1", "7N1", "7E1", "7O1"] do
      port0 = Sport.open("/tmp/tty.socat0", 9600, config)
      port1 = Sport.open("/tmp/tty.socat1", 9600, config)
      assert true == Sport.write(port0, "hello")
      assert "hello" == Sport.read(port1, 5)
      assert true == Sport.close(port0)
      assert true == Sport.close(port1)
    end

    for config <- ["8N2", "8E2", "8O2", "7N2", "7E2", "7O2"] do
      port0 = Sport.open("/tmp/tty.socat0", 9600, config)
      port1 = Sport.open("/tmp/tty.socat1", 9600, config)
      assert true == Sport.write(port0, "hello")
      assert "hello" == Sport.read(port1, 5)
      assert true == Sport.close(port0)
      assert true == Sport.close(port1)
    end
  end
end
