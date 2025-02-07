const { Joueur, rl } = require('./joueur');
const { choisirMotAleatoire } = require('./dico');
const { writeIndice , writeTour, writeFinDeTour } = require('./log');
const fs = require('fs')


const nombreJoueurs = 5;
const Joueurs = [];

for (let i = 0; i < nombreJoueurs; i++) {
    Joueurs.push(new Joueur());
}

function demanderNom(i = 0) {
    if (i < nombreJoueurs) {
        Joueurs[i].askNom(i + 1, () => {
            demanderNom(i + 1);
        });
    } else {
        commencerJeu(0, 13, 0);
    }
}

function afficherMessageScore(score) {
    if (score === 13) {
        console.log("13 Score parfait ! Y arriverez-vous encore ?");
    } else if (score === 12) {
        console.log("12 Incroyable ! Vos amis doivent être impressionnés !");
    } else if (score === 11) {
        console.log("11 Génial ! C’est un score qui se fête !");
    } else if (score >= 9 && score <= 10) {
        console.log("9-10 Waouh, pas mal du tout !");
    } else if (score >= 7 && score <= 8) {
        console.log("7-8 Vous êtes dans la moyenne. Arriverez-vous à faire mieux ?");
    } else if (score >= 4 && score <= 6) {
        console.log("4-6 C’est un bon début. Réessayez !");
    } else if (score >= 0 && score <= 3) {
        console.log("0-3 Essayez encore");
    }
}

function commencerJeu(tour = 0, nb_manche, score) {
    if (tour === 0) {
        fs.unlink('./log.txt', (err) => {
            if (err) {
                console.error(err)
                return
            }
        });
    }
    if (tour >= nb_manche) {
        console.log("Fin du jeu !");
        afficherMessageScore(score);  // Affiche le message en fonction du score
        rl.close();
        return;
    }

    choisirMotAleatoire((motATrouver) => {
        console.log(`\n Mot à deviner : ${motATrouver}\n`);

        const joueurADeviner = Joueurs[tour % nombreJoueurs];
        console.log(`${joueurADeviner.getNom()} doit deviner le mot !`);

        const indices = [];
        writeTour(motATrouver, joueurADeviner.getNom())
        function demanderIndice(i = 0) {
            if (i < nombreJoueurs) {
                if (i === tour%nombreJoueurs) {
                    demanderIndice(i + 1);
                } else {
                    nom = Joueurs[i].getNom();
                    rl.question(`Indice de ${nom}: `, (answer) => {
                        indices.push(answer);
                        writeIndice(answer,nom);
                        demanderIndice(i + 1);
                    });
                }
            } else {
                traiterIndices();
            }
        }
        
        function traiterIndices() {
            const uniqueIndices = [];
            for (let i = 0; i < indices.length; i++) {
                let isUnique = true;
                for (let j = 0; j < indices.length; j++) {
                    if (i !== j && indices[i] === indices[j]) {
                        isUnique = false;
                        break;
                    }
                    if (indices[i] === motATrouver) {
                        isUnique = false;
                        break;
                    }
                }
                if (isUnique) {
                    uniqueIndices.push(indices[i]);
                }
            }

            console.log("\n Indices retenus :");
            if (uniqueIndices.length === 0) {
                console.log("Tous les indices ont été annulés !");
            } else {
                console.log(uniqueIndices.join(", "));
            }

            rl.question(`\n${joueurADeviner.getNom()}, quelle est ta réponse ? `, (reponse) => {
                if (reponse === motATrouver) {
                    console.log("Bonne réponse !");
                    score = score + 1;
                } else {
                    console.log(`Mauvaise réponse ! Le mot était : ${motATrouver}`);
                    nb_manche = nb_manche - 1;
                }
                writeFinDeTour(joueurADeviner.getNom(), reponse);
                console.log(`tour : ${tour}`);
                console.log(`nb_manche : ${nb_manche}`);
                if (tour === nb_manche && score > 0) {
                    score = score -  1;
                }
                console.log(`Score actuel : ${score}`);
                console.log(`Cartes restantes : ${nb_manche - tour - 1}`);
                commencerJeu(tour + 1, nb_manche, score);
            });
        }

        demanderIndice();
    });
}

demanderNom();
