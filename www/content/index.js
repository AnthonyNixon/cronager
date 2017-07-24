angular.module('cronagerApp', ['ngMaterial', 'ngRoute', 'mdDataTable'])

    .config(['$mdThemingProvider', '$routeProvider', '$locationProvider', function ($mdThemingProvider, $routeProvider, $locationProvider) {
        //Set up theme for the site using Target Red
        $mdThemingProvider.theme('default')
            .primaryPalette('blue')
            .accentPalette('amber');

        $routeProvider
            .when('/', {
                templateUrl: template("home")
            });

        function template(page) {
            return 'views/' + page + '.html';
        }

        $locationProvider.html5Mode(true);
        $locationProvider.hashPrefix('!');
    }])

    // Set up routeCtrl to act as the controller for route operations
    .controller('rootCtrl', ['$scope', '$location', '$window', '$http', '$rootScope',
        function ($scope, $location, $window, $http, $rootScope) {
            $rootScope.goto = gotoInternal;
            $rootScope.gotoExternal = gotoExternal;
            // function to go to an internal link
            // takes in a link variable
            // redirects the page to an internal link
            function gotoInternal(link) {
                console.log('goto ' + link);
                $location.path(link);
            }

            // function to go to an external link
            // takes in a link variable
            // redirects the page to the external link
            function gotoExternal(link) {
                $window.open(link);
            }

        }])

    .controller('tableCtrl', ['$scope', '$http',
        function ($scope, $http) {
            $scope.data = [];
            $scope.loaded = false;

            $http.get("http://localhost:3000/jobs")
                .then(function(response){
                    console.log(response.data);
                    $scope.data = response.data.result;
                    console.log($scope.data);
                    $scope.loaded = true;
                });

        }]);