FROM  ghcr.io/oracle/oraclelinux8-instantclient:21

RUN dnf config-manager --add-repo https://mise.jdx.dev/rpm/mise.repo

RUN dnf install -y git zip unzip dnf-plugins-core go-toolset mise

RUN dnf groupinstall "Development Tools"

RUN useradd -m devcontainer
