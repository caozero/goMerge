/**
 * Created by caoping on 2017/2/4.
 */

__ = {};


__.init = function () {
    __.setFontSize();
    __.log.init();
    __.project.init();
    __.watcher.init();
    __.bindEvent();
    __.net.init();
};

__.setFontSize = function () {
    __.config.fontSize = 100;
};

__.start = function () {
    __.init();
};

__.net = {
    init: function () {
        this.initScoket();
    }
    , initScoket: function () {
        var _this = this;
        this.socket = new WebSocket('ws://127.0.0.1:6301/dataSocket');
        this.socket.onopen = function (p) {
            __.send({cmd: 'getProjectList'});
            __.send({cmd: 'getWatcherList'});
        };
        this.socket.onmessage = function (p) {
            _this.onMsg(p);
        };
        this.socket.onerror = function (p) {
            console.log('websocket 错误!');
        };
        this.socket.onclose = function (p) {
            console.log('websocket 关闭!');
            __.log.add({status: 2, msg: "websocket 关闭!"});
        };
    }
    , sendMsg: function (p) {
        this.socket.send(JSON.stringify(p));
    }
    , onMsg: function (p) {
        try {
            var m = JSON.parse(p.data);
            if (m) {
                __.log.add(m);
                if (m.cmd in __)__[m.cmd](m);
                //__.rount.addMsg(m);
            }
        } catch (e) {
            console.dir(e);
        }

    }
    , reconnect: function () {

    }
};

__.log = {

    init: function () {
        this.el = $("#logPlane");
        this.table = this.el.find('table');
    }
    , add: function (p) {
        var msg = '';
        if (typeof p == 'object') {
            var c = "success";
            switch (p.status) {
                case    1:
                    c = "warning";
                    break;
                case    2:
                    c = "error";
                    break;
            }
            msg = $('<tr class="' +
                c + '">' +
                '<td>' + p['msg'] + '</td>' +
                '</tr>');
        } else {
            msg = $('<tr>' +
                '<td>' + p + '</td>' +
                '</tr>');
        }

        this.table.prepend(msg);
    }
};

__.send = function (p) {
    __.net.sendMsg(p)
};

__.config = {
    ProjectList: null
    , WatcherList: null
    , init: function () {

    }

};

__.makeEl = function () {
    var b = $('.mainContent');
    var pl = __.config.ProjectList;
    var l = pl.length;
    while (l--) {
        var p = pl[l];
        var s = '<div class="project">'
            + '<p>' + p.fileName + '</p>' +
            '</div>';
        var d = $(s);
        b.append(d);
    }
    var wl = __.config.WatcherList;
    l = wl.length;
    while (l--) {
        var w = wl[l];
        var d = $('<div class="watcher">'
            + '<p>' + w.Root + '</p>' +
            '</div>');
        b.append(d);
    }
    __.setElPosition();
};

__.setElPosition = function () {
    var i = 0;
    $('.mainContent .project').each(function () {
        var _t = $(this);
        var x = 100;
        var y = 60 * i + 20;
        _t.stop().animate({left: x, top: y}, 'fast');
        i++;
    });
    i = 0;
    $('.mainContent .watcher').each(function () {
        var _t = $(this);
        var x = 600;
        var y = 60 * i + 20;
        _t.stop().animate({left: x, top: y}, 'fast');
        i++;
    })
};


__.project = {

    init: function () {
        this.makeEl();
    }
    , makeEl: function () {
        this.el = $('<section id="projectList"></section>');
        $('body').append(this.el);
    }
    , loadList: function (p) {
        console.dir(p);
        var d = null;
        try {
            d = JSON.parse(p.data);
        } catch (e) {
            __.log.add({status: 2, msg: '项目列表数据错误!'});
            return;
        }
        if (!d)return;
        this.list = d;
        this.el.empty();
        for (var i in d) {
            this.addProject(d[i]);
        }
    }
    , addProject: function (p) {
        var isUpdate=p['IsUpdate']?' isUpdate':'';
        var d = $('<div class="project'+isUpdate+'" id="' + p['Hex'] + '">' +
            '<p class="projectTitle">' + p['FileName'] + '</p>' +
            '<p class="topBar"><span class="button update">U</span></p>' +
            '<ul class="mto">' +
            '</ul>' +
            '</div>');
        this.el.append(d);
        var mto = d.find('ul.mto');
        for (var i in p['mto']) {
            var m = p['mto'][i];
            var li = $('<li>' +
                '<p class="title">' + m['mergeTo'] + '</p>' +
                '</li>');
            mto.append(li);
            var ul = $('<ul class="srcList"></ul>');
            for (var n in m.src) {
                var s = m.src[n];
                li2 = $('<li>' +
                    '' + s +
                    '</li>')
                ul.append(li2);
            }
            li.append(ul);
        }

    }
    ,onUpdate:function (p) {
        console.dir(p);
        var hex=p.data['ProjectHex'];
        $('#'+hex).removeClass('isUpdate');
        $('li[target-id='+hex+']').parents('li').find('p.title').removeClass('isUpdate');
    }
};

__.watcher = {
    init: function () {
        this.makeEl();
    }

    , makeEl: function () {
        this.el = $('<section id="watcherList"></section>');
        $('body').append(this.el);
    }
    , loadList: function (p) {
        console.dir(p);
        var d = null;
        try {
            d = JSON.parse(p.data);
        } catch (e) {
            __.log.add({status: 2, msg: '监控列表数据错误!'});
            return;
        }
        if (!d)return;
        this.list = d;
        this.el.empty();
        for (var i in d) {
            this.addProject(d[i]);
        }
    }
    , addProject: function (p) {
        var d = $('<div class="watcher" id="' + p['Hex'] + '">' +
            '<p class="watcherTitle">' + p['Root'] + '</p>' +
            '<ul class="folder">' +
            '</ul>' +
            '</div>');
        this.el.append(d);
        var mto = d.find('ul.folder');
        for (var i in p['FileList']) {
            var m = p['FileList'][i];
            var li = $('<li>' +
                '<p class="title" id="' + m['Hex'] + '">' + i + '</p>' +
                '</li>');
            mto.append(li);
            var ul = $('<ul class="linkPoint"></ul>');
            for (var n in m['ProjectHex']) {
                var s = m['ProjectHex'][n];
                li2 = $('<li target-id="' + s + '">' +
                    '</li>');
                ul.append(li2);
            }
            li.append(ul);
        }
    }
    , onFileModify: function (p) {
        var file = $('#' + p['data']['Hex']);
        file.addClass('isUpdate');
        file.parent('li').find('ul.linkPoint li').each(function () {
            var _t = $(this);
            var hex = _t.attr('target-id');
            $('#' + hex).addClass('isUpdate');
        })
    }
};

__.onFileModify = function (p) {
    console.dir(p);
    __.watcher.onFileModify(p)
};

__.getProjectList = function (p) {
    __.project.loadList(p)
};
__.getWatcherList = function (p) {
    __.watcher.loadList(p)
};
__.onUpdate = function (p) {
    __.project.onUpdate(p)
};

__.bindEvent = function () {
    $('.addWatchFile').click(function () {
        var s = $(this).parent().find('input[name=watchFile]').val();
        if (!s.length)return;
        __.send({cmd: 'addWatchFile', data: s});
    });
    $('body').on('click', '.button', function () {
        var _t = $(this);
        if (_t.hasClass('update')) {
            var id = _t.parents('.project').attr('id');
            __.send({cmd: 'update', hex: id});
        }
        console.log('button click');
        return false;
    })
};

$(function () {
    __.start();
});
