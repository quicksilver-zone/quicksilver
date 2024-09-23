FROM imbios/bun-node:22-slim AS deps
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get -y update && \
  apt-get install -yq openssl git ca-certificates tzdata && \
  ln -fs /usr/share/zoneinfo/Europe/London /etc/localtime && \
  dpkg-reconfigure -f noninteractive tzdata
WORKDIR /app

# Install dependencies based on the preferred package manager
COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile

# Build the app
FROM deps AS builder
WORKDIR /app
COPY . .
RUN mv next.config.prod.js next.config.js

RUN bun --bun run next build

FROM node:22-slim AS runner
WORKDIR /app

ENV NODE_ENV production
# Uncomment the following line in case you want to disable telemetry during runtime.
ENV NEXT_TELEMETRY_DISABLED 1

COPY --from=builder  /app/.next/standalone ./
COPY --from=builder  /app/.next/static .next/static
COPY --from=builder  /app/public ./public

CMD ["node", "server.js"]
