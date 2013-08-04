'user strict';

angular.module('filters', []).
  filter('truncate', function () {
  return function (text, length, end) {
    if (isNaN(length))
      length = 10;

    if (end === undefined)
      end = "...";

    if (text.length <= length || text.length - end.length <= length) {
      return text;
    }
    else {
      return String(text).substring(0, length-end.length) + end;
    }

  };
});

var progMob = angular.module('progMob', ['ui.bootstrap', 'filters'])
progMob.controller('MobsCtrl', function($scope, $dialog, $http) {

  $http.get('/user/mobs.json').success(function(mobs) {
    $scope.mobsGroupByOwner = _.groupBy(mobs, function(mob) { return mob.Owner });
    $scope.currentMob = _.first(mobs);
  });

  $scope.setCurrentMob = function(mob) {
    $scope.currentMob = mob;
  };

  $scope.closeJoinDialog = function($scope, dialog) {
    $scope.join = function(repo) {
      dialog.close(repo)
    };

    $scope.close = function() {
      dialog.close()
    };
  };

  var t = '<div class="modal-header">'+
    '<button type="button" ng-click=close() class="close fui-cross" data-dismiss="modal" aria-hidden="true"></button>' +
    '<h3>Join Room</h3>'+
    '</div>'+
    '<div class="modal-body">'+
    '<p><input type="text" class="span6" ng-model="repo" placeholder="Enter a GitHub repository" autofocus="true" /></p>'+
    '</div>'+
    '<div class="modal-footer">'+
    '<a href="#" ng-click="join(repo)" class="btn btn-wide btn-primary" >Join</button>'+
    '</div>';

  $scope.opts = {
    backdrop: true,
    keyboard: true,
    backdropClick: true,
    template:  t, // OR: templateUrl: 'path/to/view.html',
    controller: $scope.closeJoinDialog
  };

  $scope.openJoinDialog = function() {
    var d = $dialog.dialog($scope.opts);
    d.open().then(function(repo) {
      if(repo){
        console.log(repo)
      }
    });
  };
});

progMob.controller('MessagesCtrl', function($scope) {
  $scope.messages = [{
    user:
      {
      login: "jingweno",
      avatarURL: 'https://secure.gravatar.com/avatar/41740fef48f778383596392a6b2276c8?d=https://a248.e.akamai.net/assets.github.com%2Fimages%2Fgravatars%2Fgravatar-user-420.png'
    },
    text: 'fddfdfdfdfd',
    createdAt: 'May 12 3:38 PM',
  }];

  $scope.sendMessage = function(message) {
    $scope.messages.push({
      user:
        {
        login: "jingweno",
        avatarURL: 'https://secure.gravatar.com/avatar/41740fef48f778383596392a6b2276c8?d=https://a248.e.akamai.net/assets.github.com%2Fimages%2Fgravatars%2Fgravatar-user-420.png'
      },
      text: message.text,
      createdAt: 'May 12 3:38 PM',
    });
    $scope.message = {}
  };
});
