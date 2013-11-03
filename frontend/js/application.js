define(['chaplin'], function(Chaplin) {
  'use strict';

  // The application object
  // Choose a meaningful name for your application
  var Mappuri = Chaplin.Application.extend({
    // Set your application name here so the document title is set to
    // “Controller title – Site title” (see Layout#adjustTitle)
    title: 'Mappuri',
    start: function() {
      // You can fetch some data here and start app
      // (by calling supermethod) after that.
      Chaplin.Application.prototype.start.apply(this, arguments);
    }
  });

  return Mappuri;
});
