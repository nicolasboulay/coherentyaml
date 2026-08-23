# 20260823

## État du projet : analyse par codex

  C’est un prototype Go d’environ 2 000 lignes, destiné à comparer des documents YAML avec un schéma « par l’exemple » et les opérateurs Coherent, OR et Not.

  - 35 commits, dernier développement fonctionnel en mars 2021, fusion finale en janvier 2023.
  - Licence GPL-3.0.
  - Aucun tag, aucune release, aucune CI et aucun outil de qualité configuré.
  - Le binaire compile correctement avec Go 1.27.
  - Le projet déclare encore Go 1.13 et goccy/go-yaml 1.8.9.
  - Les 10 fichiers Go ne sont pas conformes à gofmt.
  - go vet relève 8 problèmes, dont deux blocs inaccessibles.
  - Les quatre fichiers Docker locaux sont non suivis et n’installent actuellement pas Go ; ils ne constituent donc pas encore un environnement de développement reproductible.

  Le dépôt public confirme l’absence de release, d’issue et de PR active : coherentyaml sur GitHub (https://github.com/nicolasboulay/coherentyaml).

  ## Santé fonctionnelle

  Le build passe, mais les tests échouent :

  - TestShallMatch : conflit sur la présence de clés supplémentaires.
  - TestCalculDeProposition2 : échec autour des lois de De Morgan.
  - TestCalculDeProposition : plusieurs tautologies échouent, principalement autour de Not.

  Le paquet moteur internal/node passe seul avec 63,3 % de couverture. Le paquet CLI échoue avec 18,6 % de couverture.

  Les problèmes les plus importants sont :

  - Le README annonce déjà que Not n’est pas fonctionnel : README.md:12.
  - NStruct.IsCoherent calcule une erreur mais retourne toujours nil : internal/node/node.go:528.
  - La comparaison des tableaux prétend être symétrique, mais ne vérifie qu’un sens.
  - No.String() renvoie "true".
  - OR est lu avec la clé OR mais réécrit sous la forme Or.
  - Les erreurs YAML appellent log.Fatal, ce qui rend le moteur difficile à tester et réutiliser : cmd/coherentyaml/ast.go:12.
  - Un fichier absent provoque finalement une panic au lieu d’une erreur propre : cmd/coherentyaml/coherentyaml.go:71.
  - La CLI annonce plusieurs fichiers, mais ignore tout argument après le deuxième.
  - Certains tests contiennent des break qui rendent une partie du code inaccessible et limitent fortement les cas réellement testés.

  ## Modernisation recommandée

  Je procéderais dans cet ordre :

  1. Créer modernisation depuis master, puis supprimer les deux anciennes branches.
  2. Fixer le contrat fonctionnel avant de refactorer :
      - clés supplémentaires autorisées ou interdites ;
      - tableaux ordonnés ou non, avec ou sans multiplicité ;
      - signification exacte de Not ;
      - gestion de null ;
      - comparaison de deux fichiers ou agrégation de plusieurs fichiers.

  3. Remettre les tests à plat et obtenir une suite verte sur le comportement décidé.
  4. Séparer le moteur de la CLI, par exemple internal/schema et cmd/coherentyaml.
  5. Remplacer les sorties fatales par des erreurs propagées, rendre l’affichage déterministe et définir les codes de sortie.
  6. Passer au minimum à Go 1.26, tester aussi Go 1.27 en CI, et mettre goccy/go-yaml à jour.
  7. Ajouter gofmt, go vet, tests, couverture et analyse de vulnérabilités dans GitHub Actions.
  8. Réécrire le README avec une syntaxe de schéma documentée et des exemples exécutables.

  Go 1.27.0 est la version stable actuelle, publiée le 19 août 2026 ; Go 1.26.7 constitue une base plus éprouvée pour le minimum supporté (historique officiel Go (https://go.dev/doc/devel/release)). goccy/go-yaml est actuellement en 1.19.2 (releases officielles (https://github.com/goccy/go-yaml/releases)). J’ai testé cette montée de version dans une copie temporaire : elle compile et conserve exactement les trois échecs
  existants, sans nouvelle régression détectée.

# 20260824

### Déjà fait ou partiellement fait

- L’image Docker utilise maintenant Go 1.27 et les fichiers Docker sont suivis par Git.
- Les commandes go test, go vet et couverture ont été ajoutées à TEST.
- Le moteur est déjà isolé dans internal/node, distinct de la CLI.
- Un hook gofmt existe pour les futurs commits.

### Priorité haute — toujours à faire

1. Définir le contrat fonctionnel avant correction :
    - politique des clés supplémentaires ;
    - ordre et multiplicité des tableaux ;
    - sémantique exacte de Not ;
    - comportement de null ;
    - prise en charge réelle de plusieurs fichiers.

2. Réparer et compléter les tests :
    - les trois groupes d’échecs consignés n’ont pas été corrigés ;
    - deux break limitent encore fortement les tests ;
    - les commandes ajoutées à TEST automatisent les contrôles, mais ne les rendent pas verts.

3. Corriger les défauts moteur toujours présents :
    - NStruct.IsCoherent() ignore encore l’erreur calculée et retourne nil dans internal/node/node.go:528 ;
    - les tableaux ne sont comparés que dans un sens dans internal/node/node.go:642 ;
    - leur multiplicité reste explicitement indéfinie ;
    - No.String() retourne toujours "true" ;
    - OR est toujours sérialisé sous Or ;
    - Not demeure non finalisé et le README l’annonce toujours non fonctionnel.

4. Fiabiliser la CLI :
    - log.Fatal est toujours employé pendant le parsing YAML dans cmd/coherentyaml/ast.go:12 ;
    - une erreur de lecture est affichée mais l’exécution continue dans cmd/coherentyaml/coherentyaml.go:70 ;
    - seuls les deux premiers arguments sont utilisés malgré l’aide [fichier2 ...] ;
    - les codes de sortie et les erreurs ne sont pas encore formalisés.

### Modernisation encore incomplète

- Aucun branchement modernisation : les deux nouveaux commits sont directement sur master, actuellement en avance de deux commits sur origin/master.
- Les anciennes branches distantes fusion et forgotten_work existent encore.
- go.mod:3 déclare Go 1.19, alors que Docker utilise 1.27 : il faut décider et documenter la version minimale supportée.
- goccy/go-yaml reste en version 1.8.9.
- Pas de GitHub Actions, de contrôle de vulnérabilités ni de publication.
- Le README reste très sommaire et ne documente pas précisément le langage de schéma.

 Je n’ai pas pu relancer les contrôles : ni Go, ni Docker ne sont disponibles dans l’environnement courant. Comme seuls Dockerfile, TEST, go.mod et go.sum ont changé depuis l’analyse, aucun correctif fonctionnel ne peut cependant avoir résolu les anomalies du moteur ou de la CLI.

  L’ordre conseillé est donc : figer les cinq décisions fonctionnelles, réécrire les tests concernés, corriger le moteur, corriger la CLI, puis terminer versionnement, CI et documentation.
