{
  inputs = {
    nixops.url = "github:nhost/nixops";
    nixpkgs.follows = "nixops/nixpkgs";
    flake-utils.follows = "nixops/flake-utils";
    nix-filter.follows = "nixops/nix-filter";
  };

  outputs = { self, nixops, nixpkgs, flake-utils, nix-filter }:
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

        nixops-lib = nixops.lib { inherit pkgs; };

        name = "mcp-nhost";
        version = "0.0.0-dev";
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

            preCheck = ''
              echo "➜ Getting access token"
              export NHOST_ACCESS_TOKEN=$(bash ${src}/get_access_token.sh)
            '';
          };
        };

        devShells = flake-utils.lib.flattenTree {
          default = nixops-lib.go.devShell {
            buildInputs = [
            ] ++ checkDeps ++ buildInputs ++ nativeBuildInputs;
          };
        };

        packages = flake-utils.lib.flattenTree rec {
          mcp-nhost = nixops-lib.go.package {
            inherit name src version ldflags buildInputs nativeBuildInputs;
          };

          docker-image = nixops-lib.go.docker-image {
            inherit name version buildInputs;
            package = mcp-nhost;
          };

          default = mcp-nhost;

        };

      }
    );
}
