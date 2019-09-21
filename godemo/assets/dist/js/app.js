
$(function () {
    $(document).on('click', '.show-todolist-modal', function () {
        $('#todolist-modal').modal('show');
    });
});

$(function () {
    $(document).on('click', '.btn-info', function () {
        $('#todolist-modal').modal('show');
    });
});

$(function () {
    $(document).on('click', '.save-change', function () {
        var title = $('.input-title').val();
        var text = $('textarea').val();

        if (title == '' || text == '') {
            return;
        }

        $.ajax({
            url: location.href.split('?')[0] + 'postTodo',
            contentType: "application/json",
            type: 'POST',
            data: JSON.stringify({ 'title': title.trim(), 'text': text.trim() }),
            dataType: 'json',
            timeout: 10000,
            cache: false,
            scriptCharset: 'utf-8',
        }).then(
            $('.group-todo').append("<li class='list-group-item'><h4 class='list listgroup-item-heading'>" + title + "</h4>\
                <p class='list-group-item-text text-wrap'>" + text + "</p><div class='buttons'>\
                <button type='button' class='btn btn-info btn-xs' data-toggle='modal' data-target='#exampleModal' title='Edit'>\
                <i class='fas fa-edit fa-xs' aria-hidden='true'></i></button><button type='button' class='btn btn-danger btn-xs move-done' title='Done'>\
                <i class='fas fa-check fa-xs clas-change' aria-hidden='true'></i></button></div></li>")
        ).then(
            $('#todolist-modal').modal('hide')
        );
    });
});

$(function () {
    $(document).on('click', '.move-done', function () {
        var done_title = $(this).parents('li').children('h4').text();
        var done_text = $(this).parents('li').children('p').text();

        $.ajax({
            url: location.href.split('?')[0] + 'done',
            contentType: "application/json",
            type: 'PATCH',
            data: JSON.stringify({ 'title': done_title.trim(), 'text': done_text.trim() }),
            dataType: 'json',
            timeout: 10000,
            cache: false,
            scriptCharset: 'utf-8',
        }).then(
            $(this).removeClass('move-done')
        ).then(
            $(this).children('.class-change').removeClass('fa-check')
        ).then(
            $(this).children('.class-change').addClass('fa-minus')
        ).then(
            $(this).addClass('delete-done')
        ).then(
            $(this).parents('li').appendTo('.group-done')
        );
    });
});

$(function () {
    $(document).on('click', '.delete-done', function () {
        var del_title = $(this).parents('li').children('h4').text();
        var del_text = $(this).parents('li').children('p').text();

        $.ajax({
            url: location.href.split('?')[0] + 'delete',
            contentType: "application/json",
            type: 'DELETE',
            data: JSON.stringify({ 'title': del_title.trim(), 'text': del_text.trim() }),
            dataType: 'json',
            timeout: 10000,
            cache: false,
            scriptCharset: 'utf-8',
        }).then(
            $(this).parents('li').remove()
        );
    });
});
