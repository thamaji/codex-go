# syntax=docker/dockerfile:1
ARG UBUNTU_VERSION=22.04
ARG _DEV_CONTAINERS_BASE_IMAGE=dev_container_auto_added_stage_label
FROM ubuntu:${UBUNTU_VERSION}

# set apt mirror server
ARG APT_MIRROR=""
RUN [ "${APT_MIRROR:-}" != "" ] && sed -i -r "s@http://(\\w+.)?archive\.ubuntu\.com/ubuntu/@${APT_MIRROR}@" || :

# enable cache that apt keeps
RUN rm -f /etc/apt/apt.conf.d/docker-clean \
    && echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' >/etc/apt/apt.conf.d/keep-cache

# update ca-certificates
RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates

# japanize
RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo 'Asia/Tokyo' >/etc/timezone \
    && apt-get install -y --no-install-recommends \
        language-pack-ja-base \
        language-pack-ja
ENV LANG=ja_JP.UTF-8

# create user
ARG USERNAME=vscode
RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        sudo \
    && groupadd --gid 5000 ${USERNAME} \
    && useradd --uid 5000 --gid 5000 --home-dir /home/${USERNAME} --create-home --shell /bin/bash ${USERNAME} \
    && echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" >>/etc/sudoers.d/${USERNAME} \
    && chmod 0440 /etc/sudoers.d/${USERNAME}

# persist bash history
RUN mkdir -p /home/${USERNAME}/.cache/bash \
    && touch /home/${USERNAME}/.cache/bash/.bash_history \
    && chown -R ${USERNAME}:${USERNAME} /home/${USERNAME}/.cache/bash \
    && ln -s /home/${USERNAME}/.cache/bash/.bash_history /home/${USERNAME}/.bash_history \
    && echo "export PROMPT_COMMAND='history -a'" >>/home/${USERNAME}/.bashrc

# persist vscode extensions
RUN mkdir -p /home/${USERNAME}/.vscode-server/extensions /home/${USERNAME}/.vscode-server-insiders \
    && chown -R ${USERNAME}:${USERNAME} /home/${USERNAME}/.vscode-server /home/${USERNAME}/.vscode-server-insiders \
    && ln -s /home/${USERNAME}/.vscode-server/extensions /home/${USERNAME}/.vscode-server-insiders/extensions

# install common dev tools
RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        bash-completion \
        curl \
        git \
        gnupg2 \
        iputils-ping \
        less \
        net-tools \
        tar \
        time \
        unzip \
        xz-utils \
        zip \
        ssh

# install golang
ARG GO_VERSION=1.25.1
RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        curl \
        git \
    && curl -fsSL --compressed "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" | tar -xz -C /usr/local
ENV GOROOT=/usr/local/go \
    GOPATH=/home/${USERNAME}/go \
    PATH=/home/${USERNAME}/go/bin:/usr/local/go/bin:${PATH}

# install Codex CLI
ARG CODEX_VERSION=rust-v0.36.0
RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        coreutils \
        curl \
        findutils \
        git \
        grep \
        ripgrep \
        sed \
        tar \
    && mkdir -p "/usr/local/codex-${CODEX_VERSION}" \
    && curl -fsSL "https://github.com/openai/codex/releases/download/${CODEX_VERSION}/codex-x86_64-unknown-linux-musl.tar.gz" \
        | tar -xz -C "/usr/local/codex-${CODEX_VERSION}" \
    && ln -s "/usr/local/codex-${CODEX_VERSION}/codex-x86_64-unknown-linux-musl" /usr/local/bin/codex

USER ${USERNAME}
