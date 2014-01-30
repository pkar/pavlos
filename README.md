# Pavlos

##requirements
```bash
brew install mongo
mongo
```

```bash
brew install go
```

```bash
sudo easy_install pymongo
sudo easy_install fabric
sudo easy_install requests
```

## Info
Everything is pulled from feedzilla as such
```bash
fab packages
GOPATH=`pwd` go run src/main/pavlos.go -logtostderr=true -load=true
```
should take not too long.


## Running
```bash
GOPATH=`pwd` go run src/main/pavlos.go -logtostderr=true
```

## API
To get recommendations, you have to provide feedback 
after getting one.
```bash
curl http://localhost:8000/relevant/{user}
```

like can be 1, 0, -1
```bash
curl --data '{"id":"524c68fe238843b2242237d8", "like": -1}' "http://localhost:8000/collect/pavlos"
```

## Example
with the server running and loaded with data.


### Get something relevant
```bash
curl http://localhost:8000/relevant
{"ID":"524c68fe238843b2242237d8","CategoryID":1314,"publish_date":"Wed, 02 Oct 2013 16:01:00 +0100","source":"IKSurfmag","source_url":"http://www.facebook.com/feeds/page.php?id=5862208995&format=rss20","summary":"CrazyFly ripper Tommy Gaunt puts together his diary from Australia, certainly\nlooks like the boys had a pretty good time over there. The West Coast never\nfails to deliver it seems! I couldn't turn down the offer of going back to\nAustralia so in November myself,...\n\n\n\nFine Times\n\nwww.iksurfmag.com\n\nCrazyFly ripper Tommy Gaunt puts together his diary from Australia, certainly\nlooks like the boys had a pretty good time over there. The West Coast never\nfails to deliver it seems! I couldn't turn down the offer of going back to\nAustralia so in November myself,...\n\n","title":"CrazyFly ripper Tommy Gaunt puts together his diary from Australia, certainly lo... (IKSurfmag)","url":"http://news.feedzilla.com/en_us/stories//331528759?count=3&client_source=api&format=json"}
```

If you curl it again you should get back the same thing. 
Thats because it is circling through the list of categories in a non random way initially(that could be changed).
It does this until the user has given feedback for all categories at which point 
whats recommended is based on a random weighted graph.

### Give feedback on the item seen to move to the next category.
```bash
curl --data '{"id":"524c68fe238843b2242237d8", "like": -1}' "http://localhost:8000/collect/pavlos"
OK
```

### Train a user to speed up the process
train pavlos to like travel, it will present things to the user and provide feedback to 
the server.
```bash
fab train_user:pavlos,keywords="greece;athens;travel"

# Train another user with similar interest in travel, but also in sports.
fab train_user:tom,keywords="sports;football;baseball;hockey;travel;cuba;china;russia;syria"

# and one more
fab train_user:barry,keywords="jesus;god;satan;religion"
```

