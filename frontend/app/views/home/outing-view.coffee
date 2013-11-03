View = require 'views/base/view'

module.exports = class OutingView extends View
    autoRender: yes
    className: 'outing-row'
    tagName: 'li'
    template: require './templates/outing'
