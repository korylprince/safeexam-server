var app = angular.module("app", ["ngRoute", "ngCookies", "ngMaterial"]);

app.config(["$routeProvider", function($routeProvider) {
    $routeProvider
        .when("/login", {
            templateUrl: "views/login.html",
            controller: "loginController",
        }).when("/code", {
            templateUrl: "views/code.html",
            controller: "codeController",
        }).when("/help", {
            templateUrl: "views/help.html",
            controller: "helpController",
        }).otherwise({ redirectTo: "/login" });
}]);

app.config(["$mdThemingProvider", function($mdThemingProvider) {
    $mdThemingProvider.theme("default")
        .primaryPalette("blue");
}]);

app.factory("session", ["$cookies", function($cookies) {
    return { 
        setID: function(id) {
            $cookies.sessionID = id;
        },
        getID: function() {
            return $cookies.sessionID || "";
        },
        deleteID : function() {
            delete $cookies.sessionID;
        }
    };
}]);

app.factory("alert", function() {
    return {
        hidden: true,
        message: "",
    };
});

app.controller("helpController", ["$scope", "$location", function($scope, $location) {
    $scope.back = function() {
        $location.path("/code");
    };
}]);

app.controller("loginController", ["$scope", "$http", "$location", "session", "alert", function($scope, $http, $location, session, alert) {
    $scope.submit = function(login) {
        $scope.alert.hidden = true;

        $http({
            method: "POST",
            url: "api/2.0/auth",
            data: login,
            headers: {
                "Accept": "application/json",
            },
        }).success(function(data, status) {
            if (status != 200) {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
                console.log("Login error: ", status, data);
                return;
            }
            if (data.SessionID == null || data.SessionID == "") {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
                console.log("Login error: ", status, data);
                return;
            }
            session.setID(data.SessionID);
            $location.path("/code");

        }).error(function(data, status) {
            $scope.alert.hidden = false;
            if (status == 401) {
                $scope.alert.message = "Bad username or password";
            } else {
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
            }
            console.log("Login error: ", status, data);
        });
    };

    $scope.login = {};

    $scope.alert = alert;

    if (session.getID() != "") {
        $location.path("/code");
    }
}]);

app.controller("codeController", ["$scope", "$http", "$location", "$interval", "session", "alert", function($scope, $http, $location, $interval, session, alert) {

    $scope.logout = function(expired) {
        if (expired) {
            $scope.alert.hidden = false;
            $scope.alert.message = "Your session expired. Please log in again";
        }
        session.deleteID();
        $location.path("/login");
    }

    $scope.updateCode = function() {
        if (typeof $scope.sessionID != "string" || $scope.sessionID == "") {
            return;
        }
        $http({
            method: "GET",
            url: "api/2.0/code",
            headers: {
                "Accept": "application/json",
                "X-Session-Key": $scope.sessionID,
            },
        }).success(function(data, status) {
            if (status != 200) {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
                console.log("Code retrieval error: ", status, data);
                return;
            }
            if (data.Code == null || data.Code == "" || data.Expires == null || data.Expires == "" || data.ServerTime == null || data.ServerTime == "") {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
                console.log("Code retrieval error: ", status, data);
                return;
            }
            $scope.data.code = data.Code;
            $scope.data.expires = Math.round(data.Expires / 1000000);
            $scope.data.offset = Math.round((new Date().getTime()) - (data.ServerTime / 1000000));
            $scope.increment()

        }).error(function(data, status) {
            if (status == 401) {
                $scope.logout(true);
            } else {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
            }
            console.log("Code retrieval error: ", status, data);
        });
        $scope.data.lastCheck = new Date().getTime();
    };

    $scope.increment = function() {
        $scope.data.now = new Date().getTime();
        $scope.updateDiff($scope.data.expires - $scope.data.now + $scope.data.offset);
        if ($scope.data.expires < ($scope.data.now + $scope.data.offset) && $scope.data.lastCheck + 1000 < $scope.data.now) {
            $scope.updateCode();
        }
    };

    $scope.updateDiff = function(diff) {
        $scope.d.hours = Math.floor(diff / (60 * 60 * 1000));
        $scope.d.minutes = Math.floor(diff / (60 * 1000)) % 60;
        $scope.d.seconds = Math.floor(diff / 1000) % 60;
        $scope.d.hS = $scope.d.hours == 1 ? "" : "s";
        $scope.d.mS = $scope.d.minutes == 1 ? "" : "s";
        $scope.d.sS = $scope.d.seconds == 1 ? "" : "s";
    };

    // setup data
    $scope.sessionID = session.getID();

    $scope.alert = alert;
    $scope.alert.hidden = true;

    $scope.data = {
        code: "", //current code
        expires: 0, //time code expires on server in UTC millis
        now: new Date().getTime(), //now in UTC millis
        offset: 0, //difference in local clock and server in millis
        lastCheck: new Date().getTime(), //last updateCode was run
    };

    $scope.d = {
        hours: 0,
        minutes: 0,
        seconds: 0,
        hS: "",
        mS: "",
        sS: "",
    };

    // check for login
    if ($scope.sessionID == "") {
        $scope.logout();
    }

    // fetch code
    $scope.updateCode();

    // start increment
    if ($scope.loop == null) {
    $scope.loop = $interval($scope.increment, 1000);
        $scope.$on('$destroy', function () {
            $interval.cancel($scope.loop);
        });
    }
}]);
