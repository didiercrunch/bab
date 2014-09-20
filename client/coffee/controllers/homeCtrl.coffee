root = this;


root.controllers.controller('homeCtrl', ['$scope', '$http', ($scope, $http) ->
    $scope.services = []
    $http.get('/webapps').then (resp) ->
        $scope.services = resp.data
])