
template = """
<div class="panel">
    <div class="row">
        <div class="large-12 columns" style="text-align:center;">
            <a href="{{ service.url }}"><img src="{{ service.image }}" style="height:100px;"></img></a>
        </div>
    </div>
    <div class="row">
        <div class="large-12 columns" style="text-align:center;">
             <h3 ><a href="{{ service.url }}"> {{ service.name }} </a></h3>
        </div>
    </div>
</div>
"""

directive = ($rootScope) ->
    directive = 
        template: template
        replace: true,
        transclude: false,
        restrict: 'E',
        scope:
            service: "="
        controller:["$scope",
            ($scope, $timeout) ->
                
            ]
    return directive
                        
directives.directive("babService", ['$rootScope', directive])
