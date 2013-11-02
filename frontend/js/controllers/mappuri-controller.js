define([
  'controllers/base/controller',
  'models/mappuri',
  'views/mappuri-view'
], function(Controller, Mappuri, MappuriView) {
  'use strict';

  var MappuriController = Controller.extend({
    show: function(params) {
      this.model = new Mappuri();
      this.view = new MappuriView({
        model: this.model,
        region: 'main'
      });
    }
  });

  return MappuriController;
});
