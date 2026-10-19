# Use a specific version of distroless based on SHA256 digest for reproducibility and security
# Using distroless for the final image to minimize attack surface and size
FROM gcr.io/distroless/base-debian12:nonroot@sha256:7f0c72cd138b442ae0deeb69c08b1acf5525439ba251a49ad93c320a061567e5 AS runner

WORKDIR /app
COPY ./dist/apps/go/meowbot .

ENTRYPOINT ["/app/main"]