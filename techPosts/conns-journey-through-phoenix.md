*** title ***
Conn's journey through Phoenix

*** date ***
1/19/16

*** tags ***
elixir, phoenix

*** description ***
Follow the conn struct's journey through a Phoenix page view

*** body ***
If you want to understand the Phoenix web framework, you really want to understand [Plug](https://github.com/elixir-lang/plug). Plug is an Elixir web request library that's hard to define succinctly because it plays several roles. It's an adapter that takes an HTTP request from the Cowboy web server and returns a struct representing both the request and the eventual response. It's also a specification for thin middleware layers (called plugs) that accept and return the Plug.Conn struct. The struct is referred to in function arguments as 'conn' so I'll mostly refer to it that way here. You can easily build your own plugs to insert into an existing framework or you could stack an entire framework on top of Plug , as Chris McCord did with Phoenix. The struct is defined [right here](https://github.com/elixir-lang/plug/blob/master/lib/plug/conn.ex). In this blog post I want to follow the journey of conn, in the context of how this blog post was rendered.

The first Phoenix plug is Endpoint, which in its 'call' function takes the base struct from the Plug library, adds the path name and your application's secret key, and then pipes to all of the plugs listed in your application's endpoint.ex file. The very first of them (assuming default configuration) is Plug.Static, which is used to bypass most of the Phoenix framework when a request for a static asset like a CSS file comes in. The bypassing part is done like this:

```elixir
def halt(%Conn{} = conn) do
  %{conn | halted: true}
end
```
In other words, if conn's halted key is set to true, the rest of the pipeline won't touch it. Right after that, Plug.RequestId generates a unique ID for the request and sets conn's req_headers key accordingly. One of the next plug modules is Plug.MethodOverride, which takes an HTTP verb from a POST request's _method parameter and uses it to set conn's method key. Let's take a look:

```elixir
# Plug.MethodOverride
@allowed_methods ~w(DELETE PUT PATCH)

defp override_method(conn, body_params) do
  method = (body_params["_method"] || "") |> String.upcase

  cond do
    method in @allowed_methods -> %{conn | method: method}
    true                       -> conn
  end
end
```
If the method from body_params pattern matches against the list of allowed methods, the method key is changed. Otherwise, conn is just returned. There are a few more "administrative" plugs like that and then finally...

```elixir
plug AdamczDotCom.Router
```

Yep, the router is yet another plug. It takes conn, transforms a few keys and spits it back out.  In a Phoenix router, you'll see a pipeline macro that looks like this:

```elixir
pipeline :browser do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :fetch_flash
    plug :protect_from_forgery
    plug :put_secure_browser_headers
end
```
Those are the default plugs, and on this site there's also an auth plug in there to stop random folks from editing this blog. Hopefully it works and I really wrote all of this. Those plugs get added to conn's before_send key. The router also adds the blog post you requested to the params key and and the appropriate controller action to the private key:

```elixir
before_send: (all the plugs from pipeline :browser)
params: %{"slug" => "conns-journey-through-phoenix"}
private: %{:phoenix_action => :show, phoenix_controller => AdamczDotCom.BlogController, :phoenix_format => "html", etc.}
```
Just as with any MVC framework, the Phoenix controller (which is of course a plug) asks the data layer for resources specific to this request. Then in the final step, it passes those resources to a render function as "assigns." For this page, "assigns" contains things I added like the blog post's title & body content, plus some things Phoenix added behind the scenes like the layout that wraps my blog template.

```elixir
def render(conn, template, assigns)
  # bunch of code
  send_resp(conn, conn.status || 200, content_type, data)
end
```
I'm glossing over how 'render' works because that's surely a blog post of its own. But at a high level it takes the completed conn, the blog template, the blog content, and hands a response to the very last plug, Plug.Conn.send_resp. And now you're looking at it.

This approach of transforming a data structure through a series of functions until it contains a completed HTML string is both simple, and fast. It's also really flexible, since it is so easy to inject your own functions (provided they uphold the plug contract) at any point along the way, or halt and skip the rest of the stack. I'm having a lot of fun learning Phoenix, and gladly welcome feedback and/or corrections if I have any of the details wrong. Feel free to send me an email, aczerepinski at Google's email service.

Thanks for reading!

-Adam