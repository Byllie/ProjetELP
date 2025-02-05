module Main exposing (..)
import Browser
import Html exposing (Html, div, input, button, text, textarea, h1, h2, pre, p)
import Html.Attributes exposing (placeholder, class, style)
import Html.Events exposing (onClick, onInput)
import ParseTcTurtle exposing (read)
import DrawTcTurtle exposing (display)


{-| Modèle de l'application.
- text: Stocke le texte saisi dans la zone de commande
- result: Contient la liste des points calculés à dessiner
-}
type alias Model =
    { text : String
    , result : List (Float, Float)
    }


{-| Initialise le modèle avec des valeurs par défaut.-}
init : Model
init =
    { text = ""
    , result = []
    }


{-| Messages possibles dans l'application.
- UpdateText: Met à jour le texte saisi (avec validation en temps réel)
- ParseCommand: Déclenche l'analyse et le calcul du dessin
-}
type Msg
    = UpdateText String
    | ParseCommand


update : Msg -> Model -> Model
update msg model =
    case msg of
        -- Met à jour le texte brut des commandes
        UpdateText newText ->
            { model | text = newText }

        -- Déclenche l'analyse syntaxique et le calcul du tracé
        ParseCommand ->
            { model | result = read model.text }


{-| Génère la vue de l'application.-}
view : Model -> Html Msg
view model =
    div [ class "container" ]
        [ h1 [] [ text "TC Tortu(r)e" ]
        , div [] [ text "Tapez votre code ci-dessous :" ]
        , textarea
            [ placeholder "exemple : [Repeat 360 [Forward 1, Left 1]]"
            , Html.Attributes.value model.text
            , onInput updateTextInput
            ]
            []
        , button
            [ onClick ParseCommand ]
            [ text "Dessiner" ]
        , display model.result
        , div [ class "examples", style "margin-top" "20px" ]
            [ h2 [] [ text "Exemples de code :" ]
            , pre []
                [ text """[Repeat 360 [ Right 1, Forward 1]]
[Forward 100, Repeat 4 [Forward 50, Left 90], Forward 100]
[Repeat 36 [Right 10, Repeat 8 [Forward 25, Left 45]]]
[Repeat 8 [Left 45, Repeat 6 [Repeat 90 [Forward 1, Left 2], Left 90]]]"""
                ]
            , p []
                [ text "Amélioration par rapport au sujet : Le dessin ne peut pas sortir du cadre (taille variable) et il occupe un maximum de place." ]
            ]
        ]


updateTextInput : String -> Msg
updateTextInput value =
    UpdateText value


-- Point d'entrée
main =
    Browser.sandbox
        { init = init
        , update = update
        , view = view
        }
