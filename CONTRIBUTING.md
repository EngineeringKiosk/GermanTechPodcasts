# Contribution Guide

## How to add a new podcast to the list

Please feel free to open a pull request if you know of a German tech podcast that is not mentioned here.  
If you're in doubt if a particular podcast is a good fit for the list, **don't open an issue, but create a pull request right away** because that's easier to handle. Thanks! :smiley:

### Criteria to be accepted

Each podcast on the list should adhere to the following criteria

- is a podcast :smiley:
- is published using a RSS feed and is not behind a paywall
- published in German :de: :austria: :switzerland:
- has published at least 3 consecutive episodes
- is at least 2 months old
- majority of content is `tech` (tech context should be clearly visible and more than 60% of all episodes)
- reasonable audio quality
- is still active (new episode published within two months) or has timeless/relevant content (e.g. a news episode won't be relevant months later)

Podcasts can be removed from the list if a criteria is no longer met.

#### What is `tech` in this context?

When it deals with things like ...

- Software, hardware, data science, Artificial intelligence (AI), user interface (UI)/user experience (UX), architecture, algorithms, security
- Dev culture, team culture of tech/product teams
- Worklife of devs and tech people (e.g. coping with stress, imposter, etc.)
- Research of topics mentioned above
- Tech influence/impact on society (Artificial intelligence (AI), security, etc)

#### What is not `tech` in this context?

When it deals (mainly) with things like ...

- Economy/finance podcasts about tech companies
- Tech topics explained for non-tech people

### Format

:warning: **The main [`README.md`](/README.md) is a rendered version of the data. Do not edit it manually.**

To add a new podcast, please create a `.yml` file in the [`/podcasts`](/podcasts) directory like `/podcasts/<podcast-name>.yml`. 
Feel free to check out a few other YAML files in that directory to see how it should look like.

| Field                        | Type               | Description                                                                                |
|------------------------------|--------------------|--------------------------------------------------------------------------------------------|
| name                         | string             | Name of the podcast                                                                        |
| website                      | string (url)       | Website (full URL) of the podcast                                                          |
| podcastIndexID               | integer            | ID of the podcast in the [PodcastIndex](https://podcastindex.org/)                         |
| rssFeed                      | string (url)       | RSS feed that is used to publish new podcast episodes                                      |
| spotify                      | string (url)       | Spotify link (full URL) of the podcast show                                                |
| description                  | string             | Description of the podcast (in German)                                                     |
| tags                         | array (of strings) | List of tags / focus themes                                                                |
| weekly_downloads_avg         | -                  | Average weekly downloads (substructure)                                                    |
| weekly_downloads_avg.value   | integer            | Number of the average weekly downloads                                                     |
| weekly_downloads_avg.updated | string (date)      | Date in format YYYY-MM-DD on when the "number of the average weekly downloads" was updated |

Finally, create a pull request with all your changes. 

**Thanks for helping out!** :tada:
