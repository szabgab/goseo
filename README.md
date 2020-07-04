# GoSEO




## Plan

* Given a hostname crawl all the pages and check various things

* Check if page returns 200
* Get all the external links and check if they all return 200
* Get all the internal links and add them to the queue
* Check for <html lang=""> tag? (Allow configuring what value to expect)
* Check for <link rel="canonical" href="..." />
* Check for <title></title>
* Check for <meta name="description" content="">
* Check for <meta name="viewport" ...>