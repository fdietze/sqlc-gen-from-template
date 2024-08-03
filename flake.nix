{
  inputs = {
    # nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    nixpkgs.url =
      "github:nixos/nixpkgs?ref=68c9ed8bbed9dfce253cc91560bf9043297ef2fe";

    # for `flake-utils.lib.eachSystem`
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:

    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ ];
          config.allowUnfree = false;
        };
      in {
        devShells = {
          default = with pkgs; pkgs.mkShellNoCC { buildInputs = [ go ]; };
        };
      });
}
