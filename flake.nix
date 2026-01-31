{
  description = "A dev shell for the beautiful hackathon";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }: let
    system = "x86_64-linux";
  in {
      devShells."${system}".default = let
        pkgs = import nixpkgs { inherit system; };
      in pkgs.mkShell {
        packages = with pkgs; [
	  nodejs_24
	  go
	];
      };
    };
}
