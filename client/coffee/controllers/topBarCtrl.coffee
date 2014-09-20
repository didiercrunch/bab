root = this;

root.controllers.controller('topBarCtrl', ['$scope', ($scope) ->
    $scope.$on('$routeChangeSuccess', (next, current) ->
        $scope.path = current.originalPath
    )
 
    $scope.path = "home"
])