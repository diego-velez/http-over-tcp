{
  description = "HTTP from TCP";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        name = "command-center";

        packages = with pkgs; [
          go
          golangci-lint
        ];

        shellHook = ''
          echo "HTTP from TCP"
        '';
      };
    };
}
