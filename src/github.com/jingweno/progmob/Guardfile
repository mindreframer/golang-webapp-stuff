ENV["PORT"] = "8080"
ENV["DATABASE_URL"] = "mongodb://localhost/progmob"
ENV["GIT_HUB_CLIENT_ID"] = "2ddf3af12b8d30259074"
ENV["GIT_HUB_CLIENT_SECRET"] = "df831a8e4a81fa7c77e0018100ccc2ebde50e3b5"

class GoWebApp
  attr_reader :name

  def initialize(name)
    @name = name
  end

  def build
    puts `go build -v`
    self
  end

  def test
    puts `go test -v ./...`
    self
  end

  def start
    Process.spawn("./#{name}")
    self
  end

  def stop
    `pgrep #{name} | xargs kill -9`
    self
  end

  def restart
    stop
    start
  end
end

$app = GoWebApp.new File.basename(File.dirname(__FILE__))
$app.build.test.restart

guard "shell" do
  watch(/(.*).go/) do |m|
    $app.test
  end

  watch(/\A(?:(?!_test).)*.(go|html|css)\z/) do |m|
    puts "Changed #{m[0]}. Rebuilding and restarting..."
    $app.build.restart
  end
end

at_exit {
  $app.stop
}
