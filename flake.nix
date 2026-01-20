{
  description = "A development shell for Golang";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };
  outputs =
    inputs@{
      flake-parts,
      self,
      ...
    }:
    # https://flake.parts/
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [ "x86_64-linux" "aarch64-darwin" ];
      perSystem = { pkgs, inputs', lib, system, ... }:
        let
          currentPath = builtins.getEnv "PWD";
        in
        {
          devShells.default = pkgs.mkShell {
            name = "dev";

            # Available packages on https://search.nixos.org/packages
            buildInputs = with pkgs; [
                gnumake
                go
                golangci-lint
            ];

            shellHook = ''
              echo "Entering a new go devshell"
              echo "Go version: $(${lib.getExe pkgs.go} version)"
              echo "golangci-lint version: $(${lib.getExe pkgs.golangci-lint} --version)"
            '';
          };
        };
    };
}
