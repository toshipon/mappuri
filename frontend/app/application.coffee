# The application object.
OutingCollection = require 'models/outing-collection'

module.exports = class Application extends Chaplin.Application
  start: ->
    # You can fetch some data here and start app
    # (by calling `super`) after that.
    super
    @initMediator()
    Chaplin.mediator.outingCollection.fetch()

  initMediator: ->
    Chaplin.mediator.outingCollection = new OutingCollection
    Chaplin.mediator.seal()

