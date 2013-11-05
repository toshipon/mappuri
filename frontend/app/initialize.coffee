Application = require 'application'
routes = require 'routes'

$.ajaxPrefilter( ( options, originalOptions, jqXHR )->
    options.url = 'http://localhost:8000/' + options.url
    options.crossDomain =
        crossDomain: true
    options.xhrFields =
        withCredentials: true
)

# Initialize the application on DOM ready event.
$ ->
  new Application {
    title: 'Brunch example application',
    controllerSuffix: '-controller',
    routes
  }
