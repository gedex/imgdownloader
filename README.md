imgdownloader
=============

Need hundred of images for your project?

## Install

~~~text
$ go get github.com/gedex/imgdownloader
~~~

## Usage

<iframe src="http://showterm.io/03c3947a3833b3e0325fe#fast" frameborder="0" style="overflow:hidden;height:100%;width:100%" height="100%" width="100%"></iframe>

* Download 1000 images, tagged with `cat`, from Flickr:

  ~~~text
  $ imgdownloader -tag="cat" -n=1000 -from="flickr" -o="./cats"
  ~~~

* From Instagram:

  ~~~text
  $ imgdownloader -tag="cat" -n=1000 -from="instagram" -o="./cats"
  ~~~

* From Google Picasa Web:

  ~~~text
  $ imgdownloader -tag="cat" -n=1000 -from="picasa" -o="./cats"
  ~~~

Currently, there are three providers: `flickr`, `instagram`, and `picasa`. You specify
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
	},
	"picasa": {
		"": ""
	}
}
~~~

To get Instagram's `access_token` you can use [ginsta](https://github.com/gedex/ginsta):

~~~text
$ go get github.com/gedex/ginsta
$ ginsta token_get # Open localhost:8080 In your browser
~~~

Currently you need to specify `picasa` config as a placeholder so that the config reader
doesn't returns an error.
