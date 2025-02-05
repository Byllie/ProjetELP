
// var Client = require('./client');
// var Server = require('./server');
var Joueur = require('./joueur');

// const HOST = 'localhost';
// const PORT = 3080;
console.log(this);
nombreJoueurs = 5;
const Joueurs = [];
for (let i = 0; i < nombreJoueurs; i++) {
    Joueurs.push(new Joueur.Joueur());
}

for (let i = 0; i < nombreJoueurs; i++) {
    Joueurs[i].setNom(`Joueur ${i + 1}`);
    console.log(Joueurs[i].getNom());
}

for (let i = 0; i < nombreJoueurs; i++) {
    console.log("Joueur " + Joueurs[i].getNom());
    Joueurs[i].askWord();
    console.log(Joueurs[i].getWord());
    
}


