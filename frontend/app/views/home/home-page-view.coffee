View = require 'views/base/view'
OutingCollectionView = require 'views/home/outing-collection-view'

module.exports = class HomePageView extends View
  autoRender: true
  className: 'home-page'
  template: require './templates/home'

  initialize: ->
    super
    @delegate 'click', '.reload-btn', @clickReloadBtn
    @delegate 'submit','#outing-form', @submitOuting
    Chaplin.mediator.outingCollection.fetch()

  render: ->
    super
    @outingCollectionView = new OutingCollectionView collection: Chaplin.mediator.outingCollection
    @subview 'outingCollectionView', @outingCollectionView

  clickReloadBtn: (e)->
    Chaplin.mediator.outingCollection.fetch()
    @render()

  submitOuting: (e)->
    $form = @$('#outing-form')
    Chaplin.mediator.outingCollection.create({
        Name: $form.find('[name=Name]').val()
        Places: $form.find('[name=Places]').val()
    })
    @render()
    return false