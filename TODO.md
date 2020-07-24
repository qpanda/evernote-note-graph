# Improvements
* Currently notes are fetched sequentially which can take considerable amount of time. Goroutines could be used to fetch several notes in parallel
* Most parameters are currently passed by pass-by-value, we could optimize further by passing pointers instead
* Parameter validation should be added to a number of functions across the codebase
* Should consider to setup Travis CI (or a similar GitHub app)
* Should explore which GitHub security apps support Go and set it up