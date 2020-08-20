function submitForm() {
    lightyear.loading('show');  // 显示
    $.ajax({
        //几个参数需要注意一下
        type: "POST",//方法类型
        dataType: "json",//预期服务器返回的数据类型
        url: "/system/profile" ,//url
        data: $('#pro_form').serialize(),
        success: function (result) {
            lightyear.loading('hide');  // 隐藏
            console.log(result);//打印服务端返回的数据(调试用)
            if (result.resultCode == 200) {
                lightyear.notify('保存成功~', 'success', 1000, 'mdi mdi-emoticon-happy', 'top', 'center');
                setTimeout(function () {
                    location.reload();
                },1000)
            }else{
                lightyear.notify('异常！'+result.errorMsg, 'warning', 1000, 'mdi mdi-emoticon-happy', 'top', 'center');
            }
        },
        error : function(result) {
            lightyear.loading('hide');  // 隐藏
            lightyear.notify('异常！' + result.errorMsg, 'warning', 1000, 'mdi mdi-emoticon-happy', 'top', 'center');
        }
    });
}

layui.use('upload', function(){
    var $ = layui.jquery
        ,upload = layui.upload;

    //普通图片上传
    var uploadInst = upload.render({
        elem: '#picker_one_pic'
        ,url: '/system/profile_upload' //改成您自己的上传接口
        ,accept: 'image' //普通文件
        ,exts: 'jpeg|jpg|gif|png' //只允许上传压缩文件
        ,size: 2 * 1024 //限制文件大小，单位 KB
        ,done: function(res){
            if(res.resultCode == 200){
                $('#file_list_one_pic').attr('src',res.msg);
            }else{
                return layer.msg('上传失败' + res.msg);
            }
        }
    });
});