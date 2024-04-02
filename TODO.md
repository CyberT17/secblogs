# TODO

- [x] `daten` needs to be sorted by `Date`
  - The `Dates` in `daten` struct currently is not sorted, this causes the resulting index.html to be displayed in random order.
- [x] Add more blogs sources
- [ ] Add better error handling
- [ ] Can we maintain order of blogs on the page
  - Right now every time the go code is run order in which the blogs are listed on the page changes. Could we somehow maintain the order? Is this due to concurrency?

## Future Features

- [ ] Consider not running the go code in github actions to generate the index.html.
- [ ] Light/Dark mode toggle
- [ ] Easily accessible list of sources
- [ ] RSS feed to import into some rss reader?
- [ ] Any improvements to make the site more accessible?
  - Look into ADA compliance standards
