CollectionView = require 'views/base/collection-view'
OutingView = require 'views/home/outing-view'

module.exports = class OutingCollectionView extends CollectionView
    autoRender: yes
    container: '.outings-render-section'
    className: 'outing-list'
    tagName: 'ul'
    template: require './templates/outing-collection'
    itemView: OutingView
