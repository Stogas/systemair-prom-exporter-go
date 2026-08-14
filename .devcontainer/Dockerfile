# any ARG values are default, override via build args in devcontainer.json instead
FROM alpine:3 AS downloader

ARG TARGETARCH TARGETOS

RUN apk add --no-cache curl xz tar

# installs to /usr/local/other-utils
RUN mkdir /usr/local/other-utils
# golangci-lint
ARG GOLANGCI_LINT_VERSION=2.12.2
RUN GOLANGCI_LINT_DOWNLOAD="golangci-lint-${GOLANGCI_LINT_VERSION}-${TARGETOS}-${TARGETARCH}" \
	&& curl -LO "https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/${GOLANGCI_LINT_DOWNLOAD}.tar.gz" \
	&& tar -xf "${GOLANGCI_LINT_DOWNLOAD}.tar.gz" -C /usr/local/other-utils --strip-components=1 "${GOLANGCI_LINT_DOWNLOAD}/golangci-lint"
# lefthook
ARG LEFTHOOK_VERSION=2.1.10
RUN case "$TARGETARCH" in \
	amd64) LEFTHOOK_ARCH="x86_64" ;; \
	arm64) LEFTHOOK_ARCH="aarch64" ;; \
	*) echo "Unsupported arch: $TARGETARCH"; exit 1 ;; \
	esac \
	&& LEFTHOOK_DOWNLOAD="lefthook_${LEFTHOOK_VERSION}_${TARGETOS}_${LEFTHOOK_ARCH}.gz" \
	&& curl -LO "https://github.com/evilmartians/lefthook/releases/download/v${LEFTHOOK_VERSION}/${LEFTHOOK_DOWNLOAD}" \
	&& gunzip -c "${LEFTHOOK_DOWNLOAD}" > /usr/local/other-utils/lefthook \
	&& chmod +x /usr/local/other-utils/lefthook

# installs to /usr/local/go
# changes most often, so is installed last
ARG GO_VERSION=1.26.6
RUN GO_DOWNLOAD="go${GO_VERSION}.${TARGETOS}-${TARGETARCH}" \
	&& curl -LO "https://go.dev/dl/${GO_DOWNLOAD}.tar.gz" \
	&& tar -xf "${GO_DOWNLOAD}.tar.gz" -C /usr/local

FROM debian:trixie-slim AS runner

RUN apt-get update && apt-get install -y sudo bash git openssh-client curl nano jq gnupg2 locales \
	&& rm -rf /var/lib/apt/lists/* \
	&& echo "en_US.UTF-8 UTF-8" > /etc/locale.gen \
	&& locale-gen
ENV LANG=en_US.UTF-8 LANGUAGE=en_US:en LC_ALL=en_US.UTF-8

RUN groupadd --gid 1000 developer \
	&& useradd --uid 1000 --gid developer --create-home --shell /bin/bash developer \
	&& echo "developer ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Go & others
COPY --from=downloader /usr/local/other-utils /usr/local/other-utils
COPY --from=downloader /usr/local/go /usr/local/go

USER developer
ENV PATH=$PATH:/home/developer/go/bin:/usr/local/go/bin:/usr/local/other-utils
