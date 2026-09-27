# Plan d’implémentation de Not logique

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Évaluer globalement `Coherent` (AND / intersection), `OR` (OR inclusif / union) et `Not` (négation / complément), avec une cohérence symétrique fondée sur l’existence d’une valeur commune.

**Architecture:** Conserver les types et l’interface `Node`. Centraliser l’évaluation des expressions logiques dans une fonction interne qui conserve toutes les contraintes d’une conjonction pendant l’exploration des alternatives. Limiter les modifications de `node.go` aux points d’entrée nécessaires et isoler l’algorithme dans un fichier.

**Tech Stack:** Go 1.23 selon `go.mod`, tests Go existants ; aucune dépendance ajoutée.

**Spec:** [Spécification de coherentyaml](../../../doc/coherentyaml.md), notamment « Logique et ensembles de valeurs » et « Cohérence globale ». Ce document prépare l’implémentation ; il ne l’exécute pas.

## Contrat et contraintes

- `Coherent(a, b)` représente `a AND b` ; `OR(a, b)` représente `a OR b`.
- `Not(a)` représente le complément des valeurs satisfaisant `a`.
- Tous les compléments sont calculés dans un même univers `U`, sans restriction implicite au type de l'opérande ni aux valeurs citées. `Coherent(G, Not(A))` représente la différence entre les ensembles de `G` et de `A`.
- `a.IsCoherentWith(b)` réussit si une même valeur satisfait les deux expressions.
- `IsCoherent()` recherche une solution, sans contrainte supplémentaire (`yes`).
- `yes` n’impose aucune contrainte ; `no` n’a aucune solution.
- `yes` et `no` sont les constantes logiques internes représentant `U` et l'ensemble vide ; les valeurs booléennes littérales `true` et `false` sont deux valeurs distinctes, jamais des alias de ces constantes.
- Une conjonction vide représente `U` ; une union vide représente l'ensemble vide.
- La réflexivité vaut pour les documents satisfiables : un document contradictoire est incohérent avec lui-même.
- Respecter la commutativité, l'associativité, l'idempotence et la distributivité des intersections et unions. Agréger progressivement des documents exige de conserver l'intersection accumulée ; des résultats booléens par paires ne suffisent pas.
- `nil` signifie cohérent ; une erreur signifie incohérent.
- Respecter la symétrie, la double négation et De Morgan. Ne pas inverser le booléen de `a.IsCoherentWith(b)` pour évaluer `Not(a)`.
- Conserver les attentes actuelles de `TestNot`, sans les adapter aux défauts de l’implémentation.
- Conserver les attentes de `TestCoherentGlobalConstraints`, y compris les cas sans `Not`, les témoins positifs, les six permutations et les comparaisons regroupées ou inversées.
- Préserver le format YAML, les constructeurs, les signatures publiques et les modifications utilisateur présentes.
- Ne pas mélanger ce changement avec une réorganisation générale, une mise à jour de dépendances ou la correction de défauts indépendants.

## Limite à lever avant d’annoncer une correction générale

Les tests discutés fixent le contrat pour les constantes logiques et les chaînes exactes. Les feuilles neutres (`StrZero`, nombres sentinelles), les structures partielles et les tableaux ne sont pas nécessairement des valeurs exactes : leur complément demande une définition de leur domaine et de leur inclusion.

Une simple incompatibilité entre deux motifs ne permet pas de calculer leur différence. Ne pas traiter les structures comme des variables booléennes indépendantes, ni utiliser `!IsCoherentWith` comme solution de repli silencieuse. La tâche 1 détermine si leur sémantique existante suffit ; sinon, consigner les exemples indécidables et faire préciser uniquement ces points avant leur implémentation. Une première livraison limitée aux expressions sur chaînes doit être annoncée comme telle, jamais comme une prise en charge complète de `Not`.

## Fichiers concernés

- `internal/node/node_test.go` : contrat existant et tests d’intégration des opérateurs.
- `internal/node/logic.go` (nouveau) : évaluateur interne et traitement des contraintes.
- `internal/node/logic_test.go` (nouveau) : tests ciblés de l’évaluation globale.
- `internal/node/node.go` : délégation depuis les opérateurs ; ajustements de dispatch démontrés nécessaires par les tests.
- `cmd/coherentyaml/coherentyaml_test.go` : validation de l’entrée YAML, uniquement si un cas manquant doit être ajouté.
- `doc/coherentyaml.md` : courte description de la sémantique retenue et de son périmètre réel.

## Points de revue

1. Conjonction globalement impossible mais dont chaque paire est compatible : tâche 2.
2. Négation d’une disjonction et disjonction contenant une négation : tâche 2.
3. Deux contraintes négatives, sans valeur positive imposée : tâche 2.
4. Comparaison inversée ou passage par `Coherent` donnant un résultat différent : tâche 3.
5. Motifs neutres, structures et tableaux : tâche 1, puis tâche 3 selon le contrat établi.

