module DrawTcTurtle exposing (display)

import Html exposing (Html, text)
import Svg exposing (..)
import Svg.Attributes exposing (..)
import Basics exposing (max)


{-| Représente les limites d'un ensemble de points dans un espace 2D.
- xmin : La coordonnée x minimale.
- xmax : La coordonnée x maximale.
- ymin : La coordonnée y minimale.
- ymax : La coordonnée y maximale.
-}
type alias Bounds =
    { xmin : Float
    , xmax : Float
    , ymin : Float
    , ymax : Float
    }


{-| Calcule les limites (bounds) d'une liste de points en 2D.
Cette fonction prend une liste de points (tuples de coordonnées x, y) et retourne
un enregistrement `Bounds` contenant les valeurs minimales et maximales pour x et y.

Exemple :
    findBounds [(0, 0), (10, 20), (-5, 5)] == { xmin = -5, xmax = 10, ymin = 0, ymax = 20 }
-}
findBounds : List (Float, Float) -> Bounds
findBounds points =
    let
        -- Extrait les coordonnées x et y des points
        xs = List.map Tuple.first points
        ys = List.map Tuple.second points

        -- Trouve les valeurs minimales et maximales pour x et y
        xmin = Maybe.withDefault 0 (List.minimum xs)
        xmax = Maybe.withDefault 0 (List.maximum xs)
        ymin = Maybe.withDefault 0 (List.minimum ys)
        ymax = Maybe.withDefault 0 (List.maximum ys)
    in
    { xmin = xmin, xmax = xmax, ymin = ymin, ymax = ymax }





{-| Affiche une liste de points sous forme de lignes SVG.
Cette fonction prend une liste de points (tuples de coordonnées x, y) et génère
un élément SVG contenant des lignes reliant les points successifs.

La taille du SVG est ajustée dynamiquement pour inclure tous les points avec une marge
de 20% autour des limites des points. Cela garantit que le dessin ne sort pas du cadre.

Exemple :
    display [(0, 0), (10, 20)] -- Génère un SVG avec une ligne de (0, 0) à (10, 20)
-}
display : List (Float, Float) -> Html msg
display points =
    let
        -- Calcule les limites des points
        bounds = findBounds points

        -- Calcule la largeur et la hauteur du dessin
        w = bounds.xmax - bounds.xmin
        h = bounds.ymax - bounds.ymin

        -- Ajoute une marge de 10% autour des limites
        padding = w * 0.1

        -- Définit la viewBox pour le SVG, en incluant la marge
        viewBoxValue =
            String.fromFloat (bounds.xmin - padding)
                ++ " "
                ++ String.fromFloat (bounds.ymin - padding)
                ++ " "
                ++ String.fromFloat (w + 2 * padding)
                ++ " "
                ++ String.fromFloat (h + 2 * padding)

        -- Convertit une paire de points en une ligne SVG
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

        -- Crée des paires de points successifs pour dessiner les lignes
        pairs =
            case points of
                [] -> []
                _ -> List.map2 Tuple.pair points (List.drop 1 points)
    in
    svg [ viewBox viewBoxValue, width "500", height "500" ] (List.map toLine pairs)
