module DrawTcTurtle exposing (draw, display)

import Html exposing (Html, text)
import Svg exposing (..)
import Svg.Attributes exposing (..)


draw : List (Float, Float) -> String
draw points =
    List.map (\(x, y) -> "(" ++ String.fromFloat x ++ ", " ++ String.fromFloat y ++ ")") points
        |> String.join ", "



display : List (Float, Float) -> Html msg
display points =
    let
        toLine (p1, p2) =
            line
                [ x1 (String.fromFloat (Tuple.first p1))
                , y1 (String.fromFloat (Tuple.second p1))
                , x2 (String.fromFloat (Tuple.first p2))
                , y2 (String.fromFloat (Tuple.second p2))
                , stroke "black"
                , strokeWidth "2"
                ]
                []

        pairs =
            case points of
                [] -> []
                _ -> List.map2 Tuple.pair points (List.drop 1 points)

    in
    svg [ viewBox "-250 -250 500 500", width "500", height "500" ]
        (List.map toLine pairs)

