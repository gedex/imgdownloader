imgdownloader
=============

Need thousand images for your project?

## Install

~~~text
$ go install github.com/gedex/imgdownloader
~~~

## Usage

* Download 1000 images, tagged with `cat,` from Flickr:

  ~~~text
  $ imgdownloader -tag cat -n 1000 -from flickr -o ./cats
  ~~~

* From Instagram:

  ~~~text
  $ imgdownloader -tag cat -n 1000 -from instagram -o ./cats
  ~~~

Currently, there are two providers: `flickr` and `instagram`. You specify
provider via `-from` parameter. Since getting list of images from providers
involves calling provider's REST API, you need to provide client credentials
for `imgdownloader` in `imgdownloader.json` either in current directory or
your `$HOME` directory with following format:

~~~json
{
	"flickr": {
		"api_key": "YOUR_FLICKR_API_KEY"
	},
	"instagram": {
		"access_token": "YOUR_ACCESS_TOKEN"
	}
}
~~~
