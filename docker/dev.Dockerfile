FROM surnet/alpine-wkhtmltopdf:3.21.3-024b2b2-small AS wkhtmltopdf

FROM golang:1.25-alpine AS base

# Install dependencies
# RUN apk add 
RUN apk add --no-cache \
    git \
    gcc \
    musl-dev \
    bash \
    dos2unix \
    make \
    curl \
    tzdata \
    libxi \
    libxtst \
    libstdc++ \
    libx11 \
    libxrender \
    libxext \
    libssl3 \
    ca-certificates \
    fontconfig \
    freetype \
    ttf-dejavu \
    ttf-droid \
    ttf-freefont \
    ttf-liberation \
    && apk add --no-cache --virtual .build-deps \
    msttcorefonts-installer \
    && update-ms-fonts \
    && fc-cache -f \
    && rm -rf /tmp/* \
    && apk del .build-deps

# =========================
# Install golang-migrate (Linux AMD64)
# =========================
ENV MIGRATE_VERSION=4.17.0

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz \
    | tar -xz \
    && mv migrate /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate


COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/wkhtmltopdf
COPY --from=wkhtmltopdf /lib/libwkhtmltox* /lib/

# Set the timezone to GMT+7
RUN cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime && \
    echo "Asia/Bangkok" > /etc/timezone

# Install Air & Swag (pakai versi yang stabil agar tidak error)
RUN go install github.com/air-verse/air@latest &&  go install github.com/swaggo/swag/cmd/swag@v1.8.12

WORKDIR /app

# Salin go.mod & go.sum dulu agar cache build tetap optimal
COPY go.mod go.sum ./
RUN go mod download

# Copy entrypoint script and convert line endings (for Windows compatibility)
COPY entrypoint.sh /entrypoint.sh
RUN dos2unix /entrypoint.sh && \
    chmod +x /entrypoint.sh

# Salin semua kode setelah dependency
COPY . .

# # Generate Swagger docs
# RUN /go/bin/swag init -g cmd/main.go --output docs --parseDependency --parseInternal=false

# Buka port aplikasi
EXPOSE 2000

ENTRYPOINT ["/entrypoint.sh"]
