# adamcz.com

This is my personal website - primarily a place to host blog posts about music and tech. It works like this:

- I write posts as markdown files in the musicPosts and techPosts directories
- Committing a new post kicks off a fresh deploy
- When the app starts up, posts are parsed and held in memory (sort of half-way between a static site generator and a traditional database driven crud app)
- Emphasis on avoiding dependencies ([zero Go dependencies](https://github.com/aczerepinski/adamcz/blob/master/go.mod) and no JavaScript at all except for code syntax highlighting) and making content creation as simple as possible