const _ = require('lodash');
const fs = require('fs');

function choisirMotAleatoire(callback) {
    fs.readFile('dico.txt', 'utf8', (err, data) => {
        if (err) {
            console.log("Erreur de lecture du fichier :", err);
            return;
        }

        const lignes = data.split('\n');
        const motATrouver = _.sample(lignes);
        callback(motATrouver);
    });
}

module.exports = {choisirMotAleatoire}
