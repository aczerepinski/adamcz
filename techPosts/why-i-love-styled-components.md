*** title ***
Why I Love Styled Components

*** date ***
8/27/17

*** tags ***
javascript

*** description ***
Some love for CSS-in-JS in a React environment

*** body ***
If you aren't familiar with [styled-components](https://www.styled-components.com/), it's a smallish library (~12kB) that allows you to write CSS directly in React components.

Prior to React's surge in popularity, keeping HTML, CSS and JavaScript in distinct files was probably the most widely practiced and least controversial way to separate concerns on the front end. In fact, in React's early days, there was a fair amount of backlash over the JSX paradigm of mixing markup and JavaScript together.

But by now, anybody who has embraced React at all has mixed JS and HTML in hundreds of components. Why not add CSS to the mix? I had been hearing about styled-components and other "CSS in JS" options, and decided to give it a try while building [DbMaj7.com](http://www.dbmaj7.com). In this post I'll share some of the reasons that I've been enjoying it. But first, a quick look at how it works:

```javascript
import styled from 'styled-components'
import { breakpoints } from './../styles'

export const PageLayout = styled.section`
  padding: 0 1rem;
  @media (min-width: ${breakpoints.medium}) {
    padding: 4rem 2rem;
  }
`
```

So starting from the top, we import the styled-components library, and then define `PageLayout` as the result of a `styled.section` tagged template call. If you aren't familiar with ES6 template literals, head over to [MDN](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Template_literals) for a quick overview. TLDR, think of `styled.section` as a function that will receive your template literal as a string wrapped in an array. If your string contains `${interpolation}`, the interpolated bits will also be available in the array, ready to be parsed by the styled library.

Under the hood, styled accepts any valid dom element (the full list is [here](https://github.com/styled-components/styled-components/blob/master/src/utils/domElements.js)) as the tag name, and returns a React component with an auto-generated className that will correspond to CSS that styled will inject into the document `<head>`.

Now that I have a styled component, I can compose it with "normal" React components like this:

```javascript
const AboutPage = () => (
  <PageLayout>
    <PageTitle>DbMaj7.com</PageTitle>
    <ArticleMarkdown body={body}/>
  </PageLayout>
)
```

So why do I like it so much?

### Easier Naming
Naming CSS classes without offending the BEM/DoCSSa experts can sometimes get pretty tricky. Especially if you have a complicated widget with several wrappers, inner containers, etc. Did I already use `.widget-wrapper` somewhere in the app?

Of course you still need to name your styled-components, but since these names don't need to have global scope, you can name un-exported components without fear of name collisions or descriptions that make sense outside the context of the page you're writing them in.

### What separation of concerns?
In my own practice, I often felt that the "separation" between HTML and CSS was more theoretical than real. I'd write Sass mixins with the best intentions of de-coupling and re-using all over the place, and then look back at the end of a project to find basically one implementation, tied to very specific markup. Of course there are exceptions, but for me this has generally been true.

### Less grepping
How often do you find yourself grepping (or `ag`ing) for class names to find all of the html, css, and js that touches your component? If you notice markup and styles changing frequently in the same pull request, there's a good chance that time could be saved by just storing them together in the same file.

### JS!
JavaScript is a *much* more powerful and fully featured language than CSS. I know and love the language, but I'm still probably barely scratching the surface of what it really offers in the context of writing styles. For sure, it gives you the flexibility to define styles as a function of props without the indirection of using JS to assign modifier class names. As an example:

```javascript
margin-left: ${props => ['a', 'b', 'd', 'e', 'g'].indexOf(props.letter) >= 0 ? '-2.36%' : 'initial'};
```

### Safe to delete + fewer styles on the page
With traditional CSS, it's always difficult to find and delete unused code. With styled components, if nobody imports it, nobody uses it. Gaining the confidence to delete feels a lot easier to me. As an added benefit, if nobody renders the component, users aren't downloading the styles at all.

### What about performance?
Any hesitation I had (and truthfully still have) about CSS in JS relates to performance. Parsing template literals to inject CSS into the document takes a non-zero amount of time, and I worry that this execution time will become noticeable when the number of styled components in an app hits some threshold number.

That's somewhat of an unfounded concern, because so far everything has worked very well for me, and I haven't come across anyone else complaining about performance yet either. Styled's website lists some big name companies using the library, but I haven't (yet) seen the tell-tale class names on their most valuable pages.

If anyone has had any negative experiences, I'd be especially interested in hearing about it. Please hit me up on twitter.

Otherwise, that's all I've got for now. Thanks for reading!

-Adam