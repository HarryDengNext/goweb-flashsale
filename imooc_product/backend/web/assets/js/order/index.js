;
var account_index_ops = {
    init:function () {
        this.eventBind();
    },
    eventBind:function () {
        var that = this;
        $('.delete').click(function () {
            that.ops('delete', $(this).attr('data'))
        });
    },

    ops:function (act, id) {
        var callback = {
            'ok': function () {
                $.ajax({
                    url: '/order/delete',
                    type: 'POST',
                    data: {
                        id: id,
                    },
                    dataType: 'json',
                    success: function (res) {
                        console.log(res)
                        var callback = null;
                        if (res.code == 200) {
                            callback = function () {
                                window.location.href = '/order/all';
                            }
                        }
                        common_ops.alert(res.msg, callback);
                    }
                })
            },
            'cancel': null,
    };
        common_ops.confirm( (act == 'delete' ? '确定删除?':'确认撤销'), callback );
    }
};
$(function () {
   account_index_ops.init();
});