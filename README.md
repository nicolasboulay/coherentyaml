# coherentyaml

The goal of coherentYaml is to create scheme for yaml. Instead of
using a description of what is expected, the schema looks like an
example with some logical operator.

A minium set of keyword have been created :
- "Coherent" (a kind of logicial 'and')
- "OR"
- "Not"

The "Not" is not fonctional.

./coherentyaml fichier1.yml fichier2.yml

The 2 data set are compared. 

"Type" are defined by 'neutral' élement (1, -1, 1.0, ""). Each key
in object must be coherent. 

Each data is coherent to it-self. Coherence is symetrical.

The tool compile with go build into cmd/coherentyaml.

## Development

Enable the versioned Git hooks after cloning the repository:

```sh
git config core.hooksPath .githooks
```

The pre-commit hook formats staged Go files with `gofmt` and adds the
formatted result back to the commit.
