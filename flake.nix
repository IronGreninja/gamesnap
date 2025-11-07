{
  description = "gamesnap: A dead simple game save snapshot cli tool";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = {self, ...} @ inputs: let
    goVersion = 25; # Change this to update the whole stack

    supportedSystems = [
      "x86_64-linux"
      "aarch64-linux"
      "x86_64-darwin"
      "aarch64-darwin"
    ];
    forEachSupportedSystem = f:
      inputs.nixpkgs.lib.genAttrs supportedSystems (
        system:
          f {
            pkgs = import inputs.nixpkgs {
              inherit system;
              overlays = [inputs.self.overlays.default];
            };
          }
      );
  in {
    overlays.default = final: prev: {
      go = final."go_1_${toString goVersion}";
    };

    devShells = forEachSupportedSystem (
      {pkgs}: {
        default = pkgs.mkShellNoCC {
          packages = with pkgs; [
            # go (version is specified by overlay)
            go

            delve

            gotools
            go-tools
            gomodifytags
            # golangci-lint
            gopls
            gotests

            goreleaser
          ];

          env = {
            CGO_ENABLED = "0"; # fixes delve
            GOTOOLCHAIN = "local";
          };
        };
      }
    );
  };
}
