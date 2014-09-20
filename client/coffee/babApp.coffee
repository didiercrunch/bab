
root = this;

angular.module('bab', [
  'ngRoute'
  'bab.filters'
  'bab.services'
  'bab.directives'
  'bab.controllers'
]).
config(['$routeProvider', ($routeProvider) ->
  $routeProvider.when('/', {templateUrl: 'partials/home.html', controller: 'homeCtrl'})
  $routeProvider.when('/login', {templateUrl: 'partials/login.html', controller: 'loginCtrl'})
  $routeProvider.when('/about', {templateUrl: 'partials/about.html', controller: 'aboutCtrl'})
  $routeProvider.otherwise({redirectTo: '/'})
])

root.filters = angular.module('bab.filters', [])
root.services = angular.module('bab.services', [])
root.directives = angular.module('bab.directives', [])
root.controllers = angular.module('bab.controllers', [])


root.filters.filter 'month',  () -> 
        (text) ->
            return ['janvier',
                    'février',
                    'mars',
                    'avril',
                    'may',
                    'juin',
                    'juillet',
                    'aout',
                    'september',
                    'octobre',
                    'novembre',
                    'decembre'][Number(text) - 1]