## Tâche 1 — Établir la référence et le périmètre atomique

- [ ] Exécuter `go test ./internal/node -run '^TestNot$' -v`, puis `go test ./...`. Conserver les noms des échecs et leurs causes avant modification.
- [ ] Exécuter `go test ./internal/node -run '^TestCoherentGlobalConstraints$' -v`. Les tests sont déjà ajoutés : relever les échecs des conjonctions impossibles et vérifier les témoins positifs et les paires, sans dupliquer ces tests.
- [ ] Relire `Leaf.IsCoherentWith`, `isNeutral`, `NStruct.IsCoherentWith` et `NArray.IsCoherentWith`, ainsi que leurs tests. Décrire les ensembles représentés, notamment les valeurs neutres et les champs absents.
- [ ] Relever les incohérences de dispatch : `NStruct.IsCoherentWith` appelle actuellement `IsCoherent()` sur un opérateur et perd ainsi la contrainte de structure. Ne pas considérer cette délégation comme une implémentation de la symétrie.
- [ ] Fixer les cas attendus pour les motifs non exacts avant leur codage ; si le contrat ne peut pas être déduit, documenter la décision manquante. Ne pas inventer une sémantique pour obtenir des tests verts.

**Livrable vérifiable :** liste de référence des échecs et contrat des atomes pris en charge. Aucun changement de production à cette étape.

## Tâche 2 — Évaluer une conjonction complète

**Interface interne prévue :** `logicalCoherence(nodes ...Node) error`, dans `logic.go`. Elle évalue la satisfiabilité du AND des arguments ; elle ne rappelle pas les méthodes publiques des opérateurs pour résoudre la même expression.

### Étape 2a — Intersection et union, sans négation

- [ ] Ajouter `TestLogicalCoherenceGlobal` dans `logic_test.go` avec des chaînes exactes distinctes : `AND(OR(a,b), OR(b,c), OR(a,c))` incohérent ; `AND(OR(a,c), OR(b,c), OR(a,b,c))` cohérent, avec `c` comme témoin. Tester les six permutations et les regroupements imbriqués.
- [ ] Ajouter `TestLogicalCoherenceIdentities` : aucun argument cohérent ; AND vide cohérent ; OR vide incohérent ; `AND(a,yes)` cohérent ; `AND(a,no)` incohérent ; `OR(a,no)` cohérent. Ajouter la réflexivité de `a` et la non-réflexivité de `AND(a,b)` pour `a != b`.
- [ ] Exécuter `go test ./internal/node -run '^TestLogicalCoherence(Global|Identities)$' -v` et constater les échecs avant implémentation. Si le symbole manque, ajouter seulement sa déclaration minimale pour obtenir des échecs d'assertion.
- [ ] Implémenter le parcours global dans `logic.go` : accumuler les contraintes AND ; explorer chaque alternative OR avec les contraintes restantes ; isoler l'état des branches. Prévoir l'ajout de la polarité de négation dans ce même évaluateur.
- [ ] Relancer la commande précédente : ces tests doivent passer avant l'étape 2b. Les tests publics existants peuvent encore échouer tant que le dispatch n'est pas branché en tâche 3.

### Étape 2b — Complément dans le même évaluateur

