Collection = require 'models/base/collection'
Outing = require 'models/outing'

module.exports = class OutingCollection extends Collection
    model: Outing
    url: 'outings'
