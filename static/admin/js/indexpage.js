function useCode() {
    $.confirm({
        title: '提示',
        content: '' +
            '<form action="" class="formName">' +
            '<div class="form-group">' +
            '<label>请输入扩容码</label>' +
            '<input type="text" autocomplete="off" class="name form-control" required />' +
            '</div>' +
            '</form>',
        buttons: {
            formSubmit: {
                text: '提交',
                btnClass: 'btn-blue',
                action: function () {
                    var code = this.$content.find('.name').val();
                    if(!code){
                        $.alert('请您输入扩容码');
                        return false;
                    }
                    $.post("/system/use_code",{"code":code},function (res) {
                        if (res.resultCode == 200){
                            $.alert('扩容成功');
                        }else{
                            $.alert(res.msg);
                        }
                    })
                }
            },
            cancel: {
                text: '取消'
            },
        },
        onContentReady: function () {
            var jc = this;
            this.$content.find('form').on('submit', function (e) {
                e.preventDefault();
                jc.$$formSubmit.trigger('click');
            });
        }
    });
}