- [ ] Ajouter des tests directs à attentes littérales dans `logic_test.go` : `Not(s1)` vrai ; `Not(Not(s1)) AND s1` vrai ; `Not(s1) AND s1` faux ; `Not(s1) AND OR(s1,s2)` vrai ; `Not(OR(s1,s2)) AND OR(s2,s3)` vrai.
- [ ] Ajouter les cas qui imposent une recherche globale : `OR(s1,s2) AND Not(s1) AND Not(s2)` faux ; `Not(s1) AND Not(s2)` vrai, avec `s3` comme témoin ; `Not(OR(s1,Not(s1)))` faux.
- [ ] Tester De Morgan avec des résultats explicites pour `s1`, `s2`, `s3` : `Not(OR(s1,s2))` et `AND(Not(s1),Not(s2))` refusent les deux premières valeurs et acceptent la troisième. `Not(AND(s1,s2))` et `OR(Not(s1),Not(s2))` acceptent les trois.
- [ ] Implémenter une exploration récursive avec une polarité de négation : inverser `yes/no`, éliminer les doubles négations, échanger AND/OR sous négation. Un AND accumule ses contraintes ; un OR explore une alternative avec toutes les contraintes restantes. Ne pas construire systématiquement une forme normale complète, qui peut grossir exponentiellement.
- [ ] Résoudre les contraintes atomiques d’une branche selon le contrat de la tâche 1. Pour les chaînes exactes : deux valeurs positives distinctes sont incompatibles ; une valeur imposée ne doit pas être exclue ; un ensemble fini d’exclusions seul laisse des chaînes possibles. Un domaine fini n'est épuisé par des exclusions que si une contrainte positive impose ce domaine ; ne jamais restreindre implicitement `U` aux booléens.
- [ ] Ajouter `TestLogicalCoherenceBooleanValues` avec des feuilles booléennes littérales distinctes de `yes/no` : `AND(true,false)` incohérent ; `AND(Not(true),false)` cohérent ; `AND(Not(true),"bonjour")` cohérent ; `AND(Not(true),Not(false))` cohérent ; `AND(OR(true,false),Not(true),Not(false))` incohérent. Vérifier aussi `Not(yes)` incohérent et `Not(no)` cohérent.
- [ ] Compléter les identités avec `Not(AND())` incohérent et `Not(OR())` cohérent. Tester la différence relative : `AND(OR(a,b),Not(a))` accepte `b`, refuse `a` et refuse `c`.
- [ ] Ajouter `TestLogicalCoherenceLaws` : comparer les formes commutées, associées, répétées et distribuées sur les témoins exacts `a`, `b`, `c`, avec des attentes littérales pour chaque forme. Par exemple `AND(OR(a,b),OR(b,c))` et sa forme distribuée acceptent seulement `b` parmi ces témoins. Tester les deux distributivités ; ne pas se contenter de comparer deux résultats produits par l'implémentation.
- [ ] Pour chaque groupe de tests de l'étape 2b, constater les échecs avant le code correspondant, puis les rendre verts. Exécuter `go test ./internal/node -run '^TestLogicalCoherence' -v` pour valider ensemble les étapes 2a et 2b.
- [ ] Exécuter les tests de l’évaluateur. Vérifier que permuter les contraintes ne change aucun résultat et que l’exploration d’une branche ne modifie pas les autres.

**Livrable vérifiable :** évaluateur testé indépendamment du dispatch existant ; aucune dépendance externe ni modification de l’AST.

## Tâche 3 — Brancher les opérateurs et valider le résultat

- [ ] Remplacer le corps de `Not.IsCoherentWith(o Node) error` par une délégation à l’évaluateur sur `n` et `o`. Conserver `Not.IsCoherent()` comme évaluation avec `yes` ; supprimer uniquement les anciennes branches et commentaires devenus faux dans ce bloc.
- [ ] Faire passer `Coherent.IsCoherentWith` par la même évaluation globale : le contrôle par paires doit disparaître pour ces expressions.
- [ ] Faire passer `OR.IsCoherentWith` par le même évaluateur afin que la branche choisie conserve le contexte de la conjonction. Vérifier `Yes`, `No` et les délégations des feuilles ; éviter les modifications là où les méthodes existantes satisfont déjà le contrat.
- [ ] Pour chaque ajustement nécessaire aux structures/tableaux, ajouter d’abord un test de comportement précis issu de la tâche 1. Garder ces ajustements dans une modification séparée à relire ; aucun changement opportuniste aux règles de champs ou de tableaux.
- [ ] Ajouter aux tests publics de `node_test.go` les cas booléens, identités vides et réflexivité de la tâche 2 : vérifier `IsCoherent`, `IsCoherentWith` dans les deux sens et le passage par `Coherent`. Cela doit détecter un dispatch qui contourne l'évaluateur.
- [ ] Exécuter `go test ./internal/node -run '^(TestNot|TestCoherentGlobalConstraints|TestLogicalCoherence.*)$' -v`, puis `go test ./internal/node` et `go test ./...`. Les nouvelles attentes doivent passer sans suppression ni assouplissement.
- [ ] Comparer chaque échec de la suite complète à la référence. Distinguer les défauts indépendants des régressions introduites ; corriger ces dernières et rapporter explicitement les autres, notamment `TestShallMatch`, `TestCalculDeProposition2`, `TestCalculDeProposition` s’ils persistent.
- [ ] Ajouter un test d’intégration YAML pour `Not` dans une disjonction si les tests existants ne couvrent pas ce chemin, et vérifier que la sérialisation reste identique.
- [ ] Vérifier la conformité à `doc/coherentyaml.md`, qui décrit déjà la sémantique cible. Documenter séparément le périmètre réellement implémenté et les limites restantes, sans affaiblir le contrat pour refléter un défaut du code. Exécuter `gofmt` et `git diff --check` sur les fichiers touchés.

## Découpage pour la relecture

Présenter trois ensembles de modifications distincts : tests et contrat ; évaluateur interne ; branchement et documentation. Les corrections requises pour les motifs structurés restent identifiables séparément. Ne pas reformater tout `node.go`, renommer les types ou modifier les diagnostics sans nécessité.

**Critère de fin :** tests du contrat verts, symétrie et recherche d’une solution commune vérifiées, aucune régression nouvelle non résolue, et périmètre des motifs pris en charge explicitement documenté. Des tests limités aux chaînes ne suffisent pas à annoncer une correction générale sur tous les `Node`.