stop that and take a look at the weights generated for the user.
```bash
fab dump_user:pavlos

********************************************************************************
pavlos
********************************************************************************
Health Weight: 0.01 Score: -8 Count: 10
General Weight: 0.018691588785 Score: -6 Count: 7
Hobbies Weight: 0.018691588785 Score: -6 Count: 9
Entertainment Weight: 0.0467289719626 Score: -3 Count: 7
Fun Stuff Weight: 0.0373831775701 Score: -4 Count: 8
Events Weight: 0.0280373831776 Score: -5 Count: 17
Industry Weight: 0.0280373831776 Score: -5 Count: 9
Internet Weight: 0.0373831775701 Score: -4 Count: 11
IT Weight: 0.0280373831776 Score: -5 Count: 5
Jobs Weight: 0.0560747663551 Score: -2 Count: 9
Law Weight: 0.0280373831776 Score: -5 Count: 12
Life Style Weight: 0.0373831775701 Score: -4 Count: 20
Music Weight: 0.018691588785 Score: -6 Count: 14
Oddly Enough Weight: 0.0280373831776 Score: -5 Count: 11
Politics Weight: 0.00934579439252 Score: -7 Count: 11
Products Weight: 0.018691588785 Score: -6 Count: 14
Art Weight: 0.018691588785 Score: -6 Count: 15
 Weight: 0.0280373831776 Score: -5 Count: 15
Business Weight: 0.018691588785 Score: -6 Count: 13
Blogs Weight: 0.0373831775701 Score: -4 Count: 11
Columnists Weight: 0.018691588785 Score: -6 Count: 10
Celebrities Weight: 0.0280373831776 Score: -5 Count: 10
Programming Weight: 0.018691588785 Score: -6 Count: 10
Religion And Spirituality Weight: 0.00934579439252 Score: -7 Count: 13
Science Weight: 0.018691588785 Score: -6 Count: 14
Shopping Weight: 0.018691588785 Score: -6 Count: 9
Society Weight: 0.0280373831776 Score: -5 Count: 12
Sports Weight: 0.0280373831776 Score: -5 Count: 8
Travel Weight: 0.990654205607 Score: 98 Count: 98
Top News Weight: 0.0467289719626 Score: -3 Count: 13
Top Blogs Weight: 0.018691588785 Score: -6 Count: 12
Technology Weight: 0.0280373831776 Score: -5 Count: 8
Video Games Weight: 0.018691588785 Score: -6 Count: 12
Video Weight: 0.0280373831776 Score: -5 Count: 14
USA Weight: 0.018691588785 Score: -6 Count: 13
Universities Weight: 0.0373831775701 Score: -4 Count: 9
World News Weight: 0.0280373831776 Score: -5 Count: 14

barry Euclidean: 0.341044375662 Pearson: -0.0582403747528
tom Euclidean: 0.526879295797 Pearson: 0.700992587527
```

Note that the topic of travel is highly weighted on a scale of 0 to 1.
Also note that the Count is much higher than any other topics count,
thats because the training script focused exclusively on travel and
put the others in the background.

Compared to tom and barry, tom might provide more insight into 
discovering better new content.

```bash
fab dump_user:tom
********************************************************************************
tom
********************************************************************************
Health Weight: 0.0222222222222 Score: -2 Count: 7
General Weight: 0.01 Score: -4 Count: 5
Hobbies Weight: 0.01 Score: -4 Count: 6
Entertainment Weight: 0.0111111111111 Score: -3 Count: 3
Fun Stuff Weight: 0.0222222222222 Score: -2 Count: 3
Events Weight: 0.01 Score: -4 Count: 10
Industry Weight: 0.0333333333333 Score: -1 Count: 8
Internet Weight: 0.0111111111111 Score: -3 Count: 4
IT Weight: 0.0222222222222 Score: -2 Count: 3
Jobs Weight: 0.0111111111111 Score: -3 Count: 9
Law Weight: 0.0333333333333 Score: -1 Count: 7
Life Style Weight: 0.0222222222222 Score: -2 Count: 2
Music Weight: 0.0222222222222 Score: -2 Count: 5
Oddly Enough Weight: 0.0222222222222 Score: -2 Count: 8
Politics Weight: 0.0333333333333 Score: -1 Count: 7
Products Weight: 0.0111111111111 Score: -3 Count: 7
Art Weight: 0.0111111111111 Score: -3 Count: 6
 Weight: 0.0444444444444 Score: 0 Count: 10
Business Weight: 0.0333333333333 Score: -1 Count: 3
Blogs Weight: 0.0333333333333 Score: -1 Count: 8
Columnists Weight: 0.0111111111111 Score: -3 Count: 5
Celebrities Weight: 0.0333333333333 Score: -1 Count: 11
Programming Weight: 0.0111111111111 Score: -3 Count: 3
Religion And Spirituality Weight: 0.0444444444444 Score: 0 Count: 8
Science Weight: 0.01 Score: -4 Count: 7
Shopping Weight: 0.0222222222222 Score: -2 Count: 5
Society Weight: 0.0111111111111 Score: -3 Count: 8
Sports Weight: 0.444444444444 Score: 36 Count: 70
Travel Weight: 0.988888888889 Score: 85 Count: 85
Top News Weight: 0.0111111111111 Score: -3 Count: 5
Top Blogs Weight: 0.0333333333333 Score: -1 Count: 2
Technology Weight: 0.0111111111111 Score: -3 Count: 7
Video Games Weight: 0.01 Score: -4 Count: 7
Video Weight: 0.0222222222222 Score: -2 Count: 4
USA Weight: 0.0222222222222 Score: -2 Count: 4
Universities Weight: 0.0111111111111 Score: -3 Count: 6
World News Weight: 0.0111111111111 Score: -3 Count: 4


pavlos Euclidean: 0.844303640707 Pearson: 0.910396552925
barry Euclidean: 0.328381310652 Pearson: -0.0278609146664
```

