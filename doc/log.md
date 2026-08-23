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
