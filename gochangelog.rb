class Gochangelog < Formula
  desc "Automatically generate changelogs from your git commit history"
  homepage "https://github.com/richard-egeli/gochangelog"
  url "https://github.com/richard-egeli/gochangelog/archive/v1.0.0.tar.gz"
  sha256 "531608f54cf7983c321f12f1b82e659a33c802f2392da4388469810a724cd117"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"gochangelog"
  end

  test do
    system "#{bin}/gochangelog", "--version"
  end
end
