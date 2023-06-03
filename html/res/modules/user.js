layui.define('form', function (exports) {
  var $ = layui.$
    , layer = layui.layer
    , laytpl = layui.laytpl
    , setter = layui.setter
    , view = layui.view
    , admin = layui.admin
    , form = layui.form
    , router = layui.router()
    , search = router.search;

  //提交
  form.on('submit(LAY-user-login-submit)', function (data) {
    //请求登入接口
    admin.req({
      url: '/mj/user/v1/email/login' //实际使用请改成服务端真实接口
      , type: "post"
      , data: JSON.stringify(data.field)
      , contentType: 'application/json'
      , done: function (res) {
        //请求成功后，写入 access_token
        layui.data(setter.tableName, {
          key: setter.request.tokenName,
          value: res.data.access_token
        });

        //登入成功的提示与跳转
        layer.msg('登录成功', {
          offset: '15px'
          , icon: 1
          , time: 1000
        }, function () {
          console.log(search.redirect)
          location.hash = search.redirect ? decodeURIComponent(search.redirect) : '/';
        });
      }
    });

  });

  //对外暴露的接口
  exports('user', {});
});