class Gochangelog < Formula
  desc "Automatically generate changelogs from your git commit history"
  homepage "https://github.com/richard-egeli/gochangelog"
  url "https://github.com/richard-egeli/gochangelog/archive/v1.0.0.tar.gz"
  sha256 "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"gochangelog"
  end

  test do
    system "#{bin}/gochangelog", "--version"
  end
end
