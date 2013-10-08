var app = angular.module('EmailTracker', ['ngResource']).
  config(function($interpolateProvider) {
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
  });

app.factory('Code', function($resource) {
  return $resource('/api/codes/:codeId', {}, {
    track: { method: "POST", params: {track: true} }
  });
});

app.controller("CodesCtrl", function(Code, $scope) {
  $scope.codes = Code.query()

  $scope.create = function() {
    var c = new Code($scope.code);
    c.$save(
      function(data) {
        $scope.codes.unshift(data)
        $scope.notice = "Successfully saved new tracking code."
        $scope.code.title = null
      }, function(resp) {
      if (resp.data.Error) {
          $scope.notice = resp.data.Message
          $scope.validationError = resp.data.Type == 'validation'
          $scope.savingError = resp.data.Type == 'error'
        }
      }
    )
  }

  $scope.track = function(codeId) {
    $scope.codes.forEach(function(code) {
      if (code.Id == codeId) {
        code.$track({id: codeId}, function(d) {
          code.Started = true
        });
      }
    });
  }

  $scope.destroy = function(codeId) {
    if (confirm('Do you really want to remove this tracking code ?')) {
      $scope.codes.forEach(function(code, index) {
        if (code.Id === codeId) {
          code.$remove({id: codeId}, function(d) {
            $scope.codes.splice(index, 1)
          })
        }
      });
    }
  }
});

//function CodesListCtrl($scope, Codes) {
//    $scope.codes = Codes
//}