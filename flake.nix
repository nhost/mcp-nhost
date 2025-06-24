{
  inputs = {
    nixops.url = "github:nhost/nixops";
    nixpkgs.follows = "nixops/nixpkgs";
    flake-utils.follows = "nixops/flake-utils";
    nix-filter.follows = "nixops/nix-filter";
    nix2container.follows = "nixops/nix2container";
  };

  outputs = { self, nixops, nixpkgs, flake-utils, nix-filter, nix2container }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [
          nixops.overlays.default
          (final: prev: {
            certbot-full = prev.certbot.overrideAttrs (old: {
              doCheck = false;
            });
          })
        ];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        src = nix-filter.lib.filter {
          root = ./.;
          include = with nix-filter.lib;[
            ".golangci.yaml"
            "go.mod"
            "go.sum"
            "nhost/auth/openapi.yaml"
            "nhost/graphql/openapi.yaml"
            "tools/cloud/schema.graphql"
            "tools/cloud/schema-with-mutations.graphql"
            (inDirectory "testdata")
            (inDirectory "graphql/testdata")
            (inDirectory "vendor")
            isDirectory
            (nix-filter.lib.matchExt "go")
          ];
        };

        nix-src = nix-filter.lib.filter {
          root = ./.;
          include = with nix-filter.lib;[
            (matchExt "nix")
          ];
        };

        checkDeps = with pkgs; [
          rover
          oapi-codegen
        ];

        buildInputs = [
        ];

        nativeBuildInputs = [
        ];

        nix2containerPkgs = nix2container.packages.${system};
        nixops-lib = nixops.lib { inherit pkgs nix2containerPkgs; };

        name = "mcp-nhost";
        version = "0.0.0-dev";
        created = "1970-01-01T00:00:00Z";
        submodule = ".";

        tags = [ ];

        ldflags = [
          "-X main.Version=${version}"
        ];
      in
      {
        checks = flake-utils.lib.flattenTree {
          nixpkgs-fmt = nixops-lib.nix.check { src = nix-src; };

          go-checks = nixops-lib.go.check {
            inherit src submodule ldflags tags buildInputs nativeBuildInputs checkDeps;
          };
        };

        devShells = flake-utils.lib.flattenTree {
          default = nixops-lib.go.devShell {
            buildInputs = [
              pkgs.skopeo
            ] ++ checkDeps ++ buildInputs ++ nativeBuildInputs;
          };
        };

        packages = flake-utils.lib.flattenTree rec {
          mcp-nhost = nixops-lib.go.package {
            inherit name src version ldflags buildInputs nativeBuildInputs;
          };

          mcp-nhost-arm64-darwin = (nixops-lib.go.package {
            inherit name src version ldflags buildInputs nativeBuildInputs;
          }).overrideAttrs (old: old // {
            env = {
              GOOS = "darwin";
              GOARCH = "arm64";
              CGO_ENABLED = "0";
            };
          });

          mcp-nhost-amd64-darwin = (nixops-lib.go.package {
            inherit name src version ldflags buildInputs nativeBuildInputs;
          }).overrideAttrs (old: old // {
            env = {
              GOOS = "darwin";
              GOARCH = "amd64";
              CGO_ENABLED = "0";
            };
          });

          mcp-nhost-arm64-linux = (nixops-lib.go.package {
            inherit name src version ldflags buildInputs nativeBuildInputs;
          }).overrideAttrs (old: old // {
            env = {
              GOOS = "linux";
              GOARCH = "arm64";
              CGO_ENABLED = "0";
            };
          });

          mcp-nhost-amd64-linux = (nixops-lib.go.package {
            inherit name src version ldflags buildInputs nativeBuildInputs;
          }).overrideAttrs (old: old // {
            env = {
              GOOS = "linux";
              GOARCH = "amd64";
              CGO_ENABLED = "0";
            };
          });

          docker-image-arm64 = nixops-lib.go.docker-image {
            inherit name version created buildInputs;
            arch = "arm64";

            package = mcp-nhost-arm64-linux;
          };

          docker-image-amd64 = nixops-lib.go.docker-image {
            inherit name version created buildInputs;
            arch = "amd64";

            package = mcp-nhost-amd64-linux;
          };

          default = mcp-nhost;

        };

      }
    );
}
