defmodule SportTest do
  use ExUnit.Case
  doctest Sport

  test "greets the world" do
    assert Sport.hello() == :world
  end
end
