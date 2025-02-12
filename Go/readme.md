# Projet Go : Implémentation parallèle de recherche de communauté

Notre projet est basé sur le travail de recherche de Arnau Prat-Pérez, David Dominguez-Sal et Josep-LLuis Larriba-Pey
"High Quality, Scalable and Parallel Community Detection
for Large Real Graph"
https://dl-acm-org.docelec.insa-lyon.fr/doi/pdf/10.1145/2566486.2568010

Nous avons mis en place un serveur TCP avec des goroutines qui récupère le ou les graphs a traités puis répartie le calcul de BestMouvement dans un pool de worker (1 pool par connections TCP).
On renvoie ensuite la liste des communautés via TCP.

### **Instructions d'exécution:**


#### Fichier à télécharger (et où les mettre)


*Si vous voulez utilisez notre graph de test :*
- [amazon-meta.txt](https://snap.stanford.edu/data/amazon-meta.html)
    - Dans go_client et sans changer le nom
- [com-amazon.ungraph.com](https://snap.stanford.edu/data/com-Amazon.html)
    - Dans go_client avec le nom **com-amazon.com** 

#### Exécution

lancer graph.go
puis envoyé un graph bien structuré via TCP : nc localhost 5828 < input.txt > output.txt
Dataset de graph: http://snap.stanford.edu/data/index.html#communities prendre [_Networks with ground-truth communities_](http://snap.stanford.edu/data/index.html#communities) undirected, le graph d'amazon à la bonne taille pour notre programme)

Nous avons fait un programme secondaire qui a partir de l'output du programme principal pour le **graph d'amazon,** fait correspondre l'id des produits avec leur nom afin de pouvoir vérifier la pertinance des communautés (utilisez send.sh pour le graph Amazon)


### Fichier de résultat 

Normalement une fois le programme fini un fichier ***communities.txt*** à été creé (ou modifié si déjà existant) dedans vous retrouverez :
- les communautés avec un nombre pour les identifier et le nombre de sommet dedans 
- les sommets appartennant à la communauté sous celle-ci
