{
  description = "lCyou's nix overlay";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }: let
    systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
    forAllSystems = nixpkgs.lib.genAttrs systems;
  in {
    overlays.default = final: prev: {
      post2bsky = final.callPackage ./pkgs/post2bsky/post2bsky.nix {};
    };

    packages = forAllSystems (system: {
      post2bsky = nixpkgs.legacyPackages.${system}.callPackage ./pkgs/post2bsky/post2bsky.nix {};
      default = self.packages.${system}.post2bsky;
    });
  };
}
