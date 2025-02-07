const readline = require('readline');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

class Joueur {
    constructor() {
        this.nom = "Joueur";
    }

    getNom() {
        return this.nom;
    }

    setNom(nom) {
        this.nom = nom;
    }

    askNom(numero, callback) {
        rl.question(`Quelle est le nom de joueur ${numero} ? `, (answer) => {
            this.setNom(answer);
            callback();
        });
    }
}

module.exports = { Joueur, rl };