Notice that now Sports and Travel have been presented to the user most frequently.
tom has similar tastes mostly to pavlos.


```bash
fab dump_user:barry

********************************************************************************
barry
********************************************************************************
Health Weight: 0.0208333333333 Score: -3 Count: 5
General Weight: 0.0208333333333 Score: -3 Count: 5
Hobbies Weight: 0.0104166666667 Score: -4 Count: 8
Entertainment Weight: 0.0208333333333 Score: -3 Count: 4
Fun Stuff Weight: 0.0208333333333 Score: -3 Count: 3
Events Weight: 0.0416666666667 Score: -1 Count: 3
Industry Weight: 0.0208333333333 Score: -3 Count: 8
Internet Weight: 0.0208333333333 Score: -3 Count: 4
IT Weight: 0.0208333333333 Score: -3 Count: 9
Jobs Weight: 0.0104166666667 Score: -4 Count: 5
Law Weight: 0.0104166666667 Score: -4 Count: 6
Life Style Weight: 0.03125 Score: -2 Count: 9
Music Weight: 0.0416666666667 Score: -1 Count: 6
Oddly Enough Weight: 0.0208333333333 Score: -3 Count: 8
Politics Weight: 0.0416666666667 Score: -1 Count: 7
Products Weight: 0.0104166666667 Score: -4 Count: 8
Art Weight: 0.0208333333333 Score: -3 Count: 4
 Weight: 0.03125 Score: -2 Count: 4
Business Weight: 0.03125 Score: -2 Count: 6
Blogs Weight: 0.0104166666667 Score: -4 Count: 7
Columnists Weight: 0.03125 Score: -2 Count: 7
Celebrities Weight: 0.0104166666667 Score: -4 Count: 7
Programming Weight: 0.0208333333333 Score: -3 Count: 4
Religion And Spirituality Weight: 0.989583333333 Score: 90 Count: 90
Science Weight: 0.0208333333333 Score: -3 Count: 4
Shopping Weight: 0.0208333333333 Score: -3 Count: 6
Society Weight: 0.0208333333333 Score: -3 Count: 7
Sports Weight: 0.01 Score: -5 Count: 8
Travel Weight: 0.0104166666667 Score: -4 Count: 8
Top News Weight: 0.01 Score: -5 Count: 6
Top Blogs Weight: 0.03125 Score: -2 Count: 4
Technology Weight: 0.0208333333333 Score: -3 Count: 7
Video Games Weight: 0.03125 Score: -2 Count: 4
Video Weight: 0.0104166666667 Score: -4 Count: 11
USA Weight: 0.0104166666667 Score: -4 Count: 6
Universities Weight: 0.01 Score: -5 Count: 7
World News Weight: 0.03125 Score: -2 Count: 4

pavlos Euclidean: 0.344573426174 Pearson: -0.0487744927358
tom Euclidean: 0.269674905828 Pearson: -0.0268624719847

```

Now barry likes religion mostly. He's not very similar to pavlos or tom, I think 
he may be "different".

## Next step, is to factor in similarites of users to make actual recommendations from people.
This is a simple step of just comparing similarity scores of users and 
adjusting the Category weight or even just recommending a specific article.

## Note this is naive in that the like dislike is just based on keyword search from a python client.

One thing to try is to scrape a users social media for keywords and try things out, but that's not 
that smart.

One other thing is to keep timestamps of when the user clicked and to adjust the weights accordingly.
So if I read a lot of politics in the morning, that would be more heavily weighted 
and shown more often to me then, but later in the evening it was technology, then that would
show more.

There are more ideas but chances are no one even got here so I'll stop...
