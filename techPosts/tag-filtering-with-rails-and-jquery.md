*** title ***
Tag Filtering with Rails and jQuery

*** date ***
7/13/15

*** tags ***
ruby

*** description ***
Add interactive tags to a Rails app

*** body ***
My current personal/fun project, [Jazz Robot](http://jazzrobot.herokuapp.com), is a web app that facilitates practicing songs in all 12 keys. It's the app I wish I had when I was a Berklee student preparing for the weekly Sunday jam at Wally's. In this blog post I'll be discussing the site's song index page, which features an interface for filtering songs by their tags (difficulty, genre, or composer). For context, I recommend visiting the page before continuing.

### Step 0: Why?

As I thought about the range of users this site might cater to (novice to professional), I anticipated different UI/UX needs and usage patterns. An advanced musician would want an enormous database of songs to select from, and would want to filter them by composer (or perhaps just command-f search by title). A novice would be overwhelmed by dozens of unfamiliar song titles, and wouldn't know where to start unless they could filter down to the easiest songs. In either use case, having a variety of filter categories does way more than just make song selection more efficient. It instantly educates the user about the contents of the site, and conveys the breadth of that content.

With a topic as ubiquitous as music, everyone already knows that songs have genres and composers. If you're designing a site about something less universally known, an interface like this implicitly teaches visitors about your business domain.

### Step 1: Create Your Tags

To get started, create a model and migration for each tag category. Unless you anticipate creating a web UI for managing tag values, I recommend the "model" generator, which does not create routes, controllers or views:

```bash
rails g model difficulty name
```

In many cases, your tag tables will have just two columns - ID and name. Generate an additional migration to add the new tag IDs to your primary model (Song in my case). If you want to enable multiple tags from the same category (as I did with genres), create a join table and store the IDs there instead.

### Step 2: Pull Your Tags into the View

You'll start by loading the tag data into instance variables in your controller file:

```ruby
def index
  @songs = Song.order(name: :asc).includes(:genres, :difficulty, :composer)
  @genres = Genre.order(name: :asc)
  @difficulties = Difficulty.all
  @composers = Composer.order(last_name: :asc)
end
```

The ".includes" method is extremely important here, because now that our songs have 3-4 tags each, there is a real risk of hitting the database 3-4 times per song, just to render a list of urls. This could be a major performance issue as the songs collection scales. If you aren't familiar with eager loading in Active Record, there are a few great optimization articles in http://rubyweekly.com/issues/252.

Once the tag values are ready, use your view file to place them on the rendered page via html data attributes:

```ruby
<% @songs.each do |song| %>
  <li class="song-link" data-difficulty="<%= song.difficulty.name" %>"...
<% end %>
```

For the genres data attribute - the tag category where I needed a many-to-many relationship - I created a simple "genres_string" method on the Song class, returning a comma separated string of the song's genres. If your data is too complicated for this approach, you can store any valid JSON directly in html data attributes.

### Step 3: Add Filter Buttons

Creating the buttons is straight forward:

```ruby
<% @difficulties.each do |difficulty| %>
  <button data-category="difficulty"...
<% end %>
```

In this simplification of the code I omitted the CSS classes used for styling, but it's important to note that all of the buttons other than difficulty, start out with a CSS class of "hidden" (display: none;).

### Step 4: Listen for Clicks

Now for the jQuery code. If your version of Rails has Turbolinks enabled (and I do recommend using it), you will want to listen for click events on both document ready and on page load:

```javascript
$(document).on('ready page:load', function(){
  songButtonListener();
});
```
When a button is clicked, hide all of the songs to clear previous filters. Then, reveal only the songs corresponding to the tag that generated the click event. The songButtonListener function grabs data from the button that was clicked (aka $(this) ) and then calls the hide and show functions:

```javascript
var hideSongs = function() {
  $('.song-link').addClass('hidden');
}

var showSongs = function(buttonCategory, dataCategory) {
  $(".song-link[data-" + dataCategory +"*='" + buttonCategory +"']").removeClass('hidden');
}
```

Note that the *= operator is used to check songs for tag names. This is a simple solution to parsing the comma separated list of genres, but it may lead to false positive matches if the list of tag values were to grow significantly in the future. If that were an issue, I would split the list into an array and check for exact matches instead.

### Step 5: Rinse & Repeat

Now that we have buttons that hide and show songs, we need to filter the buttons themselves. The logic for accomplishing this is virtually identical to steps 3 & 4 above. Create an unordered list of filter categories, with one of them "active" by default. When a filter category is clicked, add the ".hidden" class to all buttons, and remove it from buttons corresponding to the category that was clicked. I grabbed the category like this:

```javascript
var category = clickedTab.find('a').text().toLowerCase();
```
Before I close this post out, I want to point out the importance of keeping the filter UI as small and intuitive as possible. The initial problem of an overwhelmingly large index of songs would not be better replaced by an overwhelmingly large and complicated filtering interface.

Use only the vertical space you need to present filter categories in an uncluttered manner, and try to restrict the number of categories to whatever extent you can.

Thanks for reading!

-Adam