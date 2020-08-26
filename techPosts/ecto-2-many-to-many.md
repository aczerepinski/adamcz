*** title ***
Ecto 2 Many to Many Associations

*** date ***
5/9/16

*** tags ***
elixir, phoenix

*** description ***
An early look at Ecto 2's many-to-many association

*** body ***
One of the features [introduced](https://github.com/elixir-ecto/ecto/blob/master/CHANGELOG.md) in Ecto 2.0 is the many-to-many relationship; a model-level abstraction over a database's join table. Documentation for linking two models together is scattered about, such as  [here](https://hexdocs.pm/ecto/Ecto.Association.ManyToMany.html), [here](https://hexdocs.pm/ecto/Ecto.Changeset.html#put_assoc/4), and [here](https://hexdocs.pm/ecto/Ecto.Repo.html#c:preload/3), but it took me a while to piece together the complete picture of how to use this in an app.

In this blog post, I want to bring that documentation into one place and overview all of the steps needed to implement many_to_many in a typical Phoenix app. For the code examples, I'll stick to this blog's own source code and discuss how posts are linked to tags.

### Step 1: Migrations

The first thing you'll need to do is add the appropriate tables to your database. That means a table for each of the models, and a join table, which you need to create manually:

```elixir
def change do
  create table(:posts_tags) do
    add :post_id, :integer
    add :tag_id, :integer
  end
end
```

### Step 2: Models

Next, you'll need to update the schemas in both models:

```elixir
schema "posts" do
  # existing post schema...
  many_to_many :tags, AdamczDotCom.Tag, join_through: "posts_tags",
    on_replace: :delete
end

schema "tags" do
  # existing tag schema...
  many_to_many :posts, AdamczDotCom.Post, join_through: "post_tags"
end
```

join_through lets Ecto know where to save the associated ids, and on_replace lets Ecto know what to do when associations that previously existed are not included on a changeset. Note that :delete needs to be used carefully, as it deletes any row not present in an updated changeset. Attempting to remove an association without on_replace raises an error, but it can be omitted if you know you won't do that (e.g., I don't update posts at all via tags so I didn't include it on the tags schema).

Finally, [cast](https://hexdocs.pm/ecto/Ecto.Changeset.html#cast_assoc/3) the association in your changeset:

```elixir
def changeset(post, params \\ %{}) do
  post
  |> cast(params, [:title, etc...])
  |> cast_assoc(:tags)
```

### Step 3: Controller

There's a fair amount of work to be done in the controller. The simple case is loading existing data for the index or show action, so I'll present that first, useless as it may be when you don't yet have any associated data to show off.

```elixir
def index(conn, _params) do
  posts = Repo.all(Post)
    |> Repo.preload(:tags)
  render conn, "index.html", posts: posts
end
```

In order to render a blank form in the new action, you need to include an empty tags list. This is needed if your new and edit templates share the same form partial, as mine do. Also, you'll want to pull the full list of existing tags in order to present them as options on the form:

```elixir
def new(conn, _params) do
  changeset = Post.changeset(%Post{tags: []})
  tags = Repo.all(Tag)
  render conn, "new.html", changeset: changeset, tags: tags
end
```

The update action is where things got difficult for me. I haven't figured an ideal way to pass nested tags params from the form, so I'm doing a bit of extra work to process the tags params that arrive separately. If anyone knows how to omit this step, please reach out to me!

After wrangling the form data into a tags changeset, we start out with the default posts changesest, and smash the two together with Changeset.put_assoc:

```elixir
def update(conn, %{"slug" => slug, "post" => post_params, "tags" => tags}) do
  tags_to_associate = tag_changeset(tags)
  post = Repo.get_by(Post, slug: slug)
    |> Repo.preload(:tags)
  changeset = Post.changeset(post, post_params)
    |> Ecto.Changeset.put_assoc(:tags, tags_to_associate)

  case Repo.update(changeset) do
  # boilerplate after this...
end

defp tag_changeset(tags) do
  tags_to_list(tags) # converts the form params to a list of integers
  |> get_tag_structs # grabs the appropriate tags from the database
  |> Enum.map(&Ecto.Changeset.change/1) # wraps those structs in a changeset
end

defp tags_to_list(tags) do
  Enum.filter_map(tags, fn(tag) -> String.to_atom(elem(tag, 1)) == true end,
    fn(tag) -> String.to_integer(elem(tag, 0)) end)
end

defp get_tag_structs(tag_list) do
  Tag.by_id_list(tag_list) # model code here is: where([t], t.id in ^ids)
  |> Repo.all
end
```

### Step 4: Views

If you're using a form partial like I am, it won't magically have access to tags. You'll need to explicitly pass them like this:

```elixir
<%= render "form.html", changeset: @changeset,
                        tags: @tags,
                        action: blog_path(@conn, :create) %>
```

Then in the form, render *all* tags, and for each one, check to see if the current changeset says it should default to checked.

```elixir
<%= for tag <- @tags do %>
  <% checked = AdamczDotCom.BlogView.is_checked(tag, @changeset) %>
  <%= checkbox :tags, "#{tag.id}", value: tag.name, checked: checked %> <%= tag.name %>
<% end %>
```

And that’s about it. I hope this is helpful, and please reach out to me with any corrections or suggestions. Thanks for reading,

-Adam