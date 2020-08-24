package markdown

import "testing"

func TestConvertLinks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hi", "hi"},
		{"<p>hi</p>", "<p>hi</p>"},
		{
			"[click on this](https://www.adamcz.com)",
			`<a href="https://www.adamcz.com">click on this</a>`,
		},
		{
			"[ECTO](https://hexdocs.pm/ecto/Ecto.html)",
			`<a href="https://hexdocs.pm/ecto/Ecto.html">ECTO</a>`,
		},
	}

	for i, test := range tests {
		if actual := convertLinks(test.input); actual != test.expected {
			t.Errorf("%d: expected %s, got %s", i, test.expected, actual)
		}
	}
}

func TestConvertInlineCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hi", "hi"},
		{"`mix phx.new`", `<code class="inline">mix phx.new</code>`},
	}

	for i, test := range tests {
		if actual := convertInlineCode(test.input); actual != test.expected {
			t.Errorf("%d: expected %s, got %s", i, test.expected, actual)
		}
	}
}

func TestConvertElement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "<p>hello</p>"},
		{"# hi there", "<h1>hi there</h1>"},
		{"### hi", "<h3>hi</h3>"},
		// {"```elixir\ntacos = delicious\n```", `<code class="elixir">tacos = delicious</code>`},
		// {
		// 	"```elixir\n # lib/tacos/repo.ex defmodule Tacos.TacoRepo do use Ecto.Repo, otp_app: :tacos, adapter: Ecto.Adapters.Postgres end" +
		// 		"defmodule Tacos.UserRepo do use Ecto.Repo, otp_app: :tacos, adapter: Ecto.Adapters.Postgres end \n```",
		// 	`<code class="elixir">asdf</code>`,
		// },
	}

	for i, test := range tests {
		if actual := convertElement(test.input); actual != test.expected {
			t.Errorf("%d: expected %s, got %s", i, test.expected, actual)
		}
	}
}
