#!/usr/bin/env bash
set -euo pipefail

export CODEX_HOME="${CODEX_HOME:-/codex-home}"
mkdir -p "$CODEX_HOME/skills"

# un lien par skill utilisé par le projet
#ln -sfn /skill-repos/common-skills/skills/project-container-bootstrap "$CODEX_HOME/skills/project-container-bootstrap"
#ln -sfn /skill-repos/common-skills/skills/easypcb-eprj-review "$CODEX_HOME/skills/easypcb-eprj-review"

#timeout 10m codex exec --dangerously-bypass-approvals-and-sandbox --skip-git-repo-check \
#    'Installe le skill officiel OpenAI "pdf" depuis la liste curated du repo openai/skills dans $CODEX_HOME/skills, en utilisant python3 si necessaire. Ne modifie pas /workspace.'

exec "$@"
