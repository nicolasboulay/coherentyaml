# specifications

Le but est de faire une sorte de schema composable pour yaml/json.

## Design

- Un document YAML sert de schéma « par l’exemple ».
- Chaque document décrit une contrainte sur les valeurs possibles.
- Deux documents sont cohérents lorsqu'au moins une même valeur satisfait
  leurs deux contraintes. Cette relation est symétrique.
- Tout document satisfiable est cohérent avec lui-même. Un document
  contradictoire ne l'est pas, y compris avec lui-même.
- Les types sont exprimés par des valeurs neutres : 1, -1, 1.0, "", etc.
- Les mêmes mots clefs servent à exprimer la logique et les opérations
  sur les ensembles de valeurs : `Coherent`, `OR` et `Not`.
- Les objets sont comparés champ par champ.
- Les tableaux sont non ordonnés


## Logique et ensembles de valeurs

On note `U` l'univers commun des valeurs de documents considérées et `S(A)`
l'ensemble des valeurs qui satisfont une expression `A`. Dire qu'une valeur
`v` satisfait `A` revient à dire que `v` appartient à `S(A)`.

| Expression | Sens logique pour une valeur `v` | Ensemble représenté |
| --- | --- | --- |
| `Coherent(A, B, ...)` | `A(v) AND B(v) AND ...` | Intersection des ensembles |
| `OR(A, B, ...)` | `A(v) OR B(v) OR ...` | Union des ensembles |
| `Not(A)` | `NOT A(v)` | Complément `U ∖ S(A)` |

Ces deux lectures décrivent la même sémantique et se composent librement,
quel que soit le nombre d'opérandes ou leur imbrication. Une chaîne exacte
comme `a` représente le singleton `{a}` ; un motif, notamment une valeur
neutre de type, peut représenter plusieurs valeurs.

Le complément utilise le même univers `U` dans toute l'expression : il ne
se limite pas implicitement aux valeurs citées ou au type de l'opérande.
Pour exprimer un complément dans un groupe `G`, on écrit
`Coherent(G, Not(A))`, soit la différence `S(G) ∖ S(A)`.

Les constantes logiques internes `yes` et `no` représentent respectivement
`U` et l'ensemble vide. Elles ne doivent pas être confondues avec les valeurs
booléennes littérales YAML. Une conjonction vide représente `U` ; une union
vide représente l'ensemble vide.

Les lois usuelles s'appliquent : commutativité, associativité et idempotence
de l'intersection et de l'union, distributivité, double négation et lois de
De Morgan. En particulier :

- `Not(Not(A))` représente le même ensemble que `A`.
- `Not(OR(A, B))` équivaut à `Coherent(Not(A), Not(B))`.
- `Not(Coherent(A, B))` équivaut à `OR(Not(A), Not(B))`.

## Cohérence globale

`A.IsCoherent()` réussit si `S(A)` est non vide.
`A.IsCoherentWith(B)` réussit si `S(A)` et `S(B)` ont une intersection non
vide. Dans l'API Go, la réussite est représentée par `nil`, l'incohérence
par une erreur.

La négation porte sur les valeurs acceptées, pas sur le résultat du test de
cohérence : `A` et `Not(A)` peuvent tous deux être satisfiables, tout en
étant incompatibles entre eux.

Pour agréger plusieurs documents, il faut une valeur commune à **toutes**
leurs contraintes. Vérifier chaque paire indépendamment ne suffit pas.
Une agrégation progressive est possible à condition de conserver
l'intersection accumulée, et non un simple résultat booléen.

Pour trois chaînes exactes distinctes `a`, `b` et `c` :

- `Coherent(OR(a, b), Not(a), Not(b))` est incohérent : toutes les
  alternatives sont exclues.
- `Coherent(OR(a, b), OR(b, c), OR(a, c))` est incohérent : chaque paire
  a une intersection non vide, mais l'intersection des trois est vide.
- `Coherent(OR(a, b, c), Not(a), Not(b))` est cohérent : `c` est une
  valeur commune.

L'ordre des contraintes, leur regroupement et le sens de comparaison ne
doivent pas changer le résultat. L'évaluation d'une alternative de `OR`
doit tenir compte de toutes les contraintes de la conjonction qui l'entoure.

## Points atomiques à préciser

Le contrat logique ci-dessus est fixé. Les ensembles exacts représentés par
certains motifs restent à préciser avant d'annoncer leur prise en charge
complète, notamment leur complément. Ces précisions doivent respecter le
même modèle logique et ensembliste.

- Les clés supplémentaires sont-elles permises ?
  - non sauf mot clef comme '*'
- Les tableaux sont-ils des ensembles, des multi-ensembles ou des listes ?
- Comment traiter null ?
- Comment intégrer les futures regex tout en conservant symétrie et composition ?
  - trouver un exemple pour que 2 regexp soit compatible me semble un peu lourd
