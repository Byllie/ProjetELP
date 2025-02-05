var prompt = require('prompt');

class Joueur {
    constructor() {
        this.nom = "Joueur";
        this.word = null;
        this.score = 0;
    }
    getNom() {
        return this.nom;
    }
    setNom(nom) {
        this.nom = nom;
    }
    getScore() {
        return this.score;
    }
    setScore(score) {
        this.score = score;
    }
    incrementerScore() {
        this.score++;
    }
    getWord() {
        return this.word;
    }
    setWord(word) {
        this.word = word;
    }
    askWord() {
        prompt.start();
        prompt.get(['word'], (err, result) => {
            if (err) {
                console.log(err);
            } else {
                this.setWord(result.word);
            }
        });

    }
}


module.exports = { Joueur };