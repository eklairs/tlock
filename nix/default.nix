{ pkgs, lib, buildGoModule, git_version, ... }:
buildGoModule rec {
  pname = "tlock";
  version = "1.0.0";

  src = lib.cleanSource ../.;
  vendorHash = "sha256-G402CigSvloF/SI9Wbcts/So1impMUH5kroxDD/KKew=";

  ldflags = [
    "-X github.com/eklairs/tlock/tlock-internal/constants.VERSION=v${version}-${git_version}"
  ];

  excludedPackages = [ "bubbletea" ];
  
  meta = with lib; {
    description = "Two-Factor Authentication Tokens Manager in Terminal";
    homepage = "https://github.com/eklairs/tlock";
    license = licenses.mit;
  };
}
