FROM --platform=$BUILDPLATFORM node:20.9.0 AS deps
WORKDIR /app
COPY frontend/webapp/package.json frontend/webapp/yarn.lock ./
RUN yarn install 

FROM --platform=$BUILDPLATFORM node:20.9.0 AS builder
WORKDIR /webapp
COPY --from=deps /app/node_modules ./node_modules
COPY frontend/webapp .
RUN yarn build

FROM --platform=$BUILDPLATFORM golang:1.24 AS backend
WORKDIR /app
COPY . .
COPY --from=builder /webapp/out frontend/webapp/out
WORKDIR /app/frontend
ARG TARGETARCH
ARG LD_FLAGS
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -ldflags="${LD_FLAGS}" -o odigos-ui

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=backend /app/frontend/odigos-ui .
USER 65532:65532
ENTRYPOINT ["/app/odigos-ui", "--address", "0.0.0.0"]
