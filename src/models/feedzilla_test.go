package models

import (
	"testing"
)

var doRequestArticles = func(url string) ([]byte, error) {
	return []byte(`
		{
				"articles": [
						{
								"publish_date": "Tue, 01 Oct 2013 15:33:00 +0100",
								"source": "All Africa",
								"source_url": "http://allafrica.com/tools/headlines/rdf/sport/headlines.rdf",
								"summary": "[The Star]Violet Makuto with Esther Mwombe under 23 Kenya Ladies Volleyball\nteam players trainig at Kasarani Stadium.\n\n",
								"title": "Kenya: National Under-23s Ready for Top Challenge (All Africa)",
								"url": "http://news.feedzilla.com/en_us/stories/top-news/331285407?count=2&client_source=api&format=json"
						},
						{
								"publish_date": "Tue, 01 Oct 2013 14:07:00 +0100",
								"source": "All Africa",
								"source_url": "http://allafrica.com/tools/headlines/rdf/religion/headlines.rdf",
								"summary": "[This Day]From Ibrahim Abdullahi's remarks in Lagos two weeks ago, it became\nobvious that one needed to still clarify some confusion that arose from the\nonline debates that trailed, \"The Catholic I was\", which recently appeared on\nthis page. Then, last Thursday, a mail came from Frank Ihekwoaba, CEO of Nima\nCapital Advisory Partners, on the same subject matter. It was the\ncorrespondence of a decent man who, thinking that the columnist was the person\nwho disparaged Mbaise people in the online comments, simply could n\n\n",
								"title": "Nigeria: Mbaise, Fani-Kayode and 'Alaiyemore' (All Africa)",
								"url": "http://news.feedzilla.com/en_us/stories/top-news/331285392?count=2&client_source=api&format=json"
						}
				],
				"description": "Top News",
				"syndication_url": "http://news.feedzilla.com/en_us/news/top-news.rss?count=2&client_source=api",
				"title": "Feedzilla: Top News"
		}
		`), nil
}

var doRequestCategories = func(url string) ([]byte, error) {
	return []byte(`
		[
			{
				"category_id":1314,
				"display_category_name":"",
				"english_category_name":"Sports",
				"url_category_name":""
			},
			{
				"category_id":13,
				"display_category_name":
				"Art",
				"english_category_name":"Art","url_category_name":"art"
			}
		]
		`), nil
}

func TestFeedzillaCategories(t *testing.T) {
	doRequest = doRequestCategories
	categories, err := FeedzillaCategories()
	if err != nil {
		t.Fatal(err)
	}
	for _, cat := range categories {
		t.Logf("%+v", cat)
	}
}

func TestFeedzillaArticles(t *testing.T) {
	doRequest = doRequestArticles
	articles, err := FeedzillaArticles(26, 2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", articles.List[0])
	for _, art := range articles.List {
		t.Logf("%+v", art)
	}
}
