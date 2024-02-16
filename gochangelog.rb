class Gochangelog < Formula
  desc "Automatically generate changelogs from your git commit history"
  homepage "https://github.com/richard-egeli/gochangelog"
  url "https://github.com/richard-egeli/gochangelog/archive/v1.0.1.tar.gz"
  sha256 "f5bdc2f7ba1b587aca53094bb7eb0871bcbf7d6d90ab90b3a190967a9ccd217b"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"gochangelog"
  end

  test do
    system "#{bin}/gochangelog", "--version"
  end
end
