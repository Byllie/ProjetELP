module Main exposing (..)

import Browser
import Html exposing (Html, div, input, button, text)
import Html.Attributes exposing (placeholder)
import Html.Events exposing (onClick, onInput)
import ParseTcTurtle exposing (parser)
import DrawTcTurtle exposing (draw)


type alias Model =
    { text : String
    , result : List (Float, Float)
    }

init : Model
init =
    { text = ""
    , result = []
    }


type Msg
    = UpdateText String
    | ParseCommand


update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateText newText ->
            { model | text = newText }

        ParseCommand ->
            let
                parsedResult = parser model.text
            in
            { model | result = parsedResult }


view : Model -> Html Msg
view model =
    div []
        [ div [] [ text "Tapez votre code ci-dessous :" ]
        , input
            [ placeholder "exemple : [Repeat 360 [Forward 1, Left 1]]"
            , Html.Attributes.value model.text
            , onInput updateTextInput
            ] []
        , button [ onClick ParseCommand ] [ text "Dessiner" ]
        , div [] [ text ("Commande entrée : " ++ model.text) ]
        , div [] [ text ("Résultat du parsing : " ++ (draw model.result)) ]
        ]

updateTextInput : String -> Msg
updateTextInput value =
    UpdateText value


main =
    Browser.sandbox
        { init = init
        , update = update
        , view = view
        }
