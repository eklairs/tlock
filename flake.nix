{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ { self, nixpkgs, flake-parts, ... }:
  flake-parts.lib.mkFlake { inherit inputs; } {
    systems = [ "x86_64-linux" ];

    perSystem = { system, ... }: let
      pkgs = import nixpkgs { inherit system; };
      git_version = self.shortRev or self.dirtyShortRev;
    in rec {
      packages = rec {
        tlock = pkgs.callPackage ./nix { inherit git_version; };
        default = tlock;
      };

      devShells.default = pkgs.mkShell {
        inputsFrom = [ packages.default ];
      };
    };
  };
}
