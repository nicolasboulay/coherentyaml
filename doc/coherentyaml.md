# specifications

Le but est de faire une sorte de schema composable pour yaml/json.

## Design

- Un document YAML sert de schéma « par l’exemple ».
- Deux documents sont comparés par une relation de cohérence.
- Cette relation doit être :
  - symétrique ;
  - réflexive : tout document est cohérent avec lui-même.
- Les types sont exprimés par des valeurs neutres : 1, -1, 1.0, "", etc.
- Trois opérateurs logiques sont prévus :
  - Coherent : conjonction ;
  - OR : alternative ;
  - Not : négation.
- Les objets sont comparés champ par champ.
- Les tableaux sont non ordonnés


## Question à trancher

- Les clés supplémentaires sont-elles permises ?
  - non sauf mot clef comme '*'
- Quelle sémantique précise donner à Not ?
- Les tableaux sont-ils des ensembles, des multi-ensembles ou des listes ?
- Comment traiter null ?
- Comment agréger plus de deux fichiers ?
  - le faire 2 par 2.
- Comment intégrer les futures regex tout en conservant symétrie et composition ?
  - trouver un exemple pour que 2 regexp soit compatible me semble un peu lourd