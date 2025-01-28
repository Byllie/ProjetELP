module DrawTcTurtle exposing (draw)

import Html exposing (Html, text)

draw : List (Float, Float) -> String
draw points =
    List.map (\(x, y) -> "(" ++ String.fromFloat x ++ ", " ++ String.fromFloat y ++ ")") points
        |> String.join ", "
