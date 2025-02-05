

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
        this.word = prompt("Entrez un mot");
    }
}


module.exports = { Joueur };