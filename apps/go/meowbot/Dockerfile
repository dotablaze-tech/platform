# Use a specific version of distroless based on SHA256 digest for reproducibility and security
# Using distroless for the final image to minimize attack surface and size
FROM gcr.io/distroless/base-debian12:nonroot@sha256:7f0c72cd138b442ae0deeb69c08b1acf5525439ba251a49ad93c320a061567e5 AS runner

WORKDIR /app
COPY ./dist/apps/go/meowbot .

# The :nonroot base already defaults to this user; declaring it keeps the
# image non-root if the base tag ever changes.
USER nonroot:nonroot

# Probes the API /liveness endpoint; start-period covers the DB and Discord
# connect at boot.
HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
  CMD ["/app/main", "healthcheck"]

ENTRYPOINT ["/app/main"]