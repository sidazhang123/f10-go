//onLoad - populate options
$(document).ready(function () {

    $(document)
        .ajaxStart(function () {
            $('#loading').show();
        })
        .ajaxStop(function () {
            $('#loading').hide();
        })
        .ajaxError(function () {
            $('#loading').hide();
        });
    $.ajax({
        method: "POST",
        dataType: "json",
        headers: {'sno': '321razad'},
        contentType: "application/json",
        url: window.location.href.replace(/\/$/, "") + "/mgmt/load_options",
        data: JSON.stringify({}),
        success: function (data) {
            var res = '<option value="nil">Select...</option>\n';
            data.path.split(";").forEach(function (item) {
                res += '<option value="' + item + '">' + item + '</option>\n'
            });

            $('#plugin_selector').html(res);
        },
        timeout: 3000,
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            $('#testRes').html(XMLHttpRequest.responseText);
        }
    });

});


//onSelect - load source code
$(document).ready(function () {
    $(document).on('change', '#plugin_selector', function (e) {
        if (this.value === "nil") {
            return
        }
        e.preventDefault();

        $.ajax({
            method: "POST",
            dataType: "json",
            contentType: "application/json",
            headers: {'sno': '321razad'},
            url: window.location.href.replace(/\/$/, "") + "/mgmt/select",
            data: JSON.stringify({pluginPath: this.value}),
            success: function (data) {
                $('#sourceCode').val(data.sourceCode);

            },
            timeout: 3000,
            error: function (XMLHttpRequest, textStatus, errorThrown) {
                $('#testRes').html(XMLHttpRequest.responseText);

            }
        });

    });
});


//Test onClick
$(document).ready(function () {
    $("#test_btn").click(function (e) {
        e.preventDefault();

        var url = window.location.href.replace(/\/$/, "") + "/mgmt/test";

        $.ajax({
            method: "POST",
            dataType: "json",
            headers: {'sno': '321razad'},
            contentType: "application/json",
            url: url,
            data: JSON.stringify({pluginPath: $('#plugin_selector').val(), testStr: $('#testStr').val()}),
            success: function (data) {
                if (data.err === "") {
                    $('#testRes').html(data.msg);
                } else {
                    $('#testRes').html(data.err);
                }

            },

            timeout: 3000,
            error: function (XMLHttpRequest, textStatus, errorThrown) {
                $('#testRes').html(XMLHttpRequest.responseText);
            }
        });

    });
});
//Update onClick
$(document).ready(function () {
    $("#update_btn").click(function (e) {
        e.preventDefault();

        var url = window.location.href.replace(/\/$/, "") + "/mgmt/update";

        $.ajax({
            method: "POST",
            dataType: "json",
            headers: {'sno': '321razad'},
            contentType: "application/json",
            url: url,
            data: JSON.stringify({pluginPath: $('#plugin_selector').val(), sourceCode: $('#sourceCode').val()}),
            success: function (data) {
                if (data.err === "") {
                    $('#testRes').html(data.msg);
                } else {
                    $('#testRes').html(data.err);
                }

            },
            timeout: 3000,
            error: function (XMLHttpRequest, textStatus, errorThrown) {
                $('#testRes').html(XMLHttpRequest.responseText);
            }
        });

    });
});
