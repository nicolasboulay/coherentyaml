FROM golang:1.27.0-bookworm

ENV CODEX_HOME=/codex-home \
    DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
        unzip \
        ripgrep \
        jq \
        less \
        file \
        poppler-utils \
        python3 \
        python3-pip \
        python3-venv \
        nodejs \
        npm \
        make \
        build-essential \
        gcc \
        procps \
        binutils-arm-none-eabi \
        gcc-arm-none-eabi \
        libnewlib-arm-none-eabi \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://chatgpt.com/codex/install.sh \
    | CODEX_NON_INTERACTIVE=1 CODEX_INSTALL_DIR=/usr/local/bin sh

RUN python3 -m venv /opt/platformio \
    && /opt/platformio/bin/pip install --no-cache-dir platformio \
    && ln -s /opt/platformio/bin/platformio /usr/local/bin/platformio \
    && ln -s /opt/platformio/bin/pio /usr/local/bin/pio

RUN mkdir -p "$CODEX_HOME/skills" \
    && git clone --depth 1 https://github.com/obra/superpowers.git "$CODEX_HOME/superpowers" \
    && ln -s "$CODEX_HOME/superpowers/skills" "$CODEX_HOME/skills/superpowers" \
    && git config --global --add safe.directory /workspace

COPY container-install.sh /usr/local/bin/container-install.sh

WORKDIR /workspace

ENTRYPOINT ["bash", "/usr/local/bin/container-install.sh"]
CMD ["codex", "--dangerously-bypass-approvals-and-sandbox"]
