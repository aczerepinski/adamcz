*** title ***
Pretty URLs with Phoenix

*** date ***
1/13/16

*** tags ***
elixir, phoenix

*** description ***
Configure a Phoenix app to use SEO-friendly slugs

*** body ***
Edit 2016-01-17: This post originally had urls with underscores, but @_StevenNunez pointed out that Google favors hyphens. The slugify code has been updated accordingly.

---

Before we get started, take a minute to gaze at your address bar and appreciate the beauty of this page's url. adamcz.com/blog/2? Not a chance. The SEO-friendly beauty you see up there is a keyword packed slug that conveys the page's intent before you even click. They're pretty easy to set up, and well worth your while, both for the aesthetics and the SEO boost. Disclaimer: I'm writing this as I learn [Phoenix](http://www.phoenixframework.org/), so if anything here looks non-idiomatic or just plain wrong, please let me know.

First of all, create a migration to add a slug field to your database (you can call it pretty_url/permalink/linksypants if you prefer). 

```elixir
def change do
  alter table(:posts) do
    add :slug, :string
  end
end
```

Next, make a few small changes in your model. In addition to adding slug to your schema macro, you'll need a short function that does string manipulation (you can skip this if you prefer to directly type your slug-url on whatever form is submitting your source string).

```elixir
defp slugify(string) do
  string
  |> String.downcase
  |> String.replace(" ", "-")
  |> String.replace(~r/[!.?']/, "")
end
```

After that you'll want to expand your changeset, so that after it casts your source field (title in my case), it pipes the changeset to a slug creation function. Keep in mind that the changeset won't always contain changes (such as when it's being used to set up an empty form), so you'll need a conditional somewhere before you attempt string manipulation. You can read about Ecto Changeset functions [here](https://hexdocs.pm/ecto/Ecto.Changeset.html), and see the possible return values for fetch_field. I match against the expected case on form submission ( {:changes, term} ), and otherwise just return the changeset untouched.

```elixir
def changeset(model, params \\ :empty) do
  model
  |> cast(params, @required_fields, @optional_fields)
  |> create_a_slug
end

defp create_a_slug(changeset) do
  case fetch_field(changeset, :title) do
    {:changes, title} ->
      slug = title |> slugify
      put_change(changeset, :slug, slug)
    _ ->
      changeset
  end
end
```

With all that in place, head over to your router, and specify the new field that will serve as your param:

```elixir
resources "/blog", BlogController, only: [:index, ...], param: "slug"
```

And then use that param to pull from your repository. Now that you'll be accessing posts by means other than their primary key, it's probably a good idea to add an index to slug. There's an example in the Ecto Migration documentation [here](https://hexdocs.pm/ecto/Ecto.Migration.html). You'll take a minor hit on insert speeds, but it will pay off every time someone hits the show action.

```elixir
def show(conn, %("slug" => slug}) do
  post = Repo.get_by(Post, slug: slug)
  if !post.active do
    conn
    |> authenticate({})
...
```

Finally, don't forget to update any links that might previously have referenced primary key:

```elixir
<%= link "the link", to: blog_path(@conn, :show, post.slug) %>
```

That's all for now! Like I said, please let me know if there is a more efficient way to do any of this (aczerepinski at google's email service). It's been a lot of fun writing a blog about building a blog, and using the blog to write the blog about the blog. 

Thanks for reading!
-Adam