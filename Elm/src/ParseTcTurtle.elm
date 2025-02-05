module ParseTcTurtle exposing (read)

import Parser exposing (..)
-- En plus du cours, pas mal d'inspiration viennent de ce github pour gérer les instructions qui bouclent https://github.com/jinjor/elm-xml-parser/blob/master/src/XmlParser.elm
type Instruction
    = Forward Int
    | Left Int
    | Right Int
    | Repeat Int (List Instruction)

type alias Point =
    { x : Float
    , y : Float
    }

toRadians : Float -> Float
toRadians degrees =
    degrees * pi / 180


{-| Applique une instruction à l'état courant de la tortue.
Retourne le nouvel état et les points générés.

Paramètres :
- instruction : Commande à exécuter
- tortue      : État courant {position x,y, angle}

Retourne :
- (Nouvel état, Liste de points dessinés)
-}
appliquerInstruction : Instruction -> { x : Float, y : Float, angle : Float } -> ( { x : Float, y : Float, angle : Float }, List Point )
appliquerInstruction instruction tortue =
    case instruction of
        Forward distance ->
            let
                newX = tortue.x + toFloat distance * cos (toRadians tortue.angle)
                newY = tortue.y + toFloat distance * sin (toRadians tortue.angle)
                newPoint = { x = newX, y = newY }
            in
            ( { tortue | x = newX, y = newY }, [newPoint] )

        Left degrees ->
            ( { tortue | angle = tortue.angle - toFloat degrees }, [] )

        Right degrees ->
            ( { tortue | angle = tortue.angle + toFloat degrees }, [] )

        Repeat n instructions ->
            let
                ( finalTortue, points ) =
                    List.foldl
                        (\_ ( currentTortue, accPoints ) ->
                            let
                                ( newTortue, newPoints ) =
                                    appliquerInstructions instructions currentTortue
                            in
                            ( newTortue, accPoints ++ newPoints )
                        )
                        ( tortue, [] )
                        (List.repeat n ())
            in
            ( finalTortue, points )


{-| Exécute une liste d'instructions séquentiellement
-}
appliquerInstructions : List Instruction -> { x : Float, y : Float, angle : Float } -> ( { x : Float, y : Float, angle : Float }, List Point )
appliquerInstructions instructions tortue =
    List.foldl
        (\instruction ( currentTortue, accPoints ) ->
            let
                ( newTortue, newPoints ) =
                    appliquerInstruction instruction currentTortue
            in
            ( newTortue, accPoints ++ newPoints )
        )
        ( tortue, [] )
        instructions


-- PARSERS INDIVIDUELS

{-| Parse une commande Forward suivie d'un entier
Ex: "Forward 10"
-}
parseForward : Parser Instruction
parseForward =
    succeed Forward
        |. keyword "Forward"
        |. spaces
        |= int


{-| Parse une commande Left suivie d'un angle
Ex: "Left 90"
-}
parseLeft : Parser Instruction
parseLeft =
    succeed Left
        |. keyword "Left"
        |. spaces
        |= int


{-| Parse une commande Right suivie d'un angle
Ex: "Right 45"
-}
parseRight : Parser Instruction
parseRight =
    succeed Right
        |. keyword "Right"
        |. spaces
        |= int


{-| Parse une boucle Repeat avec ses instructions
Ex: "Repeat 4 [Forward 50, Left 90]"
-}
parseRepeat : Parser Instruction
parseRepeat =
    succeed Repeat
        |. keyword "Repeat"
        |. spaces
        |= int
        |. spaces
        |. symbol "["
        |. spaces
        |= lazy (\_ -> parseInstructions)
        |. spaces
        |. symbol "]"


{-| Parse une liste d'instructions séparées par des virgules
-}
parseInstructions : Parser (List Instruction)
parseInstructions =
    sequence
        { start = ""
        , separator = ","
        , end = ""
        , spaces = spaces
        , item = lazy (\_ -> parseInstruction)
        , trailing = Parser.Forbidden
        }


{-| Parse une instruction unique parmi les commandes disponibles
-}
parseInstruction : Parser Instruction
parseInstruction =
    oneOf
        [ parseForward
        , parseLeft
        , parseRight
        , parseRepeat
        ]


{-| Parse un programme complet entouré de crochets
Ex: "[Forward 100, Repeat 4 [Forward 50, Left 90]]"
-}
parseProgram : Parser (List Instruction)
parseProgram =
    succeed identity
        |. symbol "["
        |. spaces
        |= parseInstructions
        |. spaces
        |. symbol "]"
        |. end


{-
Transforme une chaîne de commandes en liste de points à dessiner.

- Commence à la position (0,0)
- Inclut le point de départ dans le résultat
- Retourne une liste vide en cas d'erreur de parsing
-}
read : String -> List (Float, Float)
read input =
    case Parser.run parseProgram input of
        Ok instructions ->
            let
                initTortue = { x = 0, y = 0, angle = 0 }
                (_, points) = appliquerInstructions instructions initTortue
            in
            List.map (\p -> (p.x, p.y)) ({ x = 0, y = 0 } :: points)

        Err _ ->
            []
