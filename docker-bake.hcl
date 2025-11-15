group "default" {
  targets = ["golang"]
}

target "golang" {
  tags = ["ukoloff/na01:golang", "ghcr.io/ukoloff/na01:golang"]
  pull = true
  dockerfile-inline = <<-EOT
    FROM golang:alpine AS build

    WORKDIR /repo
    COPY ./go/. ./go/.
    RUN go build -C go/main -o na01 -ldflags "-s -w"

    FROM gcr.io/distroless/static-debian12

    COPY --from=build /repo/go/main/na01 /bin/na01

    ENTRYPOINT ["/bin/na01", "www"]
    EOT

  labels = {
    "org.opencontainers.image.authors" = "Stanislav.Ukolov@omzglobal.com"
    "org.opencontainers.image.description" = "NetAngels DNS-01 helper for Lego Let's Encrypt / ACME client"
    "org.opencontainers.image.source" = "https://github.com/ukoloff/netangels.dns-01"
  }
}
