const exp = require('constants');
const fs = require('fs');


async function writeIndice(mot, nom){
    try {
        await fs.writeFile('./log.txt',"\t" + nom + " a proposé : " + mot + "\n",{ flag: 'a+' }, err => {});
    } catch (err) {
        console.log(err);
    } 
}

async function writeTour(mot, joueur) {
    try {
        await fs.writeFile('./log.txt', joueur + " doit deviné : " + mot + "\n",{ flag: 'a+' }, err => {});
    } catch (err) {
        console.log(err);
    } 
    
}
async function writeFinDeTour(mot, joueur) {
    try {
        await fs.writeFile('./log.txt',"\t" + joueur + " a deviné : " + mot + "\n\n",{ flag: 'a+' }, err => {});
    }
    catch (err) {
        console.log(err);
    }
}

module.exports = { writeIndice, writeTour, writeFinDeTour };