ARG RUBY_VERSION=3
FROM ruby:${RUBY_VERSION} AS builder

WORKDIR /
COPY ./Gemfile ./Gemfile

ENV PLATFORMS="aarch64-linux-gnu aarch64-linux-musl arm64-darwin ruby x86-linux-gnu x86-linux-musl x86_64-darwin x86_64-linux-gnu x86_64-linux-musl"
RUN for platform in $PLATFORMS; do \
  bundle lock --add-platform $platform ; \
  done

RUN bundle install

FROM scratch AS output
COPY --from=builder . .
CMD ["ls", "-l", "."]
